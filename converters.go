package funcv

import (
	"fmt"
	"strconv"
	"strings"
)

// StringConverter is used to convert string arguments to strings
type StringConverter struct{}

// Convert returns the given arg as-is
func (*StringConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return nil, ErrInvalidValue
	}

	return arg, nil
}

// IntegerConverter is used to convert string represented integer
// arguments to integers, the default base is decimal
type IntegerConverter struct {
	Base int // input base (0 or less defaults to decimal)
}

// Convert the given argument to integer
func (c *IntegerConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return nil, ErrInvalidValue
	}

	base := 10

	if c != nil && c.Base > 0 {
		base = c.Base
	}

	i, err := strconv.ParseInt(arg, base, 64)

	if err != nil {
		return nil, fmt.Errorf("funcv: failed to parse int var %v (%w)", arg, err)
	}

	return i, nil
}

// BoolConverter is used to convert string represented integer
// arguments to integers, "true" is converted to true, "false"
// is converted to false, it uses sensitive compare as default
type BoolConverter struct {
	Insensitive bool // how to compare the input
}

// Convert the given argument to boolean
func (c *BoolConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return true, nil
	}

	insensitive := false

	if c != nil {
		insensitive = c.Insensitive
	}

	if insensitive {
		if strings.EqualFold(arg, "true") {
			return true, nil
		}

		if strings.EqualFold(arg, "false") {
			return false, nil
		}
	} else {
		if arg == "true" {
			return true, nil
		}

		if arg == "false" {
			return false, nil
		}
	}

	return nil, fmt.Errorf("funcv: cannot convert %s to boolean", arg)
}
