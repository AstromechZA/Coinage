package transaction

import (
	"time"

	"github.com/shopspring/decimal"
)

// Amount is a value that is tagged with a currency or commodity name
type Amount struct {
	Value     *decimal.Decimal
	Commodity string
}

// TransactionLine is a line from a transaction
type TransactionLine struct {
	// Account is the source or destination this line refers to
	Account string

	// Value of this line
	// Value can be nil when this has not been fully balanced yet
	Value *Amount

	// Price that this line "cost"
	// Value can be nil when this has not been fully balanced yet
	Price *Amount
}

// Transaction represents a single transaction that occurred
type Transaction struct {
	// When is when the transaction occurred
	When time.Time

	// Arbitrary description of this transaction
	Description string

	// Labels is a user-specified key value set
	Labels map[string]string

	// Lines is the list of items in the transaction
	Lines []TransactionLine
}
