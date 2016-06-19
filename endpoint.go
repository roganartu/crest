package crest

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ENDPOINTS = make(map[string]Endpoint)
)

type Endpoint struct {
	Path string
}

func loadEndpoints() error {
	r, err := client.Get(CREST_URL)
	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return errors.New("Failed to load endpoints.")
	}

	tmp := make(map[string]interface{})
	err = json.NewDecoder(r.Body).Decode(&tmp)
	if err != nil {
		return err
	}

	// This is messy. Thanks EVE for having a gross API
	for k, v := range tmp {
		if e, ok := v.(map[string]interface{}); ok {
			for _, val := range e {
				if x, ok := val.(map[string]interface{}); ok {
					for a, b := range x {
						path, _ := b.(string)
						ENDPOINTS[k+"."+a] = Endpoint{
							Path: path,
						}
					}
				} else {
					path, _ := e["href"].(string)
					ENDPOINTS[k] = Endpoint{
						Path: path,
					}
				}
			}
		}
	}

	return nil
}
