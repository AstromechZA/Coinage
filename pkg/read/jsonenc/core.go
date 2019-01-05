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

func decodeRaw(r *raw) (*transaction.Transaction, error) {
	var err error

	out := &transaction.Transaction{
		Description: r.Description,
		Labels:      r.Labels,
		Entries:     make([]*entry.Entry, len(r.Entries)),
	}

	if out.When, err = read.ParseDateOrTime(r.When); err != nil {
		return nil, err
	}

	for i, e := range r.Entries {
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

func DecodeTransaction(reader io.Reader) (*transaction.Transaction, error) {
	tmp := new(raw)
	if err := json.NewDecoder(reader).Decode(tmp); err != nil {
		return nil, err
	}
	return decodeRaw(tmp)
}

func DecodeTransactions(reader io.Reader) ([]*transaction.Transaction, error) {
	var tmp []raw
	var err error

	if err := json.NewDecoder(reader).Decode(&tmp); err != nil {
		return nil, err
	}

	out := make([]*transaction.Transaction, len(tmp))
	for i, tt := range tmp {
		out[i], err = decodeRaw(&tt)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
