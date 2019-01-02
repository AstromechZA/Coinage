package transaction

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
)

const minCommodityLen = 1
const maxCommodityLen = 8
const minAccountNamePartLen = 2
const maxAccountNamePartLen = 100

func CheckAccountNamePart(input string) error {
	if len(input) < minAccountNamePartLen {
		return fmt.Errorf("it is shorter than %d", minAccountNamePartLen)
	}
	if len(input) > maxAccountNamePartLen {
		return fmt.Errorf("it is longer than %d", maxAccountNamePartLen)
	}
	m := regexp.MustCompile(`[[:space:]:\p{C}\p{Z}]`).FindStringIndex(input)
	if m != nil {
		return fmt.Errorf("it contains bad character `%s` at position %d", input[m[0]:m[1]], m[0])
	}
	return nil
}

func CheckAccountName(input string) error {
	if len(input) == 0 {
		return fmt.Errorf("it is empty")
	}
	parts := strings.Split(input, ":")
	for i, p := range parts {
		if err := CheckAccountNamePart(p); err != nil {
			return fmt.Errorf("part %d is invalid: %s", i+1, err)
		}
	}
	return nil
}

func CheckCommodity(input string) error {
	if len(input) < minCommodityLen {
		return fmt.Errorf("it is shorter than %d", minCommodityLen)
	}
	if len(input) > maxCommodityLen {
		return fmt.Errorf("it is longer than %d", maxCommodityLen)
	}
	m := regexp.MustCompile(`[^\p{Sc}\p{L}]`).FindStringIndex(input)
	if m != nil {
		return fmt.Errorf("it contains bad character `%s` at position %d", input[m[0]:m[1]], m[0])
	}
	return nil
}

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

func StringToLine(input string) (*TransactionLine, error) {
	m := linePattern.FindStringSubmatch(input)
	if m == nil {
		return nil, fmt.Errorf("did not match the correct line format")
	}
	accountName, value, valueCommodity, price, priceCommodity, each := m[1], m[2], m[3], m[4], m[5], m[6]

	output := &TransactionLine{
		Account: accountName,
	}

	if err := CheckAccountName(output.Account); err != nil {
		return nil, fmt.Errorf("account name is invalid: %s", err)
	}

	output.Value = new(Amount)

	if err := CheckCommodity(valueCommodity); err != nil {
		return nil, fmt.Errorf("value commodity is invalid: %s", err)
	}
	output.Value.Commodity = valueCommodity

	if len(value) != 0 {
		d, err := decimal.NewFromString(value)
		if err != nil {
			return nil, fmt.Errorf("value `%s` is invalid: %s", value, err)
		}
		output.Value.Value = &d
	}

	if len(price) != 0 || len(priceCommodity) > 0 {
		output.Price = new(Amount)

		d, err := decimal.NewFromString(price)
		if err != nil {
			return nil, fmt.Errorf("price `%s` is invalid: %s", price, err)
		}
		output.Price.Value = &d

		if err := CheckCommodity(priceCommodity); err != nil {
			return nil, fmt.Errorf("price commodity is invalid: %s", err)
		}
		output.Price.Commodity = priceCommodity
	}

	if len(each) != 0 {
		vv := output.Price.Value.Mul(*output.Value.Value)
		output.Price.Value = &vv
	}

	return output, nil
}
