package crest

import (
	"fmt"
	"math/big"
	"time"
)

type CustomFloat struct {
	big.Float
}

type CustomTime struct {
	time.Time
}

const (
	customTimeLayout = "2006-01-02T15:04:05"
)

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	ct.Time, err = time.Parse(customTimeLayout, string(b))
	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Time.Format(customTimeLayout))), nil
}

func (ct *CustomTime) IsSet() bool {
	return ct.UnixNano() != (time.Time{}).UnixNano()
}

func (cf *CustomFloat) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	cf.UnmarshalText(b)
	return
}

func (cf *CustomFloat) MarshalJSON() ([]byte, error) {
	return cf.MarshalText()
}
