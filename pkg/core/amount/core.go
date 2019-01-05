package amount

import (
	"fmt"

	"github.com/AstromechZA/coinage/pkg/decext"

	"github.com/AstromechZA/coinage/pkg/core/commodity"
	"github.com/ericlagergren/decimal"
)

// Amount is a value that is tagged with a currency or commodity name
type Amount struct {
	Value     *decimal.Big
	Commodity commodity.Commodity
}

func New(c commodity.Commodity, v *decimal.Big) *Amount {
	if v == nil || !v.IsFinite() {
		panic(fmt.Errorf("refusing to create non finite amount %s", v))
	}
	return &Amount{
		Commodity: c,
		Value:     v,
	}
}

func NewZero(c commodity.Commodity) *Amount {
	return &Amount{
		Commodity: c,
		Value:     new(decimal.Big),
	}
}

func NewNil(c commodity.Commodity) *Amount {
	return &Amount{
		Commodity: c,
		Value:     nil,
	}
}

func (a Amount) Empty() bool {
	return a.Value == nil
}

func (a Amount) Copy() *Amount {
	return New(a.Commodity, decext.Copy(a.Value))
}
