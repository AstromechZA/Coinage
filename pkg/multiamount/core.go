package multiamount

import (
	"errors"
	"fmt"

	"github.com/ericlagergren/decimal"

	"github.com/AstromechZA/coinage/pkg/amount"
	"github.com/AstromechZA/coinage/pkg/commodity"
)

// MultiAmount is a value that is spread over multiple commodities
type MultiAmount map[commodity.Commodity]*decimal.Big

func New() *MultiAmount {
	v := make(MultiAmount)
	return &v
}

func (ma *MultiAmount) Has(c commodity.Commodity) bool {
	_, ok := (*ma)[c]
	return ok
}

func (ma *MultiAmount) Ensure(c commodity.Commodity) *decimal.Big {
	x, ok := (*ma)[c]
	if !ok {
		x = new(decimal.Big)
		(*ma)[c] = x
	}
	return x
}

func (ma *MultiAmount) AddAmount(amount *amount.Amount) *decimal.Big {
	if amount == nil {
		panic(errors.New("refusing to add nil amount"))
	}
	return ma.Add(amount.Commodity, amount.Value)
}

func (ma *MultiAmount) Add(c commodity.Commodity, v *decimal.Big) *decimal.Big {
	if v == nil || !v.IsFinite() {
		panic(fmt.Errorf("refusing to add non finite value %s", v))
	}
	current := ma.Ensure(c)
	return current.Add(current, v)
}

func (ma *MultiAmount) Sub(c commodity.Commodity, v *decimal.Big) *decimal.Big {
	if v == nil || !v.IsFinite() {
		panic(fmt.Errorf("refusing to subtract non finite value %s", v))
	}
	current := ma.Ensure(c)
	return current.Sub(current, v)
}
