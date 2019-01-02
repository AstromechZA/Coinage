package transaction

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// IsFullyBalanced attempts to determine whether the given transaction is balanced
func IsFullyBalanced(transaction *Transaction) (ok bool, err error) {
	commodityValues := make(map[string]decimal.Decimal)

	for _, l := range transaction.Lines {
		vv, ok := commodityValues[l.Value.Commodity]
		if !ok {
			vv = decimal.NewFromFloat(0)
		}
		commodityValues[l.Value.Commodity] = vv.Add(*l.Value.Value)
	}

	for commodity, value := range commodityValues {
		if !value.Equal(decimal.Zero) {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), commodity)
		}
	}

	return true, nil
}

func EnsureBalanced(transaction *Transaction) (changed bool, err error) {
	commodityValues := make(map[string]decimal.Decimal)
	commodityBalanceLines := make(map[string]int)

	for i, l := range transaction.Lines {
		if l.Value.Value == nil {
			if _, ok := commodityBalanceLines[l.Value.Commodity]; ok {
				return false, fmt.Errorf("multiple no-value lines for commodity `%s`", l.Value.Commodity)
			}
			commodityBalanceLines[l.Value.Commodity] = i
		} else {
			if l.Price == nil {
				l.Price = &Amount{Value: l.Value.Value, Commodity: l.Value.Commodity}
			}
		}

		if l.Price != nil {
			vv, ok := commodityValues[l.Price.Commodity]
			if !ok {
				vv = decimal.Zero
			}
			commodityValues[l.Price.Commodity] = vv.Add(*l.Price.Value)
		}
	}

	for commodity, value := range commodityValues {
		bi, ok := commodityBalanceLines[commodity]
		if ok {
			vv := value.Neg()
			l := transaction.Lines[bi]
			l.Value.Value = &vv
			l.Price = &Amount{Value: &vv, Commodity: l.Value.Commodity}
			changed = true
			continue
		}
		if !value.Equal(decimal.Zero) {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), commodity)
		}
	}

	return changed, nil
}
