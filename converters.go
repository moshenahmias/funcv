package funcv

import (
	"fmt"
	"strconv"
	"strings"
)

type StringConverter struct{}

func (*StringConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return nil, ErrInvalidValue
	}

	return arg, nil
}

type IntegerConverter struct {
	base int
}

func (c *IntegerConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return nil, ErrInvalidValue
	}

	base := 10

	if c != nil && c.base > 0 {
		base = c.base
	}

	i, err := strconv.ParseInt(arg, base, 64)

	if err != nil {
		return nil, fmt.Errorf("funcv: failed to parse int var %v (%w)", arg, err)
	}

	return i, nil
}

type BoolConverter struct {
	insensitive bool
}

func (c *BoolConverter) Convert(arg string) (interface{}, error) {
	if arg == "" {
		return true, nil
	}

	insensitive := false

	if c != nil {
		insensitive = c.insensitive
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
