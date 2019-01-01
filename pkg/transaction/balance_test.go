package transaction

import (
	"math/big"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestIsBalanced_when_empty(t *testing.T) {
	balanced, err := IsBalanced(&Transaction{})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_when_exact_match(t *testing.T) {
	balanced, err := IsBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: big.NewFloat(-0.333), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: big.NewFloat(0.333), Commodity: "GBP"}},
		},
	})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_two_from(t *testing.T) {
	balanced, err := IsBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: big.NewFloat(-0.333), Commodity: "GBP"}},
			{Account: "From", Value: &Amount{Value: big.NewFloat(-0.333), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: big.NewFloat(0.666), Commodity: "GBP"}},
		},
	})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_simple_not_balanced(t *testing.T) {
	balanced, err := IsBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: big.NewFloat(-0.2), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: big.NewFloat(0.1), Commodity: "GBP"}},
		},
	})
	if assert.ShouldBeFalse(t, balanced) {
		assert.Equal(t, err.Error(), "unbalanced value of -0.1 for GBP")
	}
}
