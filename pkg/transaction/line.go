package transaction

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
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

func StringToLine(input string) (*TransactionLine, error) {
	output := new(TransactionLine)

	parts := regexp.MustCompile(`[\t ]+`).Split(input, -1)
	if len(parts) == 0 {
		return nil, fmt.Errorf("expected at least an account name but got none")
	}

	output.Account, parts = parts[0], parts[1:]
	if err := CheckAccountName(output.Account); err != nil {
		return nil, fmt.Errorf("account name is invalid: %s", err)
	}

	// the line can have no value or price
	if len(parts) == 0 {
		return output, nil
	}

	// parse the value as the next component
	value, _, err := big.ParseFloat(parts[0], 10, 53, big.AwayFromZero)
	if err != nil {
		return nil, fmt.Errorf("value `%s` is invalid: %s", parts[0], err)
	}
	output.Value = &Amount{
		Value: value,
	}
	parts = parts[1:]

	// must have a commodity if it has a value
	if len(parts) == 0 {
		return nil, fmt.Errorf("value is missing a commodity symbol")
	}
	if err := CheckCommodity(parts[0]); err != nil {
		return nil, fmt.Errorf("value commodity is invalid: %s", err)
	}
	output.Value.Commodity = parts[0]
	parts = parts[1:]

	// the line can have price
	if len(parts) == 0 {
		return output, nil
	}

	if parts[0] != "for" {
		return nil, fmt.Errorf("expected next part to be `for` but got `%s`", parts[0])
	}
	parts = parts[1:]
	if len(parts) == 0 {
		return nil, fmt.Errorf("expected price value and commodity after `for`")
	}

	value, _, err = big.ParseFloat(parts[0], 10, 53, big.AwayFromZero)
	if err != nil {
		return nil, fmt.Errorf("price `%s` is invalid: %s", parts[0], err)
	}
	if value.Sign() == -1 {
		return nil, fmt.Errorf("price cannot be negative")
	}
	output.Price = &Amount{
		Value: value,
	}
	parts = parts[1:]

	// must have a commodity if it has a value
	if len(parts) == 0 {
		return nil, fmt.Errorf("price is missing a commodity symbol")
	}
	if err := CheckCommodity(parts[0]); err != nil {
		return nil, fmt.Errorf("price commodity is invalid: %s", err)
	}
	output.Price.Commodity = parts[0]
	parts = parts[1:]

	// the line can skip the each
	if len(parts) == 0 {
		return output, nil
	}

	if parts[0] == "each" {
		output.Price.Value.Mul(output.Price.Value, output.Value.Value)
		parts = parts[1:]
	}

	if len(parts) > 0 {
		return nil, fmt.Errorf("has unparsed content: `%s`", parts)
	}
	return output, nil
}
