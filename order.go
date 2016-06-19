package crest

import (
	"encoding/json"
	"fmt"
	"math/big"
)

type orderWrapper struct {
	Count int     `json:"totalCount"`
	Pages int     `json:"pageCount"`
	Items []Order `json:"items"`
}

type Order struct {
	Volume   big.Int     `json:"volume"`
	Buy      bool        `json:"buy"`
	Price    CustomFloat `json:"price"`
	Date     CustomTime  `json:"issued"`
	Location Location    `json:"location"`
}

const (
	ORDERS_URL = "https://crest-tq.eveonline.com/market/%d/orders/?type=https://crest-tq.eveonline.com/inventory/types/%d/"
)

func LoadOrders(region, itemType int) ([]Order, error) {
	url := fmt.Sprintf(ORDERS_URL, region, itemType)

	r, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	var wrap orderWrapper
	err = json.NewDecoder(r.Body).Decode(&wrap)
	if err != nil {
		return nil, err
	}

	return wrap.Items, nil
}
