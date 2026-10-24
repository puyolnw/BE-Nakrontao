// Package recommend builds a visit plan from interests and a time budget (scope 6.x).
// Rule-based + content-based, as agreed in scope: no trained model, no live traffic.
package recommend

import (
	"encoding/json"
	"math"
	"net/http"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/reclover/sinakhontao/internal/httpx"
)

type Place struct {
	ID           int64             `json:"id"`
	Category     string            `json:"category"`
	Name         map[string]string `json:"name"`
	Lat          float64           `json:"lat"`
	Lon          float64           `json:"lon"`
	VisitMinutes int               `json:"visit_minutes"`
	Tags         []string          `json:"tags"`
}

type Request struct {
	Interests []string `json:"interests"` // categories or tags, e.g. sacred, history, food
	Minutes   int      `json:"minutes"`
	StartLat  *float64 `json:"start_lat"`
	StartLon  *float64 `json:"start_lon"`
	Mode      string   `json:"mode"` // walk | drive
}

type Stop struct {
	Place       Place    `json:"place"`
	TravelMin   int      `json:"travel_minutes"`
	ArriveAtMin int      `json:"arrive_at_minute"`
	Reasons     []string `json:"reasons"` // machine codes, translated by the client (scope 6.7)
	Score       float64  `json:"score"`
}

type Plan struct {
	Stops        []Stop `json:"stops"`
	TotalMinutes int    `json:"total_minutes"`
}

const anchorCategory = "sacred" // the shrine is the reason people come; it always gets a bonus

func score(p Place, interests map[string]bool) (float64, []string) {
	var s float64
	var reasons []string
	if interests[p.Category] {
		s += 3
		reasons = append(reasons, "match_category:"+p.Category)
	}
	for _, t := range p.Tags {
		if interests[t] {
			s += 1
			reasons = append(reasons, "match_tag:"+t)
		}
	}
	if p.Category == anchorCategory {
		s += 1.5
		reasons = append(reasons, "key_sacred_site")
	}
	if p.VisitMinutes <= 30 {
		s += 0.3
		reasons = append(reasons, "short_visit")
	}
	return s, reasons
}

func matches(p Place, interests map[string]bool) bool {
	if interests[p.Category] {
		return true
	}
	for _, t := range p.Tags {
		if interests[t] {
			return true
		}
	}
	return false
}

// Haversine distance in km.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const r = 6371.0
	rad := math.Pi / 180
	dLat, dLon := (lat2-lat1)*rad, (lon2-lon1)*rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * r * math.Asin(math.Sqrt(a))
}

// travelMinutes estimates door-to-door time with a 1.3 detour factor on straight-line distance.
func travelMinutes(km float64, mode string) int {
	speed := 4.5 // km/h walking
	if mode == "drive" {
		speed = 35
	}
	return int(math.Ceil(km * 1.3 / speed * 60))
}

// Build greedily picks the next stop maximising score per minute spent (travel + visit), within budget.
func Build(places []Place, req Request) Plan {
	interests := map[string]bool{}
	for _, i := range req.Interests {
		interests[i] = true
	}
	type cand struct {
		p       Place
		score   float64
		reasons []string
	}
	var cands []cand
	for _, p := range places {
		s, rs := score(p, interests)
		// Only places matching an interest qualify; the shrine always qualifies as the anchor site.
		if matches(p, interests) || p.Category == anchorCategory {
			cands = append(cands, cand{p, s, rs})
		}
	}
	sort.SliceStable(cands, func(i, j int) bool { return cands[i].score > cands[j].score })

	var plan Plan
	if len(cands) == 0 {
		plan.Stops = []Stop{}
		return plan
	}
	var curLat, curLon float64
	hasStart := req.StartLat != nil && req.StartLon != nil
	if hasStart {
		curLat, curLon = *req.StartLat, *req.StartLon
	}
	used := make([]bool, len(cands))
	elapsed := 0
	for {
		best, bestVal, bestTravel := -1, -1.0, 0
		for i, c := range cands {
			if used[i] {
				continue
			}
			travel := 0
			if hasStart || len(plan.Stops) > 0 {
				travel = travelMinutes(Haversine(curLat, curLon, c.p.Lat, c.p.Lon), req.Mode)
			}
			if elapsed+travel+c.p.VisitMinutes > req.Minutes {
				continue
			}
			val := c.score / float64(travel+c.p.VisitMinutes)
			if val > bestVal {
				best, bestVal, bestTravel = i, val, travel
			}
		}
		if best < 0 {
			break
		}
		used[best] = true
		c := cands[best]
		elapsed += bestTravel
		reasons := c.reasons
		if bestTravel > 0 && bestTravel <= 10 {
			reasons = append(append([]string(nil), reasons...), "nearby")
		}
		plan.Stops = append(plan.Stops, Stop{Place: c.p, TravelMin: bestTravel, ArriveAtMin: elapsed, Reasons: reasons, Score: c.score})
		elapsed += c.p.VisitMinutes
		curLat, curLon = c.p.Lat, c.p.Lon
	}
	if plan.Stops == nil {
		plan.Stops = []Stop{}
	}
	plan.TotalMinutes = elapsed
	return plan
}

type Handler struct{ DB *pgxpool.Pool }

func (h *Handler) Plan(w http.ResponseWriter, r *http.Request) error {
	var req Request
	if err := httpx.Decode(w, r, &req); err != nil {
		return err
	}
	if req.Minutes < 15 || req.Minutes > 24*60 {
		return httpx.BadRequest("minutes must be 15-1440")
	}
	if len(req.Interests) == 0 || len(req.Interests) > 20 {
		return httpx.BadRequest("choose 1-20 interests")
	}
	if req.Mode != "drive" {
		req.Mode = "walk"
	}
	var raw []byte
	err := h.DB.QueryRow(r.Context(), `SELECT coalesce(jsonb_agg(jsonb_build_object(
		'id', id, 'category', category, 'name', name, 'lat', lat, 'lon', lon, 'visit_minutes', visit_minutes, 'tags', tags)), '[]')
		FROM places WHERE published`).Scan(&raw)
	if err != nil {
		return err
	}
	var places []Place
	if err := json.Unmarshal(raw, &places); err != nil {
		return err
	}
	return httpx.JSON(w, http.StatusOK, Build(places, req))
}
