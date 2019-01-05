package decext

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

func NewFromString(value string) (*decimal.Big, error) {
	d, ok := new(decimal.Big).SetString(value)
	if !ok {
		return nil, fmt.Errorf("`%s` is an invalid decimal", value)
	}
	if !d.IsFinite() {
		return nil, fmt.Errorf("`%s` is not supported", value)
	}
	return d, nil
}

func MustNewFromString(value string) *decimal.Big {
	v, err := NewFromString(value)
	if err != nil {
		panic(err)
	}
	return v
}
