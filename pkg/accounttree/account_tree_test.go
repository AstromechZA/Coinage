package accounttree

import (
	"testing"
	"time"

	"github.com/AstromechZA/coinage/pkg/amount"
	"github.com/AstromechZA/coinage/pkg/assert"
	"github.com/AstromechZA/coinage/pkg/transaction"
	"github.com/ericlagergren/decimal"
)

func TestAccountTree_Insert(t *testing.T) {
	at := New()
	trans := &transaction.Transaction{
		When:        time.Time{},
		Description: "Some description",
		Entries: []*transaction.Entry{
			{Account: "A:X", Value: *amount.New("USD", decimal.New(100, 0))},
			{Account: "B:Y", Value: *amount.New("USD", decimal.New(-100, 0))},
		},
	}
	err := at.Insert(trans)
	assert.Equal(t, err, nil)
	assert.ShouldEqual(t, at.TreeTotals["USD"].String(), "0")
	assert.ShouldEqual(t, at.Accounts["A"].TreeTotals["USD"].String(), "100")
	assert.ShouldEqual(t, at.Accounts["B"].TreeTotals["USD"].String(), "-100")
}
