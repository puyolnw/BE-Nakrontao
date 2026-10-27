package content

import (
	"encoding/json"
	"testing"
)

func ptr[T any](v T) *T { return &v }

func validShop() Shop {
	return Shop{Name: I18n{"th": "ผ้าไหม"}, Kind: "product"}
}

func TestShopValidate(t *testing.T) {
	cases := []struct {
		name string
		edit func(*Shop)
		ok   bool
	}{
		{"minimal", func(*Shop) {}, true},
		{"price range", func(s *Shop) { s.PriceFrom, s.PriceTo = ptr(350), ptr(2500) }, true},
		{"price reversed", func(s *Shop) { s.PriceFrom, s.PriceTo = ptr(500), ptr(100) }, false},
		{"negative price", func(s *Shop) { s.PriceFrom = ptr(-1) }, false},
		{"seller ok", func(s *Shop) {
			s.Sellers = Sellers{{Name: "ครัวเรือนที่ 1", Village: "1", Phone: "0800000000", Lat: ptr(15.51), Lon: ptr(103.38)}}
		}, true},
		{"seller without name", func(s *Shop) { s.Sellers = Sellers{{Phone: "0800000000"}} }, false},
		{"seller lat without lon", func(s *Shop) { s.Sellers = Sellers{{Name: "x", Lat: ptr(15.5)}} }, false},
		{"seller bad lat", func(s *Shop) { s.Sellers = Sellers{{Name: "x", Lat: ptr(95.0), Lon: ptr(103.0)}} }, false},
		{"bad unit language", func(s *Shop) { s.PriceUnit = I18n{"fr": "pièce"} }, false},
		{"too many sellers", func(s *Shop) {
			s.Sellers = make(Sellers, 51)
			for i := range s.Sellers {
				s.Sellers[i].Name = "x"
			}
		}, false},
	}
	for _, c := range cases {
		s := validShop()
		c.edit(&s)
		if err := s.Validate(); (err == nil) != c.ok {
			t.Errorf("%s: ok=%v, err=%v", c.name, c.ok, err)
		}
	}
}

func TestSellersMarshalNil(t *testing.T) {
	b, err := json.Marshal(validShop())
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	if string(m["sellers"]) != "[]" {
		t.Fatalf("sellers = %s, want []", m["sellers"])
	}
}
