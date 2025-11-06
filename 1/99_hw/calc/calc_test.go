package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestGetValueCalculatedByRPN(t *testing.T) {
	var cases = []struct {
		expected float64
		er       error
		input    string
	}{
		{expected: 5,
			er:    nil,
			input: "2 3 +\n",
		},

		{expected: 26,
			er:    nil,
			input: "2 3 * 4 5 * +\n",
		},

		{expected: (2. / (3 - (4 + (5 * 6)))),
			er:    nil,
			input: "2 3 4 5 6 * + - /\n",
		},

		{expected: 15,
			er:    nil,
			input: "1 2 3 4 + * + =\n",
		},

		{expected: 21,
			er:    nil,
			input: "1 2 + 3 4 + * =\n",
		},

		{expected: 0,
			er:    errors.New("1"),
			input: "1 2 3 * * * * = = = =\n",
		},

		{expected: 0,
			er:    errors.New("2"),
			input: "1 2 3 * gq34g 3 5 -\n",
		},

		{expected: 0,
			er:    errors.New("2"),
			input: "1 2 3 * A 3 5 -\n",
		},

		{expected: 0,
			er:    errors.New("3"),
			input: "1 2 + 0 /\n",
		},

		{expected: 1,
			er:    nil,
			input: "1\n",
		},
	}
	for _, item := range cases {
		result, err := GetValueCalculatedByRPN(item.input)
		if !reflect.DeepEqual(err, item.er) {
			t.Error("expected error", item.er, "have", err)
		}
		if result != item.expected {
			t.Error("expected", item.expected, "have", result)
		}
	}
}
