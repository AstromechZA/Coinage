package commodity

import (
	"fmt"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestCheck(t *testing.T) {
	for _, c := range []struct {
		input    Commodity
		expected error
	}{
		{"ZAR", nil},
		{"AAPL", nil},
		{"aapl", nil},
		{"£", nil},
		{"$", nil},
		{"", fmt.Errorf("it is shorter than 1")},
		{"muchtoolong", fmt.Errorf("it is longer than 8")},
	} {
		t.Run(string(c.input), func(t *testing.T) {
			assert.ShouldEqual(t, Check(c.input), c.expected)
		})
	}
}
