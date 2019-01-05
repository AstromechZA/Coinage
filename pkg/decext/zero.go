package decext

import "github.com/ericlagergren/decimal"

func IsZero(d *decimal.Big) bool {
	return d.CmpAbs(new(decimal.Big)) == 0
}
