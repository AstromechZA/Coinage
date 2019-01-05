package yamlenc

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"

	"github.com/AstromechZA/coinage/pkg/core/transaction"
	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
	"github.com/AstromechZA/coinage/pkg/read"
	"github.com/AstromechZA/coinage/pkg/read/line"
)

type raw struct {
	When        string            `yaml:"when"`
	Description string            `yaml:"description"`
	Labels      map[string]string `yaml:"labels"`
	Entries     []string          `yaml:"entries"`
}

func decodeInner(decoder *yaml.Decoder) (*transaction.Transaction, error) {
	var tmp raw
	var err error

	if err := decoder.Decode(&tmp); err != nil {
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

func DecodeTransaction(reader io.Reader) (*transaction.Transaction, error) {
	return decodeInner(yaml.NewDecoder(reader))
}

func DecodeTransactions(reader io.Reader) ([]*transaction.Transaction, error) {
	decoder := yaml.NewDecoder(reader)
	out := make([]*transaction.Transaction, 0)

	for {
		t, err := decodeInner(decoder)
		if err != nil {
			if err == io.EOF {
				return out, nil
			}
			return nil, err
		}
		out = append(out, t)
	}
}
