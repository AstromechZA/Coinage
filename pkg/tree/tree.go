package tree

import (
	"fmt"

	"github.com/AstromechZA/coinage/pkg/core/multiamount"
	"github.com/AstromechZA/coinage/pkg/core/transaction"
	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
)

type EntryRef struct {
	*entry.Entry
	*transaction.Transaction
}

type AccountTree struct {
	Accounts map[string]*AccountTree
	Entries  []EntryRef

	Totals     multiamount.MultiAmount
	TreeTotals multiamount.MultiAmount
}

func New() *AccountTree {
	return &AccountTree{
		Accounts:   make(map[string]*AccountTree),
		Entries:    make([]EntryRef, 0),
		Totals:     make(multiamount.MultiAmount),
		TreeTotals: make(multiamount.MultiAmount),
	}
}

func (a *AccountTree) Insert(t *transaction.Transaction) error {
	for i, e := range t.Entries {
		if e.Value.Value == nil {
			return fmt.Errorf("line %d has nil value", i+1)
		}
		if err := a.insertLine(t, e, e.Account); err != nil {
			return fmt.Errorf("failed to add Line %d: %s", i+1, err)
		}
	}
	return nil
}

func (a *AccountTree) insertLine(t *transaction.Transaction, entry *entry.Entry, accountParts []string) error {
	if len(accountParts) == 0 {
		a.Totals.Add(entry.Value.Commodity, entry.Value.Value)
		a.Entries = append(a.Entries, EntryRef{Entry: entry, Transaction: t})
	} else {
		a.TreeTotals.Add(entry.Value.Commodity, entry.Value.Value)
		first, remainder := accountParts[0], accountParts[1:]

		_, ok := a.Accounts[first]
		if !ok {
			a.Accounts[first] = New()
		}
		return a.Accounts[first].insertLine(t, entry, remainder)
	}
	return nil
}

func (a *AccountTree) Lookup(prefix []string) (*AccountTree, error) {
	if len(prefix) == 0 {
		return a, nil
	}
	next, ok := a.Accounts[prefix[0]]
	if ok {
		return next.Lookup(prefix[1:])
	}
	return nil, fmt.Errorf("account does not exist")
}

// DepthFirst visits the node and supports multiple iteration types
func (a *AccountTree) DepthFirst(
	prefix []string,
	preOrder func([]string, *AccountTree) bool,
	inOrder func([]string, *AccountTree) bool,
	postOrder func([]string, *AccountTree) bool,
) bool {
	if preOrder != nil && preOrder(prefix, a) {
		return true
	}
	for n, acc := range a.Accounts {
		name := append(prefix, n)
		if acc.DepthFirst(name, preOrder, inOrder, postOrder) {
			return true
		}
		if inOrder != nil && inOrder(prefix, a) {
			return true
		}
	}
	if postOrder != nil && postOrder(prefix, a) {
		return true
	}
	return false
}

func (a *AccountTree) DepthFirstPreOrder(prefix []string, f func([]string, *AccountTree) bool) bool {
	return a.DepthFirst(prefix, f, nil, nil)
}

func (a *AccountTree) DepthFirstInOrder(prefix []string, f func([]string, *AccountTree) bool) bool {
	return a.DepthFirst(prefix, nil, f, nil)
}

func (a *AccountTree) DepthFirstPostOrder(prefix []string, f func([]string, *AccountTree) bool) bool {
	return a.DepthFirst(prefix, nil, nil, f)
}
