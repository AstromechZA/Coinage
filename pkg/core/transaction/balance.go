package transaction

import (
	"fmt"

	"github.com/ericlagergren/decimal"

	"github.com/AstromechZA/coinage/pkg/core/commodity"

	"github.com/AstromechZA/coinage/pkg/decext"

	"github.com/AstromechZA/coinage/pkg/core/multiamount"
)

func (transaction *Transaction) Balance() (changed bool, err error) {
	commodityValues := multiamount.New()
	commodityBalanceLines := make(map[commodity.Commodity]int)

	for i, l := range transaction.Entries {

		if l.Value.Value == nil {

			if l.Price.Value != nil {
				return false, fmt.Errorf("line %d: no-value line cannot have a price", i+1)
			}
			if l.Price.Commodity != "" && l.Price.Commodity != l.Value.Commodity {
				return false, fmt.Errorf("line %d: no-value commidity does not match no-value price commodity", i+1)
			}

			if _, ok := commodityBalanceLines[l.Value.Commodity]; ok {
				return false, fmt.Errorf("line %d: multiple no-value lines for commodity `%s`", i+1, l.Value.Commodity)
			}
			commodityBalanceLines[l.Value.Commodity] = i
			continue
		}

		if l.Price.Value == nil {
			l.Price.Value = new(decimal.Big).Abs(l.Value.Value)
			l.Price.Commodity = l.Value.Commodity
		}

		if l.Value.Value.Signbit() {
			commodityValues.Sub(l.Price.Commodity, l.Price.Value)
		} else {
			commodityValues.Add(l.Price.Commodity, l.Price.Value)
		}
	}

	for c, value := range *commodityValues {
		bi, ok := commodityBalanceLines[c]
		if ok {
			l := transaction.Entries[bi]
			l.Value.Value = value.Neg(value)

			// if the price of the balance line was also unsure, then copy it over
			if l.Price.Value == nil {
				l.Price.Value = new(decimal.Big).Abs(l.Value.Value)
				l.Price.Commodity = l.Value.Commodity
			}
			changed = true
			continue
		}
		if !decext.IsZero(value) {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), c)
		}
	}

	return changed, nil
}
