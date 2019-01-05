package decext

import "github.com/ericlagergren/decimal"

func Copy(d *decimal.Big) *decimal.Big {
	return new(decimal.Big).Set(d)
}
