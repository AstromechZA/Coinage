package entry

import "github.com/AstromechZA/coinage/pkg/core/amount"

// Entry is a Line from a transaction
type Entry struct {
	// Account is the source or destination this Line refers to
	Account []string

	// Value of this Line in the account
	// Value.Value can be nil when this has not been fully balanced yet
	Value amount.Amount

	// Price that this Line "cost" to balance (used to translate between currencies)
	// Price can be nil when this has not been fully balanced yet
	Price amount.Amount
}
