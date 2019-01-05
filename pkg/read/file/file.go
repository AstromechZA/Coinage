package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AstromechZA/coinage/pkg/core/transaction"
	"github.com/AstromechZA/coinage/pkg/read/jsonenc"
	"github.com/AstromechZA/coinage/pkg/read/yamlenc"
)

func TransactionsFromFile(p string, m func(t *transaction.Transaction) error) error {
	s, err := os.Stat(p)
	if err != nil {
		return fmt.Errorf("cannot read file: %s", err)
	}
	if s.IsDir() {
		return fmt.Errorf("path %s is a directory", p)
	}

	x := filepath.Ext(p)

	f, err := os.Open(p)
	if err != nil {
		return fmt.Errorf("cannot read file: %s", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	var transactions []*transaction.Transaction
	switch strings.ToLower(x) {
	case ".yaml":
		transactions, err = yamlenc.DecodeTransactions(f)
	case ".json":
		fallthrough
	case "":
		transactions, err = jsonenc.DecodeTransactions(f)
	default:
		return fmt.Errorf("cannot read transactions from a %s file", x)
	}
	if err != nil {
		return fmt.Errorf("failed to read transactions: %s", err)
	}

	for i, tt := range transactions {
		if err = m(tt); err != nil {
			return fmt.Errorf("failed to process transaction %d: %s", i+1, err)
		}
	}

	return nil
}

func TransactionsFromDir(p string, m func(t *transaction.Transaction) error) error {
	return filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		x := filepath.Ext(p)
		switch strings.ToLower(x) {
		case ".yaml":
		case ".json":
		case "":
		default:
			return nil
		}
		return TransactionsFromFile(path, m)
	})
}
