package line

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AstromechZA/coinage/pkg/core/accountname"
	"github.com/AstromechZA/coinage/pkg/core/amount"
	"github.com/AstromechZA/coinage/pkg/core/commodity"
	"github.com/AstromechZA/coinage/pkg/core/transaction/entry"
	"github.com/AstromechZA/coinage/pkg/decext"
)

var linePattern = regexp.MustCompile(regexp.MustCompile(`\s`).ReplaceAllString(`
^
(?P<account>[^\s]+)
\s+
(?:
	(?P<value>[^\s]+)
	\s+
)?
(?P<valuecommodity>[^\s]+)
$
`, ""))

func StringToEntry(input string) (*entry.Entry, error) {
	m := linePattern.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("did not match the correct Line format")
	}
	accountName, value, valueCommodity := m[1], m[2], commodity.Commodity(m[3])

	output := &entry.Entry{
		Account: strings.Split(accountName, ":"),
	}

	if err := accountname.Check(output.Account); err != nil {
		return nil, fmt.Errorf("account name is invalid: %s", err)
	}

	if err := commodity.Check(valueCommodity); err != nil {
		return nil, fmt.Errorf("value commodity is invalid: %s", err)
	}
	output.Value = *amount.NewNil(valueCommodity)

	if len(value) != 0 {
		d, err := decext.NewFromString(value)
		if err != nil {
			return nil, err
		}
		output.Value.Value = d
	}

	return output, nil
}

func MustStringToEntry(input string) *entry.Entry {
	t, err := StringToEntry(input)
	if err != nil {
		panic(fmt.Errorf("could not create entry from '%s'", err))
	}
	return t
}
