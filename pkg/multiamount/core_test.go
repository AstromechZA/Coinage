package multiamount

import (
	"testing"

	"github.com/ericlagergren/decimal"

	"github.com/AstromechZA/coinage/pkg/amount"
	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestNew(t *testing.T) {
	ma := New()
	assert.Equal(t, len(*ma), 0)
}

func TestMultiAmount_Ensure(t *testing.T) {
	ma := New()
	assert.Equal(t, len(*ma), 0)
	x := ma.Ensure("USD")
	assert.Equal(t, len(*ma), 1)
	assert.Equal(t, x.Cmp(new(decimal.Big)), 0)
	x.Add(x, decimal.New(100, 0))

	y := ma.Ensure("USD")
	assert.Equal(t, len(*ma), 1)
	assert.Equal(t, y.Cmp(new(decimal.Big)), 1)
	assert.Equal(t, y.String(), "100")
}

func TestMultiAmount_Has(t *testing.T) {
	ma := New()
	assert.False(t, ma.Has("USD"))
	ma.Ensure("USD")
	assert.True(t, ma.Has("USD"))
}

func TestMultiAmount_Add(t *testing.T) {
	ma := New()
	ma.Add("USD", decimal.New(42, 0))
	assert.True(t, ma.Has("USD"))
	x := ma.Add("USD", decimal.New(42, 3))
	assert.Equal(t, x.String(), "42.042")
}

func TestMultiAmount_Add_nil(t *testing.T) {
	defer func() {
		x := recover()
		assert.Equal(t, x.(error).Error(), "refusing to add non finite value <nil>")
	}()
	ma := New()
	ma.Add("USD", nil)
}

func TestMultiAmount_Add_non_finite(t *testing.T) {
	defer func() {
		x := recover()
		assert.Equal(t, x.(error).Error(), "refusing to add non finite value -Infinity")
	}()
	ma := New()
	ma.Add("USD", new(decimal.Big).SetInf(true))
}

func TestMultiAmount_AddAmount(t *testing.T) {
	ma := New()
	ma.AddAmount(&amount.Amount{
		Value:     decimal.New(11, 1),
		Commodity: "USD",
	})
	assert.True(t, ma.Has("USD"))
	assert.Equal(t, (*ma)["USD"].String(), "1.1")
}

func TestMultiAmount_AddAmount_nil(t *testing.T) {
	defer func() {
		x := recover()
		assert.Equal(t, x.(error).Error(), "refusing to add nil amount")
	}()
	ma := New()
	ma.AddAmount(nil)
}
