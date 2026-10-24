// Package content serves the CMS entities (places, knowledge, events, shops, news, routes, faq).
//
// Each entity is a Go struct with a Validate method. Writes are decoded strictly into that struct,
// validated, re-marshalled and inserted with jsonb_populate_record, so only whitelisted columns can
// ever be written and SQL stays generic.
package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/auth"
	"github.com/reclover/sinakhontao/internal/httpx"
)

type Validator interface{ Validate() error }

// Filter turns query parameters into an extra SQL condition. Placeholders start at $1.
type Filter func(r *http.Request) (cond string, args []any)

type Resource[T Validator] struct {
	DB         *pgxpool.Pool
	Name       string   // audit entity name
	Table      string   // trusted constant
	Columns    []string // writable columns, trusted constants
	PublicCond string   // condition for public listing
	Order      string
	HasUpdated bool
	Filter     Filter
	OnChange   func(ctx context.Context, id int64) // e.g. re-index RAG
}

func (res *Resource[T]) cols() string { return strings.Join(res.Columns, ", ") }

// source is what a read selects from: the public site sees only id, the writable columns and
// updated_at, so a column added to the table later is never exposed by accident.
func (res *Resource[T]) source(public bool) string {
	if !public {
		return res.Table
	}
	cols := "id, " + res.cols()
	if res.HasUpdated {
		cols += ", updated_at"
	}
	return fmt.Sprintf("(SELECT %s FROM %s)", cols, res.Table)
}

func (res *Resource[T]) list(w http.ResponseWriter, r *http.Request, public bool) error {
	conds := []string{"TRUE"}
	var args []any
	if public && res.PublicCond != "" {
		conds = append(conds, res.PublicCond)
	}
	if res.Filter != nil {
		if c, a := res.Filter(r); c != "" {
			// Renumber placeholders: filters are written from $1.
			for i := len(a); i >= 1; i-- {
				c = strings.ReplaceAll(c, fmt.Sprintf("$%d", i), fmt.Sprintf("$%d", len(args)+i))
			}
			conds = append(conds, c)
			args = append(args, a...)
		}
	}
	q := fmt.Sprintf(`SELECT coalesce(jsonb_agg(to_jsonb(t) ORDER BY %s), '[]') FROM (SELECT * FROM %s x WHERE %s LIMIT 500) t`,
		res.Order, res.source(public), strings.Join(conds, " AND "))
	var out json.RawMessage
	if err := res.DB.QueryRow(r.Context(), q, args...).Scan(&out); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write(out)
	return err
}

func (res *Resource[T]) get(w http.ResponseWriter, r *http.Request, public bool) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	cond := "TRUE"
	if public && res.PublicCond != "" {
		cond = res.PublicCond
	}
	var out json.RawMessage
	err = res.DB.QueryRow(r.Context(), fmt.Sprintf(`SELECT to_jsonb(t) FROM %s t WHERE id=$1 AND %s`, res.source(public), cond), id).Scan(&out)
	if errors.Is(err, pgx.ErrNoRows) {
		return httpx.NotFound()
	}
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(out)
	return err
}

func (res *Resource[T]) decode(w http.ResponseWriter, r *http.Request) ([]byte, error) {
	var v T
	if err := httpx.Decode(w, r, &v); err != nil {
		return nil, err
	}
	if err := v.Validate(); err != nil {
		return nil, httpx.BadRequest(err.Error())
	}
	return json.Marshal(v)
}

func (res *Resource[T]) create(w http.ResponseWriter, r *http.Request) error {
	body, err := res.decode(w, r)
	if err != nil {
		return err
	}
	var id int64
	q := fmt.Sprintf(`INSERT INTO %[1]s (%[2]s) SELECT %[2]s FROM jsonb_populate_record(NULL::%[1]s, $1) RETURNING id`, res.Table, res.cols())
	if err := res.DB.QueryRow(r.Context(), q, body).Scan(&id); err != nil {
		return mapPgError(err)
	}
	res.changed(r, "create", id)
	return httpx.JSON(w, http.StatusCreated, map[string]int64{"id": id})
}

func (res *Resource[T]) update(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	body, err := res.decode(w, r)
	if err != nil {
		return err
	}
	set := fmt.Sprintf(`(%[1]s) = (SELECT %[1]s FROM jsonb_populate_record(NULL::%[2]s, $2))`, res.cols(), res.Table)
	if res.HasUpdated {
		set += ", updated_at = now()"
	}
	tag, err := res.DB.Exec(r.Context(), fmt.Sprintf(`UPDATE %s SET %s WHERE id=$1`, res.Table, set), id, body)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return httpx.NotFound()
	}
	res.changed(r, "update", id)
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (res *Resource[T]) delete(w http.ResponseWriter, r *http.Request) error {
	id, err := httpx.PathID(r)
	if err != nil {
		return err
	}
	tag, err := res.DB.Exec(r.Context(), fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, res.Table), id)
	if err != nil {
		return mapPgError(err)
	}
	if tag.RowsAffected() == 0 {
		return httpx.NotFound()
	}
	res.changed(r, "delete", id)
	return httpx.JSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (res *Resource[T]) changed(r *http.Request, action string, id int64) {
	if u := auth.FromContext(r.Context()); u != nil {
		auth.Audit(r.Context(), res.DB, u.ID, action, res.Name, id, nil)
	}
	if res.OnChange != nil {
		res.OnChange(r.Context(), id)
	}
}

func mapPgError(err error) error {
	s := err.Error()
	switch {
	case strings.Contains(s, "SQLSTATE 23505"):
		return &httpx.Error{Status: http.StatusConflict, Code: "conflict", Message: "duplicate code"}
	case strings.Contains(s, "SQLSTATE 23514"), strings.Contains(s, "SQLSTATE 23503"):
		return httpx.BadRequest("value violates a constraint")
	}
	return err
}

// Mount registers public read routes under /api/v1/<path> and admin CRUD under /api/v1/admin/<path>.
func (res *Resource[T]) Mount(mux *http.ServeMux, path string, admin func(http.Handler) http.Handler) {
	mux.Handle("GET /api/v1/"+path, httpx.Wrap(func(w http.ResponseWriter, r *http.Request) error { return res.list(w, r, true) }))
	mux.Handle("GET /api/v1/"+path+"/{id}", httpx.Wrap(func(w http.ResponseWriter, r *http.Request) error { return res.get(w, r, true) }))
	mux.Handle("GET /api/v1/admin/"+path, admin(httpx.Wrap(func(w http.ResponseWriter, r *http.Request) error { return res.list(w, r, false) })))
	mux.Handle("GET /api/v1/admin/"+path+"/{id}", admin(httpx.Wrap(func(w http.ResponseWriter, r *http.Request) error { return res.get(w, r, false) })))
	mux.Handle("POST /api/v1/admin/"+path, admin(httpx.Wrap(res.create)))
	mux.Handle("PUT /api/v1/admin/"+path+"/{id}", admin(httpx.Wrap(res.update)))
	mux.Handle("DELETE /api/v1/admin/"+path+"/{id}", admin(httpx.Wrap(res.delete)))
}
