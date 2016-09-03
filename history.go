package crest

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type historyWrapper struct {
	Count int        `json:"totalCount"`
	Pages int        `json:"pageCount"`
	Items []*History `json:"items"`
}

type History struct {
	Volume big.Int     `json:"volume"`
	Count  big.Int     `json:"orderCount"`
	Low    CustomFloat `json:"lowPrice"`
	High   CustomFloat `json:"highPrice"`
	Avg    CustomFloat `json:"avgPrice"`
	Date   CustomTime  `json:"date"`
}

const (
	HISTORY_URL = "https://crest-tq.eveonline.com/market/%d/history/?type=https://crest-tq.eveonline.com/inventory/types/%d/"
)

func LoadHistory(region, itemType int) ([]*History, error) {
	url := fmt.Sprintf(HISTORY_URL, region, itemType)

	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	var wrap historyWrapper
	err = json.NewDecoder(r.Body).Decode(&wrap)
	if err != nil {
		return nil, err
	}

	return wrap.Items, nil
}
