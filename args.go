package funcv

import (
	"fmt"
	"io"
	"strings"
)

type constant struct {
	text        string
	insensitive bool
}

func (c *constant) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, nil, ErrArgNotFound
	}

	if c.insensitive {
		if !strings.EqualFold(args[0], c.text) {
			return args, nil, ErrArgNotFound
		}
	} else {
		if args[0] != c.text {
			return args, nil, ErrArgNotFound
		}
	}

	return args[1:], nil, nil
}

func (*constant) WriteTo(w io.Writer) (int64, error) {
	return 0, nil
}

func (c *constant) String() string {
	return c.text
}

type strVar struct {
	name string
	desc string
}

func (*strVar) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, nil, ErrArgNotFound
	}

	s, err := new(StringConverter).Convert(args[0])

	if err != nil {
		return args, nil, err
	}

	return args[1:], []interface{}{s}, nil
}

func (v *strVar) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "\n\t%s\t%s", v.name, v.desc)
	return int64(n), err
}

func (v *strVar) String() string {
	return fmt.Sprintf("<%s>", v.name)
}

type intVar struct {
	name string
	desc string
	base int
}

func (v *intVar) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, nil, ErrArgNotFound
	}

	i, err := (&IntegerConverter{v.base}).Convert(args[0])

	if err != nil {
		return args, nil, err
	}

	return args[1:], []interface{}{i}, nil
}

func (v *intVar) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "\n\t%s\t%s (base: %d)", v.name, v.desc, v.base)
	return int64(n), err
}

func (v *intVar) String() string {
	return fmt.Sprintf("<%s>", v.name)
}

type defStrVar struct {
	name string
	desc string
	def  string
}

func (v *defStrVar) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, []interface{}{v.def}, nil
	}

	s, err := new(StringConverter).Convert(args[0])

	if err != nil {
		return args, nil, err
	}

	return args[1:], []interface{}{s}, nil
}

func (v *defStrVar) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "\n\t%s\t%s (default: %s)", v.name, v.desc, v.def)
	return int64(n), err
}

func (v *defStrVar) String() string {
	return fmt.Sprintf("[%s]", v.name)
}

type defIntVar struct {
	name string
	desc string
	def  int
	base int
}

func (v *defIntVar) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, []interface{}{v.def}, nil
	}

	i, err := (&IntegerConverter{v.base}).Convert(args[0])

	if err != nil {
		return args, nil, err
	}

	return args[1:], []interface{}{i}, nil
}

func (v *defIntVar) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "\n\t%s\t%s (base: %d, default: %d)", v.name, v.desc, v.base, v.def)
	return int64(n), err
}

func (v *defIntVar) String() string {
	return fmt.Sprintf("[%s]", v.name)
}
