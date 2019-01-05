package transaction

import (
	"fmt"

	"github.com/AstromechZA/coinage/pkg/core/commodity"

	"github.com/AstromechZA/coinage/pkg/decext"

	"github.com/AstromechZA/coinage/pkg/core/multiamount"
)

func (transaction *Transaction) Balance() (changed bool, err error) {
	commodityValues := multiamount.New()
	commodityBalanceLines := make(map[commodity.Commodity]int)

	for i, l := range transaction.Entries {

		if l.Value.Value == nil {
			if _, ok := commodityBalanceLines[l.Value.Commodity]; ok {
				return false, fmt.Errorf("line %d: multiple no-value lines for commodity `%s`", i+1, l.Value.Commodity)
			}
			commodityBalanceLines[l.Value.Commodity] = i
			continue
		}
		commodityValues.Add(l.Value.Commodity, l.Value.Value)
	}

	for c, value := range *commodityValues {
		bi, ok := commodityBalanceLines[c]
		if ok {
			l := transaction.Entries[bi]
			l.Value.Value = value.Neg(value)
			changed = true
			continue
		}
		if !decext.IsZero(value) {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), c)
		}
	}

	return changed, nil
}
