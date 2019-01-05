package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/AstromechZA/coinage/pkg/core/commodity"

	"github.com/AstromechZA/coinage/pkg/core/transaction"
	"github.com/AstromechZA/coinage/pkg/read/file"
	"github.com/AstromechZA/coinage/pkg/tree"
)

func mainInner() error {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	if fs.NArg() != 1 {
		return fmt.Errorf("expected a single filepath as a source of transactions")
	}

	transactions := make([]*transaction.Transaction, 0)
	accounts := tree.New()

	if err := file.TransactionsFromDir(fs.Arg(0), func(t *transaction.Transaction) error {
		transactions = append(transactions, t)
		if e := accounts.Insert(t); e != nil {
			return e
		}
		return nil
	}); err != nil {
		return err
	}

	accounts.DepthFirstPreOrder([]string{}, func(account []string, node *tree.AccountTree) bool {
		if len(account) > 0 {
			keys := make([]commodity.Commodity, 0, len(node.TreeTotals))
			for c := range node.TreeTotals {
				keys = append(keys, c)
			}
			for _, c := range keys {
				fmt.Printf("%20s %s%s\n", node.TreeTotals[c].String()+" "+string(c), strings.Repeat("  ", len(account)), account[len(account)-1])
			}
		}
		return false
	})
	fmt.Println(strings.Repeat("-", 30))
	keys := make([]string, 0, len(accounts.TreeTotals))
	for c := range accounts.TreeTotals {
		keys = append(keys, string(c))
	}
	sort.Strings(keys)
	for _, c := range keys {
		fmt.Printf("%20s\n", accounts.TreeTotals[commodity.Commodity(c)].String()+" "+c)
	}

	return nil
}

func main() {
	if err := mainInner(); err != nil {
		if _, err = fmt.Fprintf(os.Stderr, "Error: %s\n", err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
