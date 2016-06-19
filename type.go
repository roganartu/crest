package crest

import (
	"encoding/json"
	"net/url"
	"os"
	"strconv"
)

type typeWrapper struct {
	Count int                      `json:"totalCount"`
	Pages int                      `json:"pageCount"`
	Items []marketGroupTypeWrapper `json:"items"`
}

// Thanks EVE for having the most confusingly organised data structures around.
// Sure, they're not hard to understand but boy was whoever designed the original
// ones confused about implementation. What's a MarketGroup? Oh, could be a few
// things really, depends on context (why?!)
type marketGroupTypeWrapper struct {
	Type Type `json:"type"`
}

type Type struct {
	ID   int    `json:"id"`
	Path string `json:"href"`
	Name string `json:"name"`
}

func loadMarketTypesForGroup(path string) ([]Type, error) {
	return nil, nil
}

func LoadMarketTypes() ([]Type, error) {
	types := make([]Type, 0)

	// Try and load types from cache
	f, err := os.Open("market_types.json")
	if err == nil {
		types := make([]Type, 0)
		err = json.NewDecoder(f).Decode(&types)
		f.Close()
		if err == nil {
			return types, nil
		}
	}

	if len(ENDPOINTS) < 1 {
		err := loadEndpoints()
		if err != nil {
			return nil, err
		}
	}

	u, err := url.Parse(ENDPOINTS["marketTypes"].Path)
	if err != nil {
		return nil, err
	}

	page := 1
	for {
		q := u.Query()
		q.Del("page")
		q.Add("page", strconv.Itoa(page))
		u.RawQuery = q.Encode()
		r, err := client.Get(u.String())
		if err != nil {
			return nil, err
		}

		var wrap typeWrapper
		err = json.NewDecoder(r.Body).Decode(&wrap)
		if err != nil {
			return nil, err
		}

		for _, v := range wrap.Items {
			types = append(types, v.Type)
		}

		page += 1
		if page > wrap.Pages {
			break
		}
	}

	// Cache the types for later use
	f, err = os.Create("market_types.json")
	if err == nil {
		err = json.NewEncoder(f).Encode(types)
		f.Close()
	}

	return types, nil
}
