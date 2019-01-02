package transaction

import (
	"fmt"
	"testing"

	"github.com/AstromechZA/coinage/pkg/decext"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestIsBalanced_when_empty(t *testing.T) {
	balanced, err := IsFullyBalanced(&Transaction{})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_when_exact_match(t *testing.T) {
	balanced, err := IsFullyBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.333"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("0.333"), Commodity: "GBP"}},
		},
	})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_two_from(t *testing.T) {
	balanced, err := IsFullyBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.333"), Commodity: "GBP"}},
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.333"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("0.666"), Commodity: "GBP"}},
		},
	})
	if assert.ShouldEqual(t, err, nil) {
		assert.True(t, balanced)
	}
}

func TestIsBalanced_simple_not_balanced(t *testing.T) {
	balanced, err := IsFullyBalanced(&Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("0.1"), Commodity: "GBP"}},
		},
	})
	if assert.ShouldBeFalse(t, balanced) {
		assert.Equal(t, err.Error(), "unbalanced value of -0.1 for GBP")
	}
}

func TestEnsureBalanced_basic(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: nil, Commodity: "GBP"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Lines[1].Value.Value.String(), "0.2")
}

func TestEnsureBalanced_multiple_no_change(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("0.2"), Commodity: "ZAR"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "ZAR"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, false)
	assert.Equal(t, err, nil)
}

func TestEnsureBalanced_addition_of_values(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: nil, Commodity: "GBP"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Lines[2].Value.Value.String(), "0.4")
}

func TestEnsureBalanced_multiple_commodities(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("0.2"), Commodity: "ZAR"}},
			{Account: "To", Value: &Amount{Value: nil, Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: nil, Commodity: "ZAR"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Lines[2].Value.Value.String(), "0.2")
	assert.Equal(t, transaction.Lines[3].Value.Value.String(), "-0.2")
}

func TestEnsureBalanced_bad_balance(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: decext.MustNewFromString("0.3"), Commodity: "GBP"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("unbalanced value of 0.1 for GBP"))
}

func TestEnsureBalanced_bad_wildcards(t *testing.T) {
	transaction := &Transaction{
		Lines: []TransactionLine{
			{Account: "From", Value: &Amount{Value: decext.MustNewFromString("-0.2"), Commodity: "GBP"}},
			{Account: "To", Value: &Amount{Value: nil, Commodity: "GBP"}},
			{Account: "To2", Value: &Amount{Value: nil, Commodity: "GBP"}},
		},
	}
	changed, err := EnsureBalanced(transaction)
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("multiple no-value lines for commodity `GBP`"))
}
