package accountname

import (
	"fmt"
	"regexp"
)

const BadCharacterSet = `[[:space:]:\p{C}\p{Z}]`
const MinAccountNamePartLen = 2
const MaxAccountNamePartLen = 100

func CheckPart(input string) error {
	if len(input) < MinAccountNamePartLen {
		return fmt.Errorf("it is shorter than %d", MinAccountNamePartLen)
	}
	if len(input) > MaxAccountNamePartLen {
		return fmt.Errorf("it is longer than %d", MaxAccountNamePartLen)
	}
	m := regexp.MustCompile(BadCharacterSet).FindStringIndex(input)
	if m != nil {
		return fmt.Errorf("it contains bad character `%s` at position %d", input[m[0]:m[1]], m[0])
	}
	return nil
}

func Check(input []string) error {
	if len(input) == 0 {
		return fmt.Errorf("it is empty")
	}
	for i, p := range input {
		if err := CheckPart(p); err != nil {
			return fmt.Errorf("part %d is invalid: %s", i+1, err)
		}
	}
	return nil
}
