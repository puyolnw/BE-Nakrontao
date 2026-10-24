package content

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// I18n is a translated text: {"th": "...", "zh": "...", "en": "..."}. Thai is mandatory.
type I18n map[string]string

func (t I18n) check(field string, maxLen int, required bool) error {
	for k, v := range t {
		if k != "th" && k != "zh" && k != "en" {
			return fmt.Errorf("%s: unsupported language %q", field, k)
		}
		if len(v) > maxLen {
			return fmt.Errorf("%s.%s: too long", field, k)
		}
	}
	if required && strings.TrimSpace(t["th"]) == "" {
		return fmt.Errorf("%s.th is required", field)
	}
	return nil
}

// List is a string array that encodes nil as [] so NOT NULL array columns accept it.
type List []string

func (l List) MarshalJSON() ([]byte, error) {
	if l == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(l))
}

var codeRe = regexp.MustCompile(`^[A-Za-z0-9_-]{1,40}$`)

func checkURL(field, s string) error {
	if s == "" {
		return nil
	}
	u, err := url.Parse(s)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http" && !strings.HasPrefix(s, "/uploads/")) {
		return fmt.Errorf("%s must be an http(s) URL", field)
	}
	return nil
}

func checkURLs(field string, list List) error {
	if len(list) > 30 {
		return fmt.Errorf("%s: too many items", field)
	}
	for _, s := range list {
		if err := checkURL(field, s); err != nil {
			return err
		}
	}
	return nil
}

func checkLen(field, s string, max int) error {
	if len(s) > max {
		return fmt.Errorf("%s: too long", field)
	}
	return nil
}

func firstErr(errs ...error) error { return errors.Join(errs...) }

