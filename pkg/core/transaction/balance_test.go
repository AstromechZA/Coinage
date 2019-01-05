package transaction

import (
	"fmt"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
	"github.com/AstromechZA/coinage/pkg/read/line"
)

func TestEnsureBalanced_basic(t *testing.T) {
	transaction := &Transaction{
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("To		GBP"),
		},
	}
	changed, err := transaction.Balance()
	assert.ShouldEqual(t, changed, true)
	assert.ShouldEqual(t, err, nil)
	assert.ShouldEqual(t, transaction.Entries[1].Value.Value.String(), "0.2")
}

func TestEnsureBalanced_multiple_no_change(t *testing.T) {
	transaction := &Transaction{
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("From		0.2 ZAR"),
			line.MustStringToEntry("To		0.2 GBP"),
			line.MustStringToEntry("To		-0.2 ZAR"),
		},
	}
	changed, err := transaction.Balance()
	assert.ShouldEqual(t, changed, false)
	assert.ShouldEqual(t, err, nil)
}

func TestEnsureBalanced_addition_of_values(t *testing.T) {
	transaction := &Transaction{
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("To		GBP"),
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, true)
	assert.Equal(t, err, nil)
	assert.Equal(t, transaction.Entries[2].Value.Value.String(), "0.4")
}

func TestEnsureBalanced_multiple_commodities(t *testing.T) {
	transaction := &Transaction{
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("From		0.2 ZAR"),
			line.MustStringToEntry("To		GBP"),
			line.MustStringToEntry("To		ZAR"),
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
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("To		0.3 GBP"),
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("unbalanced value of 0.1 for GBP"))
}

func TestEnsureBalanced_bad_wildcards(t *testing.T) {
	transaction := &Transaction{
		Entries: []*entry.Entry{
			line.MustStringToEntry("From		-0.2 GBP"),
			line.MustStringToEntry("To		GBP"),
			line.MustStringToEntry("To		GBP"),
		},
	}
	changed, err := transaction.Balance()
	assert.Equal(t, changed, false)
	assert.Equal(t, err, fmt.Errorf("line 3: multiple no-value lines for commodity `GBP`"))
}
