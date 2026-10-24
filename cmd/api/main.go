package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/reclover/sinakhontao/internal/auth"
	"github.com/reclover/sinakhontao/internal/config"
	"github.com/reclover/sinakhontao/internal/content"
	"github.com/reclover/sinakhontao/internal/db"
	"github.com/reclover/sinakhontao/internal/eval"
	"github.com/reclover/sinakhontao/internal/httpx"
	"github.com/reclover/sinakhontao/internal/insights"
	"github.com/reclover/sinakhontao/internal/media"
	"github.com/reclover/sinakhontao/internal/rag"
	"github.com/reclover/sinakhontao/internal/recommend"
)

func main() {
	seed := flag.String("seed", os.Getenv("SEED_FILE"), "path to a SQL file to load once after migrating (demo data); default $SEED_FILE")
	flag.Parse()
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	if err := run(*seed); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run(seedFile string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	if err := httpx.SetTrustedProxies(cfg.TrustedProxy); err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool); err != nil {
		return err
	}
	if err := auth.Bootstrap(ctx, pool, cfg.BootstrapAdminEmail, cfg.BootstrapAdminPassword); err != nil {
		return err
	}
	if seedFile != "" {
		var n int
		pool.QueryRow(ctx, `SELECT count(*) FROM knowledge`).Scan(&n)
		if n == 0 {
			sql, err := os.ReadFile(seedFile)
			if err != nil {
				return err
			}
			if _, err := pool.Exec(ctx, string(sql)); err != nil {
				return err
			}
			slog.Info("seed loaded", "file", seedFile)
		}
	}

	var embedder rag.Embedder = rag.HashEmbedder{}
	minVec := 0.0 // offline embedder: keyword evidence required
	if cfg.EmbedURL != "" && cfg.EmbedKey != "" {
		embedder = &rag.OpenAIEmbedder{URL: cfg.EmbedURL, Key: cfg.EmbedKey, Model: cfg.EmbedModel}
		minVec = rag.DefaultMinVecSimRemote
	}
	var llm rag.LLM = rag.Extractive{}
	if cfg.AnthropicKey != "" {
		llm = &rag.Anthropic{Key: cfg.AnthropicKey, Model: cfg.AnthropicModel}
	}
	slog.Info("ai providers", "embedder", embedder.Name(), "llm", llm.Name())

	ragSvc := &rag.Service{
		DB: pool, Retriever: &rag.Retriever{DB: pool, Embedder: embedder}, LLM: llm, Salt: cfg.IPHashSalt, MinVecSim: minVec,
		PerClient: httpx.NewLimiter(cfg.ChatPerMinute, cfg.ChatPerMinute), Daily: httpx.NewDailyCap(cfg.ChatDailyCap),
	}
	if err := ragSvc.Reindex(ctx); err != nil {
		slog.Error("initial reindex failed", "err", err)
	}

	go func() {
		t := time.NewTicker(6 * time.Hour)
		defer t.Stop()
		for {
			rag.Purge(ctx, pool, cfg.RetentionDays)
			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
	}()

	authSvc := auth.NewService(pool, cfg.SecureCookies)
	editor := authSvc.Require("admin", "editor")
	admin := authSvc.Require("admin")
	ins := &insights.Handler{DB: pool, Salt: cfg.IPHashSalt, Track: httpx.NewLimiter(120, 60), Votes: httpx.NewLimiter(3, 3)}
	ev := &eval.Handler{DB: pool, RAG: ragSvc}
	rec := &recommend.Handler{DB: pool}
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "uploads"
	}
	med := &media.Handler{Dir: uploadDir}

	mux := http.NewServeMux()
	w := httpx.Wrap
	mux.HandleFunc("GET /api/health", func(rw http.ResponseWriter, r *http.Request) {
		if err := pool.Ping(r.Context()); err != nil {
			httpx.JSON(rw, http.StatusServiceUnavailable, map[string]string{"status": "db_down"})
			return
		}
		httpx.JSON(rw, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Public
	mux.Handle("GET /api/v1/settings", w(ins.GetSettings))
	mux.Handle("POST /api/v1/chat", w(ragSvc.Chat))
	mux.Handle("POST /api/v1/chat/{id}/feedback", w(ragSvc.Feedback))
	mux.Handle("POST /api/v1/plan", w(rec.Plan))
	mux.Handle("GET /api/v1/survey", w(ins.SurveyQuestions))
	mux.Handle("POST /api/v1/survey", w(ins.SubmitSurvey))
	mux.Handle("POST /api/v1/track", w(ins.TrackView))
	content.Register(mux, pool, editor, ragSvc.ReindexAsync)

	// Auth
	mux.Handle("POST /api/v1/auth/login", w(authSvc.Login))
	mux.Handle("POST /api/v1/auth/logout", w(authSvc.Logout))
	mux.Handle("GET /api/v1/auth/me", editor(w(authSvc.Me)))
	mux.Handle("POST /api/v1/auth/password", editor(w(authSvc.ChangePassword)))

	// Admin
	mux.Handle("GET /api/v1/admin/users", admin(w(authSvc.ListUsers)))
	mux.Handle("POST /api/v1/admin/users", admin(w(authSvc.CreateUser)))
	mux.Handle("PUT /api/v1/admin/users/{id}", admin(w(authSvc.UpdateUser)))
	mux.Handle("PUT /api/v1/admin/settings/{key}", admin(w(ins.PutSetting)))
	mux.Handle("GET /api/v1/admin/audit", admin(w(ins.AuditLog)))
	mux.Handle("POST /api/v1/admin/upload", editor(w(med.Upload)))
	mux.Handle("POST /api/v1/admin/reindex", editor(w(ragSvc.ReindexHandler)))
	mux.Handle("GET /api/v1/admin/dashboard", editor(w(ins.Dashboard)))
	mux.Handle("GET /api/v1/admin/survey-report", editor(w(ins.SurveyReport)))
	mux.Handle("GET /api/v1/admin/export/{kind}", editor(w(ins.Export)))
	mux.Handle("GET /api/v1/admin/eval/questions", editor(w(ev.ListQuestions)))
	mux.Handle("POST /api/v1/admin/eval/questions/import", editor(w(ev.ImportQuestions)))
	mux.Handle("GET /api/v1/admin/eval/runs", editor(w(ev.ListRuns)))
	mux.Handle("POST /api/v1/admin/eval/runs", admin(w(ev.StartRun)))
	mux.Handle("GET /api/v1/admin/eval/runs/{id}", editor(w(ev.RunResults)))
	mux.Handle("PUT /api/v1/admin/eval/results/{id}", editor(w(ev.Grade)))
	mux.Handle("GET /api/v1/admin/eval/export", editor(w(ev.ExportRuns)))

	mux.Handle("GET /uploads/", med.Serve())
	if cfg.StaticDir != "" {
		mux.Handle("GET /", spa(cfg.StaticDir))
	}

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           httpx.Logging(httpx.SecurityHeaders(httpx.APILimit(httpx.NewLimiter(300, 120), httpx.NewLimiter(40, 20), mux))),
		MaxHeaderBytes:    32 << 10,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
	go func() {
		<-ctx.Done()
		sctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		srv.Shutdown(sctx)
	}()
	slog.Info("listening", "addr", cfg.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// spa serves the built SvelteKit app and falls back to index.html for client-side routes.
// The shell's inline bootstrap script is allowed by its SHA-256 hash, so the CSP needs no 'unsafe-inline'.
func spa(dir string) http.Handler {
	mime.AddExtensionType(".webmanifest", "application/manifest+json")
	files := http.FileServer(http.Dir(dir))
	index := filepath.Join(dir, "index.html")
	var mu sync.Mutex
	var csp string
	var cspMod time.Time
	shellCSP := func() string { // recomputed only when index.html changes (e.g. a redeploy)
		mu.Lock()
		defer mu.Unlock()
		if st, err := os.Stat(index); err == nil && !st.ModTime().Equal(cspMod) {
			csp, cspMod = httpx.ShellCSP(index), st.ModTime()
		}
		return csp
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			httpx.WriteError(w, r, httpx.NotFound())
			return
		}
		p := filepath.Join(dir, filepath.FromSlash(filepath.Clean("/"+r.URL.Path)))
		if st, err := os.Stat(p); err == nil && !st.IsDir() && filepath.Base(p) != "index.html" {
			if strings.HasPrefix(r.URL.Path, "/_app/immutable/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			if r.URL.Path == "/sw.js" {
				w.Header().Set("Cache-Control", "no-cache")
			}
			files.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Security-Policy", shellCSP())
		http.ServeFile(w, r, index)
	})
}
