package transaction

import (
	"fmt"
	"testing"

	"github.com/ericlagergren/decimal"

	"github.com/AstromechZA/coinage/pkg/amount"
	"github.com/AstromechZA/coinage/pkg/assert"
	"github.com/AstromechZA/coinage/pkg/decext"
)

func TestEnsureBalanced_basic(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "To", Value: *amount.NewNil("GBP")},
		},
	}
	changed, err := transaction.Balance()
	assert.ShouldEqual(t, changed, true)
	assert.ShouldEqual(t, err, nil)
	assert.ShouldEqual(t, transaction.Entries[1].Value.Value.String(), "0.2")
}

func TestEnsureBalanced_multiple_no_change(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "From", Value: *amount.New("ZAR", decext.MustNewFromString("0.2"))},
			{Account: "To", Value: *amount.New("GBP", decext.MustNewFromString("0.2"))},
			{Account: "To", Value: *amount.New("ZAR", decext.MustNewFromString("-0.2"))},
		},
	}
	changed, err := transaction.Balance()
	assert.ShouldEqual(t, changed, false)
	assert.ShouldEqual(t, err, nil)
}

func TestEnsureBalanced_addition_of_values(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "To", Value: *amount.NewNil("GBP")},
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Entries[2].Value.Value.String(), "0.4")
}

func TestEnsureBalanced_multiple_commodities(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "From", Value: *amount.New("ZAR", decext.MustNewFromString("0.2"))},
			{Account: "To", Value: *amount.NewNil("GBP")},
			{Account: "To", Value: *amount.NewNil("ZAR")},
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Entries[2].Value.Value.String(), "0.2")
	assert.Equal(t, transaction.Entries[3].Value.Value.String(), "-0.2")
}

func TestEnsureBalanced_bad_balance(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "To", Value: *amount.New("GBP", decext.MustNewFromString("0.3"))},
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("unbalanced value of 0.1 for GBP"))
}

func TestEnsureBalanced_bad_wildcards(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "From", Value: *amount.New("GBP", decext.MustNewFromString("-0.2"))},
			{Account: "To", Value: *amount.NewNil("GBP")},
			{Account: "To2", Value: *amount.NewNil("GBP")},
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("line 3: multiple no-value lines for commodity `GBP`"))
}

func TestEnsureBalanced_stock_example(t *testing.T) {
	transaction := &Transaction{
		Entries: []*Entry{
			{Account: "Stock", Value: *amount.New("ORCL", decext.MustNewFromString("100")), Price: *amount.New("$", decimal.New(50, 0))},
			{Account: "Stock", Value: *amount.New("AAPL", decext.MustNewFromString("100")), Price: *amount.New("$", decimal.New(30, 0))},
			{Account: "Cash", Value: *amount.NewNil("$")},
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, err, nil)
	assert.ShouldEqual(t, changed, true)
	assert.ShouldEqual(t, transaction.Entries[2].Value.Value.String(), "-80")
	assert.ShouldEqual(t, string(transaction.Entries[2].Value.Commodity), "$")
	assert.ShouldEqual(t, transaction.Entries[2].Price.Value.String(), "80")
	assert.ShouldEqual(t, string(transaction.Entries[2].Price.Commodity), "$")
}
