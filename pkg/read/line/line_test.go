package line

import (
	"fmt"
	"testing"

	"github.com/AstromechZA/coinage/pkg/assert"
)

func TestStringToLine_with_value(t *testing.T) {
	l, err := StringToEntry("Assets:Cash    100 USD")
	if assert.ShouldEqual(t, err, nil) {
		assert.Equal(t, l.Account, []string{"Assets", "Cash"})
		assert.Equal(t, l.Value.Value.String(), "100")
		assert.Equal(t, string(l.Value.Commodity), "USD")
	}
}

func TestStringToLine_with_negative_value(t *testing.T) {
	l, err := StringToEntry("Assets:Cash    -3.333 £")
	if assert.ShouldEqual(t, err, nil) {
		assert.Equal(t, l.Account, []string{"Assets", "Cash"})
		assert.Equal(t, l.Value.Value.String(), "-3.333")
		assert.Equal(t, string(l.Value.Commodity), "£")
	}
}

func TestStringToLine_with_price(t *testing.T) {
	l, err := StringToEntry("Assets:Cash    3.333 £  for 10 GBP")
	if assert.ShouldEqual(t, err, nil) {
		assert.Equal(t, l.Account, []string{"Assets", "Cash"})
		assert.Equal(t, l.Value.Value.String(), "3.333")
		assert.Equal(t, string(l.Value.Commodity), "£")
		assert.Equal(t, l.Price.Value.String(), "10")
		assert.Equal(t, string(l.Price.Commodity), "GBP")
	}
}

func TestStringToLine_with_each_price(t *testing.T) {
	l, err := StringToEntry("Assets:Cash    3.333 £  for 10 GBP each")
	if assert.ShouldEqual(t, err, nil) {
		assert.Equal(t, l.Account, []string{"Assets", "Cash"})
		assert.Equal(t, l.Value.Value.String(), "3.333")
		assert.Equal(t, string(l.Value.Commodity), "£")
		assert.Equal(t, l.Price.Value.String(), "33.330")
		assert.Equal(t, string(l.Price.Commodity), "GBP")
	}
}

func TestStringToLine_errors(t *testing.T) {
	for _, c := range []struct {
		input    string
		expected error
	}{
		{"", fmt.Errorf("did not match the correct Line format")},
		{"invalid::", fmt.Errorf("did not match the correct Line format")},
		{"invalid:: 100 USD", fmt.Errorf("account name is invalid: part 2 is invalid: it is shorter than 2")},
		{"Assets:Cash 100", fmt.Errorf("value commodity is invalid: it contains bad character `1` at position 0")},
		{"Assets:Cash 100 __", fmt.Errorf("value commodity is invalid: it contains bad character `_` at position 0")},
		{"Assets:Cash flerp USD", fmt.Errorf("`flerp` is an invalid decimal")},
		{"Assets:Cash 100 USD for", fmt.Errorf("did not match the correct Line format")},
		{"Assets:Cash 100 USD for 1", fmt.Errorf("did not match the correct Line format")},
		{"Assets:Cash 100 USD for 1 __", fmt.Errorf("price commodity is invalid: it contains bad character `_` at position 0")},
		{"Assets:Cash 100 USD for 1 GBP hi", fmt.Errorf("did not match the correct Line format")},
		{"Assets:Cash 100 USD for flerp GBP", fmt.Errorf("`flerp` is an invalid decimal")},
	} {
		t.Run(c.input, func(t *testing.T) {
			_, err := StringToEntry(c.input)
			assert.ShouldEqual(t, err, c.expected)
		})
	}
}
