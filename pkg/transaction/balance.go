package transaction

import (
	"fmt"
	"math/big"
)

const errorFloatPrecision = 6

// IsFullyBalanced attempts to determine whether the given transaction is balanced
func IsFullyBalanced(transaction *Transaction) (ok bool, err error) {
	commodityValues := make(map[string]*big.Float)

	for _, l := range transaction.Lines {
		_, ok := commodityValues[l.Value.Commodity]
		if !ok {
			commodityValues[l.Value.Commodity] = big.NewFloat(0)
		}
		commodityValues[l.Value.Commodity].Add(commodityValues[l.Value.Commodity], l.Value.Value)
	}

	for commodity, value := range commodityValues {
		if value.Cmp(big.NewFloat(0)) != 0 {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.Text('f', errorFloatPrecision), commodity)
		}
	}

	return true, nil
}

func EnsureBalanced(transaction *Transaction) (changed bool, err error) {
	commodityValues := make(map[string]*big.Float)
	commodityBalanceLines := make(map[string]int)

	for i, l := range transaction.Lines {
		if l.Value.Value == nil {
			if _, ok := commodityBalanceLines[l.Value.Commodity]; ok {
				return false, fmt.Errorf("multiple no-value lines for commodity `%s`", l.Value.Commodity)
			}
			commodityBalanceLines[l.Value.Commodity] = i
		} else {
			if l.Price == nil {
				l.Price = &Amount{Value: &big.Float{}, Commodity: l.Value.Commodity}
				l.Price.Value.Copy(l.Value.Value)
			}
		}

		if l.Price != nil {
			vv, ok := commodityValues[l.Price.Commodity]
			if !ok {
				commodityValues[l.Price.Commodity] = big.NewFloat(0)
				vv = commodityValues[l.Price.Commodity]
			}
			vv.Add(vv, l.Price.Value)
		}
	}

	for commodity, value := range commodityValues {
		bi, ok := commodityBalanceLines[commodity]
		if ok {
			l := transaction.Lines[bi]
			l.Value.Value = new(big.Float).Neg(value)
			l.Price = &Amount{Value: new(big.Float).Copy(l.Value.Value), Commodity: l.Value.Commodity}
			changed = true
			continue
		}
		if value.Cmp(big.NewFloat(0)) != 0 {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.Text('f', errorFloatPrecision), commodity)
		}
	}

	return changed, nil
}
