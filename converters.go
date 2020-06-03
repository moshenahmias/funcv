package funcv

import (
	"fmt"
	"reflect"
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

// IsSupported returns true if the given value is a string
func (*StringConverter) IsSupported(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.String
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

// IsSupported returns true if the given value is an interger
// or can be converted to integer
func (*IntegerConverter) IsSupported(v interface{}) bool {
	return reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf(int64(0)))
}

// BooleanConverter is used to convert string represented integer
// arguments to integers, "true" is converted to true, "false"
// is converted to false, it uses sensitive compare as default
type BooleanConverter struct {
	Insensitive bool // how to compare the input
}

// Convert the given argument to boolean
func (c *BooleanConverter) Convert(arg string) (interface{}, error) {
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

// IsSupported returns true if the given value is a boolean
// or can be converted to boolean
func (*BooleanConverter) IsSupported(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Bool
}

// FloatConverter is used to convert string represented float
// arguments to float
type FloatConverter struct{}

// Convert the given argument to float
func (*FloatConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return nil, ErrInvalidValue
	}

	i, err := strconv.ParseFloat(arg, 64)

	if err != nil {
		return nil, fmt.Errorf("funcv: failed to parse float var %v (%w)", arg, err)
	}

	return i, nil
}

// IsSupported returns true if the given value is a float
// or can be converted to float
func (*FloatConverter) IsSupported(v interface{}) bool {
	return reflect.TypeOf(v).ConvertibleTo(reflect.TypeOf(float64(0)))
}
