package transaction

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

// IsFullyBalanced attempts to determine whether the given transaction is balanced
func IsFullyBalanced(transaction *Transaction) (ok bool, err error) {
	commodityValues := make(map[string]*decimal.Big)

	for _, l := range transaction.Lines {
		vv, ok := commodityValues[l.Value.Commodity]
		if !ok {
			vv = new(decimal.Big)
		}
		commodityValues[l.Value.Commodity] = vv.Add(vv, l.Value.Value)
	}

	for commodity, value := range commodityValues {
		if value.Cmp(new(decimal.Big)) != 0 {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), commodity)
		}
	}

	return true, nil
}

func EnsureBalanced(transaction *Transaction) (changed bool, err error) {
	commodityValues := make(map[string]*decimal.Big)
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
				vv = new(decimal.Big)
			}
			commodityValues[l.Price.Commodity] = vv.Add(vv, l.Price.Value)
		}
	}

	for commodity, value := range commodityValues {
		bi, ok := commodityBalanceLines[commodity]
		if ok {
			l := transaction.Lines[bi]
			l.Value.Value = value.Neg(value)
			l.Price = &Amount{Value: new(decimal.Big).Copy(l.Value.Value), Commodity: l.Value.Commodity}
			changed = true
			continue
		}
		if value.Cmp(new(decimal.Big)) != 0 {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.String(), commodity)
		}
	}

	return changed, nil
}
