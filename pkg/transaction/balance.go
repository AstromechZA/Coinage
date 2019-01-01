package transaction

import (
	"fmt"
	"math/big"
)

// IsBalanced attempts to determine whether the given transaction is balanced
func IsBalanced(transaction *Transaction) (ok bool, err error) {
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
			return false, fmt.Errorf("unbalanced value of %s for %s", value.Text('f', -1), commodity)
		}
	}

	return true, nil
}
