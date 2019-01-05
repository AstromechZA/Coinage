package line

import (
	"fmt"
	"regexp"

	"github.com/AstromechZA/coinage/pkg/transaction"

	"github.com/AstromechZA/coinage/pkg/accountname"

	"github.com/AstromechZA/coinage/pkg/amount"
	"github.com/AstromechZA/coinage/pkg/decext"

	"github.com/AstromechZA/coinage/pkg/commodity"
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
(?:
	\s+
	for
	\s+
	(?P<price>[^\s]+)
	\s+
	(?P<pricecommodity>[^\s]+)
	(?P<each>
		\s+
		each
	)?
)?
$
`, ""))

func StringToEntry(input string) (*transaction.Entry, error) {
	m := linePattern.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("did not match the correct Line format")
	}
	accountName, value, valueCommodity, price, priceCommodity, each := m[1], m[2], commodity.Commodity(m[3]), m[4], commodity.Commodity(m[5]), m[6]

	output := &transaction.Entry{
		Account: accountName,
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

	if len(price) != 0 || len(priceCommodity) > 0 {
		d, err := decext.NewFromString(price)
		if err != nil {
			return nil, err
		}
		output.Price.Value = d.Abs(d)

		if err := commodity.Check(priceCommodity); err != nil {
			return nil, fmt.Errorf("price commodity is invalid: %s", err)
		}
		output.Price.Commodity = priceCommodity
	}

	if len(each) != 0 {
		output.Price.Value.Mul(output.Price.Value, output.Value.Value)
	}

	return output, nil
}
