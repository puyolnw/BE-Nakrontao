package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/httpx"
)

const (
	cookieName  = "snt_session"
	sessionTTL  = 12 * time.Hour
	maxFailures = 5
	lockFor     = 15 * time.Minute
)

type User struct {
	ID     int64  `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

type Service struct {
	DB            *pgxpool.Pool
	SecureCookies bool
	loginLimit    *httpx.Limiter
}

func NewService(db *pgxpool.Pool, secure bool) *Service {
	return &Service{DB: db, SecureCookies: secure, loginLimit: httpx.NewLimiter(10, 10)}
}

type ctxKey struct{}

func FromContext(ctx context.Context) *User {
	u, _ := ctx.Value(ctxKey{}).(*User)
	return u
}

func hashToken(tok string) []byte {
	s := sha256.Sum256([]byte(tok))
	return s[:]
}

// Require loads the session user and enforces one of the given roles.
// Mutating requests must also carry X-Requested-With, which a cross-site form cannot set.
func (s *Service) Require(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Header.Get("X-Requested-With") != "fetch" {
				httpx.WriteError(w, r, httpx.Forbidden())
				return
			}
			c, err := r.Cookie(cookieName)
			if err != nil {
				httpx.WriteError(w, r, httpx.Unauthorized())
				return
			}
			var u User
			err = s.DB.QueryRow(r.Context(), `
				SELECT u.id, u.email, u.name, u.role, u.active FROM sessions s JOIN users u ON u.id = s.user_id
				WHERE s.token_hash = $1 AND s.expires_at > now() AND u.active`, hashToken(c.Value)).
				Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.Active)
			if err != nil {
				httpx.WriteError(w, r, httpx.Unauthorized())
				return
			}
			if len(roles) > 0 && !hasRole(u.Role, roles) {
				httpx.WriteError(w, r, httpx.Forbidden())
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, &u)))
		})
	}
}

func hasRole(role string, allowed []string) bool {
	for _, a := range allowed {
		if a == role {
			return true
		}
	}
	return false
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("X-Requested-With") != "fetch" {
		return httpx.Forbidden() // blocks login CSRF from cross-site forms
	}
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if !s.loginLimit.Allow(httpx.ClientIP(r)) {
		return httpx.TooMany()
	}
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	invalid := &httpx.Error{Status: http.StatusUnauthorized, Code: "invalid_credentials", Message: "email or password is incorrect"}

	var id int64
	var hash string
	var active bool
	var lockedUntil *time.Time
	err := s.DB.QueryRow(r.Context(), `SELECT id, password_hash, active, locked_until FROM users WHERE email=$1`, in.Email).
		Scan(&id, &hash, &active, &lockedUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		VerifyPassword(dummyHash, in.Password) // equalise timing for unknown accounts
		return invalid
	}
	if err != nil {
		return err
	}
	if lockedUntil != nil && lockedUntil.After(time.Now()) {
		VerifyPassword(dummyHash, in.Password) // same answer and timing as a wrong password
		return invalid
	}
	if !active || !VerifyPassword(hash, in.Password) {
		_, err := s.DB.Exec(r.Context(), `UPDATE users SET failed_logins = failed_logins + 1,
			locked_until = CASE WHEN failed_logins + 1 >= $2 THEN now() + $3::interval END WHERE id=$1`,
			id, maxFailures, lockFor.String())
		if err != nil {
			return err
		}
		return invalid
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return err
	}
	tok := base64.RawURLEncoding.EncodeToString(raw)
	exp := time.Now().Add(sessionTTL)
	if _, err := s.DB.Exec(r.Context(), `UPDATE users SET failed_logins=0, locked_until=NULL WHERE id=$1`, id); err != nil {
		return err
	}
	if _, err := s.DB.Exec(r.Context(), `INSERT INTO sessions(token_hash, user_id, expires_at) VALUES ($1,$2,$3)`, hashToken(tok), id, exp); err != nil {
		return err
	}
	s.DB.Exec(r.Context(), `DELETE FROM sessions WHERE expires_at < now()`)
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: tok, Path: "/api", Expires: exp,
		HttpOnly: true, Secure: s.SecureCookies, SameSite: http.SameSiteStrictMode})
	Audit(r.Context(), s.DB, id, "login", "user", id, nil)
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Service) Logout(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("X-Requested-With") != "fetch" {
		return httpx.Forbidden()
	}
	if c, err := r.Cookie(cookieName); err == nil {
		if _, err := s.DB.Exec(r.Context(), `DELETE FROM sessions WHERE token_hash=$1`, hashToken(c.Value)); err != nil {
			return err
		}
	}
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/api", MaxAge: -1,
		HttpOnly: true, Secure: s.SecureCookies, SameSite: http.SameSiteStrictMode})
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Service) Me(w http.ResponseWriter, r *http.Request) error {
	return httpx.JSON(w, http.StatusOK, FromContext(r.Context()))
}

func (s *Service) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	var in struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	u := FromContext(r.Context())
	var hash string
	if err := s.DB.QueryRow(r.Context(), `SELECT password_hash FROM users WHERE id=$1`, u.ID).Scan(&hash); err != nil {
		return err
	}
	if !VerifyPassword(hash, in.Current) {
		return httpx.BadRequest("current password is incorrect")
	}
	nh, err := HashPassword(in.New)
	if err != nil {
		return httpx.BadRequest(err.Error())
	}
	if _, err := s.DB.Exec(r.Context(), `UPDATE users SET password_hash=$2 WHERE id=$1`, u.ID, nh); err != nil {
		return err
	}
	// Invalidate every other session of this user.
	c, _ := r.Cookie(cookieName)
	s.DB.Exec(r.Context(), `DELETE FROM sessions WHERE user_id=$1 AND token_hash<>$2`, u.ID, hashToken(c.Value))
	Audit(r.Context(), s.DB, u.ID, "change_password", "user", u.ID, nil)
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Users management (admin only).

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) error {
	rows, err := s.DB.Query(r.Context(), `SELECT id, email, name, role, active FROM users ORDER BY id`)
	if err != nil {
		return err
	}
	users, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (User, error) {
		var u User
		return u, row.Scan(&u.ID, &u.Email, &u.Name, &u.Role, &u.Active)
	})
	if err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, users)
}

type userInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	Active   *bool  `json:"active"`
	Password string `json:"password"`
}

func (in *userInput) validate(create bool) error {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Name = strings.TrimSpace(in.Name)
	if !strings.Contains(in.Email, "@") || len(in.Email) > 200 {
		return httpx.BadRequest("invalid email")
	}
	if in.Name == "" || len(in.Name) > 200 {
		return httpx.BadRequest("name is required")
	}
	if in.Role != "admin" && in.Role != "editor" {
		return httpx.BadRequest("role must be admin or editor")
	}
	if create && in.Password == "" {
		return httpx.BadRequest("password is required")
	}
	return nil
}

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var in userInput
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if err := in.validate(true); err != nil {
		return err
	}
	hash, err := HashPassword(in.Password)
	if err != nil {
		return httpx.BadRequest(err.Error())
	}
	var id int64
	err = s.DB.QueryRow(r.Context(), `INSERT INTO users(email,name,role,password_hash) VALUES ($1,$2,$3,$4)
		ON CONFLICT (email) DO NOTHING RETURNING id`, in.Email, in.Name, in.Role, hash).Scan(&id)
	if errors.Is(err, pgx.ErrNoRows) {
		return &httpx.Error{Status: http.StatusConflict, Code: "conflict", Message: "email already exists"}
	}
	if err != nil {
		return err
	}
	Audit(r.Context(), s.DB, FromContext(r.Context()).ID, "create", "user", id, map[string]string{"email": in.Email, "role": in.Role})
	return httpx.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	var in userInput
	if err := httpx.Decode(w, r, &in); err != nil {
		return err
	}
	if err := in.validate(false); err != nil {
		return err
	}
	me := FromContext(r.Context())
	active := in.Active == nil || *in.Active
	if id == me.ID && (!active || in.Role != "admin") {
		return httpx.BadRequest("you cannot demote or suspend yourself")
	}
	tag, err := s.DB.Exec(r.Context(), `UPDATE users SET email=$2, name=$3, role=$4, active=$5 WHERE id=$1`,
		id, in.Email, in.Name, in.Role, active)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return httpx.NotFound()
	}
	if in.Password != "" {
		hash, err := HashPassword(in.Password)
		if err != nil {
			return httpx.BadRequest(err.Error())
		}
		if _, err := s.DB.Exec(r.Context(), `UPDATE users SET password_hash=$2 WHERE id=$1`, id, hash); err != nil {
			return err
		}
	}
	if !active || in.Password != "" {
		s.DB.Exec(r.Context(), `DELETE FROM sessions WHERE user_id=$1`, id)
	}
	Audit(r.Context(), s.DB, me.ID, "update", "user", id, map[string]any{"role": in.Role, "active": active})
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// Bootstrap creates the first admin when the users table is empty.
func Bootstrap(ctx context.Context, db *pgxpool.Pool, email, password string) error {
	if email == "" || password == "" {
		return nil
	}
	var n int
	if err := db.QueryRow(ctx, `SELECT count(*) FROM users`).Scan(&n); err != nil || n > 0 {
		return err
	}
	hash, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, `INSERT INTO users(email,name,role,password_hash) VALUES ($1,'Administrator','admin',$2)`,
		strings.ToLower(email), hash)
	return err
}

// Audit records an important change (scope 1.6). Failures are logged by the caller's request log only.
func Audit(ctx context.Context, db *pgxpool.Pool, userID int64, action, entity string, entityID int64, detail any) {
	var d []byte
	if detail != nil {
		d, _ = json.Marshal(detail)
	}
	db.Exec(ctx, `INSERT INTO audit_log(user_id, action, entity, entity_id, detail) VALUES ($1,$2,$3,$4,$5)`,
		userID, action, entity, strconv.FormatInt(entityID, 10), d)
}

var dummyHash, _ = HashPassword("timing-equaliser-password")
