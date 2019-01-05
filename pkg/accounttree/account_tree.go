package accounttree

import (
	"fmt"
	"strings"

	"github.com/AstromechZA/coinage/pkg/transaction"
)

type EntryRef struct {
	*transaction.Entry
	*transaction.Transaction
}

type AccountTree struct {
	Accounts map[string]*AccountTree
	Entries  []EntryRef
}

func (a *AccountTree) Insert(t *transaction.Transaction) error {
	for i, e := range t.Entries {
		if e.Value.Value == nil {
			return fmt.Errorf("line %d has nil value", i+1)
		}
		if err := a.insertLine(t, e, strings.Split(e.Account, ":")); err != nil {
			return fmt.Errorf("failed to add Line %d: %s", i+1, err)
		}
	}
	return nil
}

func (a *AccountTree) insertLine(t *transaction.Transaction, entry *transaction.Entry, accountParts []string) error {
	if len(accountParts) == 0 {
		a.Entries = append(a.Entries, EntryRef{Entry: entry, Transaction: t})
	} else {
		first, remainder := accountParts[0], accountParts[1:]

		_, ok := a.Accounts[first]
		if !ok {
			a.Accounts[first] = new(AccountTree)
		}
		return a.Accounts[first].insertLine(t, entry, remainder)
	}
	return nil
}

func (a *AccountTree) DepthFirstVisit(prefix []string, f func([]string, *AccountTree) bool) bool {
	for n, acc := range a.Accounts {
		name := append(prefix, n)
		if acc.DepthFirstVisit(name, f) {
			return true
		}
	}
	return f(prefix, a)
}
