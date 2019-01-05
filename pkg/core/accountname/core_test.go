package accountname

import (
	"fmt"
	"strings"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestCheckAccountNamePart(t *testing.T) {
	for _, c := range []struct {
		input    string
		expected error
	}{
		{"Somepart", nil},
		{"", fmt.Errorf("it is shorter than 2")},
		{strings.Repeat("a", 101), fmt.Errorf("it is longer than 100")},
		{"has a space", fmt.Errorf("it contains bad character ` ` at position 3")},
		{"unicodeĤĤĤtest", nil},
		{"😀", nil},
	} {
		t.Run(c.input, func(t *testing.T) {
			assert.ShouldEqual(t, CheckPart(c.input), c.expected)
		})
	}
}

func TestCheckAccountName(t *testing.T) {
	for _, c := range []struct {
		input    []string
		expected error
	}{
		{[]string{"Somepart"}, nil},
		{[]string{}, fmt.Errorf("it is empty")},
		{[]string{"something", "with", "spaces"}, nil},
		{[]string{"bad", "thing", "", ""}, fmt.Errorf("part 3 is invalid: it is shorter than 2")},
	} {
		t.Run(strings.Join(c.input, ":"), func(t *testing.T) {
			assert.ShouldEqual(t, Check(c.input), c.expected)
		})
	}
}
