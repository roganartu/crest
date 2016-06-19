package crest

import (
	"encoding/json"
	"fmt"
	"os"
)

type marketGroupWrapper struct {
	Count int           `json:"totalCount"`
	Pages int           `json:"pageCount"`
	Items []MarketGroup `json:"items"`
}

type MarketGroup struct {
	ID          int               `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Path        string            `json:"href"`
	TypesPath   map[string]string `json:"types"`
	Types       []Type            `json:"_types"`
}

func loadMarketGroups() ([]MarketGroup, error) {
	if len(ENDPOINTS) < 1 {
		err := loadEndpoints()
		if err != nil {
			return nil, err
		}
	}

	// Try and load groups from cache
	f, err := os.Open("market_groups.json")
	if err == nil {
		groups := make([]MarketGroup, 0)
		err = json.NewDecoder(f).Decode(&groups)
		f.Close()
		if err == nil {
			return groups, nil
		}
	}

	r, err := client.Get(ENDPOINTS["marketGroups"].Path)
	if err != nil {
		return nil, err
	}

	var wrap marketGroupWrapper
	err = json.NewDecoder(r.Body).Decode(&wrap)
	if err != nil {
		return nil, err
	}

	for _, group := range wrap.Items {
		fmt.Println(group.Name)
		// Populate types
		group.Types, err = loadMarketTypesForGroup(group.TypesPath["href"])
		if err != nil {
			return nil, err
		}
	}

	// Cache the groups for later use
	f, err = os.Create("market_groups.json")
	if err == nil {
		err = json.NewEncoder(f).Encode(wrap.Items)
		f.Close()
	}

	return wrap.Items, nil
}
