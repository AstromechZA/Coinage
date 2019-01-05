package transaction

import (
	"time"

	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
)

// Transaction represents a single transaction that occurred
type Transaction struct {
	// When is when the transaction occurred
	When time.Time

	// Arbitrary description of this transaction
	Description string

	// Labels is a user-specified key value set
	Labels map[string]string

	// Entries is the list of items in the transaction
	Entries []*entry.Entry
}
