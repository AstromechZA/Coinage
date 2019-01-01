package transaction

import (
	"fmt"
	"math/big"
)

// IsBalanced attempts to determine whether the given transaction is balanced
func IsBalanced(transaction *Transaction) (ok bool, err error) {
	commodityValues := make(map[string]*big.Float)

	for i, t := range transaction.Lines {
		if t.Value == nil {
			return false, fmt.Errorf("line %d does not contain a value", i+1)
		}

		_, ok := commodityValues[t.Value.Commodity]
		if !ok {
			commodityValues[t.Value.Commodity] = big.NewFloat(0)
		}
		commodityValues[t.Value.Commodity].Add(commodityValues[t.Value.Commodity], t.Value.Value)
	}

	for commodity, value := range commodityValues {
		if value.Cmp(big.NewFloat(0)) != 0 {
			return false, fmt.Errorf("unbalanced value of %s for %s", value.Text('f', -1), commodity)
		}
	}

	return true, nil
}
