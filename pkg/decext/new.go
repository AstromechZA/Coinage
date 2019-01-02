package decext

import (
	"fmt"

	"github.com/shopspring/decimal"
)

func NewPtr(value int64, exp int32) *decimal.Decimal {
	v := decimal.New(value, exp)
	return &v
}

func MustNewFromString(value string) *decimal.Decimal {
	v, err := decimal.NewFromString(value)
	if err != nil {
		panic(fmt.Errorf("could not convert '%s' to decimal: %s", value, err))
	}
	return &v
}
