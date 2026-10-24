package content

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func search(param string, cols ...string) Filter {
	return func(r *http.Request) (string, []any) {
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		cat := r.URL.Query().Get(param)
		var conds []string
		var args []any
		if q != "" && len(q) <= 100 {
			args = append(args, "%"+likeEscaper.Replace(q)+"%")
			var or []string
			for _, c := range cols {
				or = append(or, c+"::text ILIKE $1")
			}
			conds = append(conds, "("+strings.Join(or, " OR ")+")")
		}
		if cat != "" && param != "" && len(cat) <= 40 {
			args = append(args, cat)
			conds = append(conds, param+" = $"+itoa(len(args)))
		}
		return strings.Join(conds, " AND "), args
	}
}

// likeEscaper stops user input from acting as ILIKE wildcards (expensive full scans, data probing).
var likeEscaper = strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`)

func itoa(i int) string { return string(rune('0' + i)) }

func eventFilter(r *http.Request) (string, []any) {
	from, err1 := time.Parse("2006-01-02", r.URL.Query().Get("from"))
	to, err2 := time.Parse("2006-01-02", r.URL.Query().Get("to"))
	switch {
	case err1 == nil && err2 == nil:
		return "ends_at >= $1 AND starts_at < $2", []any{from, to.AddDate(0, 0, 1)}
	case err1 == nil:
		return "ends_at >= $1", []any{from}
	}
	return "", nil
}

// Register mounts every CMS resource. reindex is called after knowledge/faq changes.
func Register(mux *http.ServeMux, db *pgxpool.Pool, editor func(http.Handler) http.Handler, reindex func(context.Context)) {
	onKB := func(ctx context.Context, _ int64) { reindex(ctx) }

	(&Resource[Knowledge]{DB: db, Name: "knowledge", Table: "knowledge",
		Columns:    []string{"doc_code", "category", "kind", "title", "body", "source", "informant", "reviewer", "status", "cover_url", "video_url", "sort"},
		PublicCond: "status = 'published'", Order: "sort, id", HasUpdated: true,
		Filter: search("kind", "title", "body"), OnChange: onKB}).Mount(mux, "knowledge", editor)

	(&Resource[Place]{DB: db, Name: "place", Table: "places",
		Columns:    []string{"code", "category", "name", "description", "lat", "lon", "open_hours", "visit_minutes", "fee", "contact", "tags", "images", "source", "published"},
		PublicCond: "published", Order: "id", HasUpdated: true,
		Filter: search("category", "name", "description")}).Mount(mux, "places", editor)

	(&Resource[Route]{DB: db, Name: "route", Table: "routes",
		Columns:    []string{"name", "description", "place_ids", "minutes", "published"},
		PublicCond: "published", Order: "id", HasUpdated: true}).Mount(mux, "routes", editor)

	(&Resource[Event]{DB: db, Name: "event", Table: "events",
		Columns:    []string{"title", "description", "preparation", "starts_at", "ends_at", "place_id", "cover_url", "contact", "source", "published"},
		PublicCond: "published", Order: "starts_at", HasUpdated: true, Filter: eventFilter}).Mount(mux, "events", editor)

	(&Resource[Shop]{DB: db, Name: "shop", Table: "shops",
		Columns:    []string{"name", "kind", "products", "description", "lat", "lon", "contact", "link", "images", "published"},
		PublicCond: "published", Order: "id", HasUpdated: true,
		Filter: search("kind", "name", "products")}).Mount(mux, "shops", editor)

	(&Resource[News]{DB: db, Name: "news", Table: "news",
		Columns:    []string{"title", "body", "cover_url", "video_url", "is_banner", "published_at"},
		PublicCond: "published_at IS NOT NULL AND published_at <= now()", Order: "published_at DESC NULLS FIRST, id DESC", HasUpdated: true,
		Filter: search("", "title", "body")}).Mount(mux, "news", editor)

	(&Resource[FAQ]{DB: db, Name: "faq", Table: "faq",
		Columns:    []string{"category", "question", "answer", "doc_code", "reviewer", "published"},
		PublicCond: "published", Order: "category, id", OnChange: onKB}).Mount(mux, "faq", editor)
}
