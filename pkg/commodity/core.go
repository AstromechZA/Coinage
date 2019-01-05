package commodity

import (
	"fmt"
	"regexp"
)

type Commodity string

const charSet = `\p{Sc}\p{L}`
const BadCharacterSet = `[^` + charSet + `]`
const MinCommodityLen = 1
const MaxCommodityLen = 8

func Check(input Commodity) error {
	if len(input) < MinCommodityLen {
		return fmt.Errorf("it is shorter than %d", MinCommodityLen)
	}
	if len(input) > MaxCommodityLen {
		return fmt.Errorf("it is longer than %d", MaxCommodityLen)
	}
	m := regexp.MustCompile(BadCharacterSet).FindStringIndex(string(input))
	if m != nil {
		return fmt.Errorf("it contains bad character `%s` at position %d", input[m[0]:m[1]], m[0])
	}
	return nil
}

func IsValid(input Commodity) (bool, error) {
	if err := Check(input); err != nil {
		return false, err
	}
	return true, nil
}
