package jsonenc

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/AstromechZA/coinage/pkg/core/transaction"
	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
	"github.com/AstromechZA/coinage/pkg/read"
	"github.com/AstromechZA/coinage/pkg/read/line"
)

type raw struct {
	When        string            `json:"when"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	Entries     []string          `json:"entries"`
}

func DecodeTransaction(reader io.Reader) (*transaction.Transaction, error) {
	var tmp raw
	var err error

	if err := json.NewDecoder(reader).Decode(&tmp); err != nil {
		return nil, err
	}

	out := &transaction.Transaction{
		Description: tmp.Description,
		Labels:      tmp.Labels,
		Entries:     make([]*entry.Entry, len(tmp.Entries)),
	}

	if out.When, err = read.ParseDateOrTime(tmp.When); err != nil {
		return nil, err
	}

	for i, e := range tmp.Entries {
		out.Entries[i], err = line.StringToEntry(e)
		if err != nil {
			return nil, err
		}
	}

	if _, err = out.Balance(); err != nil {
		return nil, fmt.Errorf("transaction does not balance: %s", err)
	}

	return out, nil
}

func DecodeTransactions(reader io.Reader) ([]*transaction.Transaction, error) {
	var tmp []raw
	var err error

	if err := json.NewDecoder(reader).Decode(&tmp); err != nil {
		return nil, err
	}

	out := make([]*transaction.Transaction, len(tmp))
	for i, tt := range tmp {
		out[i] = &transaction.Transaction{
			Description: tt.Description,
			Labels:      tt.Labels,
			Entries:     make([]*entry.Entry, len(tt.Entries)),
		}

		if out[i].When, err = read.ParseDateOrTime(tt.When); err != nil {
			return nil, err
		}

		for i, e := range tt.Entries {
			out[i].Entries[i], err = line.StringToEntry(e)
			if err != nil {
				return nil, err
			}
		}

		if _, err = out[i].Balance(); err != nil {
			return nil, fmt.Errorf("transaction %d does not balance: %s", i, err)
		}
	}
	return out, nil
}
