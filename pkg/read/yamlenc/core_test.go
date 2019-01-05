package yamlenc

import (
	"strings"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestDecodeTransaction(t *testing.T) {
	transaction, err := DecodeTransaction(strings.NewReader(`
when: 2018-01-01
description: I bought something at a shop
labels:
  shop: blah
entries:
- Assets:Cash           -100 £
- Expenses:Groceries    £
`))
	assert.Equal(t, err, nil)
	assert.ShouldEqual(t, len(transaction.Entries), 2)
	assert.ShouldEqual(t, transaction.Entries[1].Value.Value.String(), "100")
}

func TestDecodeTransactions(t *testing.T) {
	transactions, err := DecodeTransactions(strings.NewReader(`
when: 2018-01-01
description: I bought something at a shop
labels:
  shop: blah
entries:
- Assets:Cash           -100 £
- Expenses:Groceries    £

---

when: 2018-01-01
description: I bought something at a shop
labels:
  shop: blah
entries:
- Assets:Cash           -200 £
- Expenses:Groceries    £
`))
	assert.Equal(t, err, nil)
	assert.ShouldEqual(t, len(transactions), 2)
	assert.ShouldEqual(t, len(transactions[0].Entries), 2)
	assert.ShouldEqual(t, transactions[0].Entries[1].Value.Value.String(), "100")
	assert.ShouldEqual(t, transactions[1].Entries[1].Value.Value.String(), "200")
}
