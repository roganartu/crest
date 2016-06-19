package crest

import (
	"fmt"
)

const (
	CREST_URL = "https://crest-tq.eveonline.com"
)

func Init() error {
	types, err := LoadMarketTypes()
	if err != nil {
		return err
	}

	forge := 10000002
	for _, t := range types {
		fmt.Println(t.Name)
		orders, err := LoadOrders(forge, t.ID)
		if err != nil {
			return err
		}
		fmt.Printf("\t%d orders\n", len(orders))
	}

	return nil
}
