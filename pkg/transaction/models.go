package transaction

import (
	"time"

	"github.com/AstromechZA/coinage/pkg/amount"
)

// Entry is a Line from a transaction
type Entry struct {
	// Account is the source or destination this Line refers to
	Account string

	// Value of this Line in the account
	// Value.Value can be nil when this has not been fully balanced yet
	Value amount.Amount

	// Price that this Line "cost" to balance (used to translate between currencies)
	// Price can be nil when this has not been fully balanced yet
	Price amount.Amount
}

// Transaction represents a single transaction that occurred
type Transaction struct {
	// When is when the transaction occurred
	When time.Time

	// Arbitrary description of this transaction
	Description string

	// Labels is a user-specified key value set
	Labels map[string]string

	// Entries is the list of items in the transaction
	Entries []*Entry
}