func oneOf(field, v string, allowed ...string) error {
	for _, a := range allowed {
		if v == a {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %s", field, strings.Join(allowed, ", "))
}

type Knowledge struct {
	DocCode   string `json:"doc_code"`
	Category  string `json:"category"`
	Kind      string `json:"kind"`
	Title     I18n   `json:"title"`
	Body      I18n   `json:"body"`
	Source    string `json:"source"`
	Informant string `json:"informant"`
	Reviewer  string `json:"reviewer"`
	Status    string `json:"status"`
	CoverURL  string `json:"cover_url"`
	VideoURL  string `json:"video_url"`
	Sort      int    `json:"sort"`
}

func (k Knowledge) Validate() error {
	if !codeRe.MatchString(k.DocCode) {
		return errors.New("doc_code must be 1-40 letters, digits, - or _")
	}
	return firstErr(
		oneOf("kind", k.Kind, "history", "belief", "ritual", "practice", "community"),
		oneOf("status", k.Status, "draft", "review", "published"),
		k.Title.check("title", 300, true), k.Body.check("body", 60_000, true),
		checkLen("category", k.Category, 100), checkLen("source", k.Source, 2000),
		checkLen("informant", k.Informant, 300), checkLen("reviewer", k.Reviewer, 300),
		checkURL("cover_url", k.CoverURL), checkURL("video_url", k.VideoURL),
	)
}

type Place struct {
	Code         string  `json:"code"`
	Category     string  `json:"category"`
	Name         I18n    `json:"name"`
	Description  I18n    `json:"description"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	OpenHours    string  `json:"open_hours"`
	VisitMinutes int     `json:"visit_minutes"`
	Fee          string  `json:"fee"`
	Contact      string  `json:"contact"`
	Tags         List    `json:"tags"`
	Images       List    `json:"images"`
	Source       string  `json:"source"`
	Published    bool    `json:"published"`
}

func (p Place) Validate() error {
	if !codeRe.MatchString(p.Code) {
		return errors.New("code must be 1-40 letters, digits, - or _")
	}
	if p.Lat < -90 || p.Lat > 90 || p.Lon < -180 || p.Lon > 180 {
		return errors.New("lat/lon out of range")
	}
	if p.VisitMinutes <= 0 || p.VisitMinutes > 600 {
		return errors.New("visit_minutes must be 1-600")
	}
	if len(p.Tags) > 20 {
		return errors.New("too many tags")
	}
	return firstErr(
		oneOf("category", p.Category, "sacred", "culture", "history", "community", "food", "stay"),
		p.Name.check("name", 200, true), p.Description.check("description", 20_000, false),
		checkLen("open_hours", p.OpenHours, 200), checkLen("fee", p.Fee, 200), checkLen("contact", p.Contact, 300),
		checkURLs("images", p.Images), checkLen("source", p.Source, 2000),
	)
}

type Route struct {
	Name        I18n    `json:"name"`
	Description I18n    `json:"description"`
	PlaceIDs    []int64 `json:"place_ids"`
	Minutes     int     `json:"minutes"`
	Published   bool    `json:"published"`
}

func (r Route) Validate() error {
	if len(r.PlaceIDs) < 1 || len(r.PlaceIDs) > 30 {
		return errors.New("route needs 1-30 places")
	}
	if r.Minutes <= 0 || r.Minutes > 24*60 {
		return errors.New("minutes must be 1-1440")
	}
	return firstErr(r.Name.check("name", 200, true), r.Description.check("description", 5000, false))
}

type Event struct {
	Title       I18n      `json:"title"`
	Description I18n      `json:"description"`
	Preparation I18n      `json:"preparation"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	PlaceID     *int64    `json:"place_id"`
	CoverURL    string    `json:"cover_url"`
	Contact     string    `json:"contact"`
	Source      string    `json:"source"`
	Published   bool      `json:"published"`
}

func (e Event) Validate() error {
	if e.StartsAt.IsZero() || e.EndsAt.Before(e.StartsAt) {
		return errors.New("ends_at must be after starts_at")
	}
	return firstErr(
		e.Title.check("title", 300, true), e.Description.check("description", 20_000, false),
		e.Preparation.check("preparation", 20_000, false),
		checkLen("contact", e.Contact, 300), checkLen("source", e.Source, 2000), checkURL("cover_url", e.CoverURL),
	)
}

type Shop struct {
	Name        I18n     `json:"name"`
	Kind        string   `json:"kind"`
	Products    I18n     `json:"products"`
	Description I18n     `json:"description"`
	Lat         *float64 `json:"lat"`
	Lon         *float64 `json:"lon"`
	Contact     string   `json:"contact"`
	Link        string   `json:"link"`
	Images      List     `json:"images"`
	Published   bool     `json:"published"`
}

func (s Shop) Validate() error {
	return firstErr(
		oneOf("kind", s.Kind, "shop", "food", "stay", "product"),
		s.Name.check("name", 200, true), s.Products.check("products", 5000, false),
		s.Description.check("description", 10_000, false), checkLen("contact", s.Contact, 300),
		checkURL("link", s.Link), checkURLs("images", s.Images),
	)
}

type News struct {
	Title       I18n       `json:"title"`
	Body        I18n       `json:"body"`
	CoverURL    string     `json:"cover_url"`
	VideoURL    string     `json:"video_url"`
	IsBanner    bool       `json:"is_banner"`
	PublishedAt *time.Time `json:"published_at"`
}

func (n News) Validate() error {
	return firstErr(n.Title.check("title", 300, true), n.Body.check("body", 60_000, false),
		checkURL("cover_url", n.CoverURL), checkURL("video_url", n.VideoURL))
}

type FAQ struct {
	Category  string `json:"category"`
	Question  I18n   `json:"question"`
	Answer    I18n   `json:"answer"`
	DocCode   string `json:"doc_code"`
	Reviewer  string `json:"reviewer"`
	Published bool   `json:"published"`
}

func (f FAQ) Validate() error {
	return firstErr(f.Question.check("question", 1000, true), f.Answer.check("answer", 10_000, true),
		checkLen("category", f.Category, 100), checkLen("doc_code", f.DocCode, 40), checkLen("reviewer", f.Reviewer, 300))
}
