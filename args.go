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

type variable struct {
	name string
	desc string
	conv Converter
	def  interface{}
}

func (v *variable) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		if v.def != nil {
			return args, []interface{}{v.def}, nil
		}

		return args, nil, ErrArgNotFound
	}

	p, err := v.conv.Convert(args[0])

	if err != nil {
		return args, nil, err
	}

	return args[1:], []interface{}{p}, nil
}

func (v *variable) WriteTo(w io.Writer) (int64, error) {
	if v.def != nil {
		n, err := fmt.Fprintf(w, "\n\t%s\t%s (default: %v)", v.name, v.desc, v.def)
		return int64(n), err
	}

	n, err := fmt.Fprintf(w, "\n\t%s\t%s", v.name, v.desc)

	return int64(n), err
}

func (v *variable) String() string {
	if v.def != nil {
		return fmt.Sprintf("[%s]", v.name)
	}

	return fmt.Sprintf("<%s>", v.name)
}

type variadic struct {
	name string
	desc string
	conv Converter
}

func (v *variadic) Extract(args []string) ([]string, []interface{}, error) {
	if len(args) == 0 {
		return args, nil, nil
	}

	var params []interface{}

	for i, arg := range args {
		p, err := v.conv.Convert(arg)

		if err != nil {
			return args[i:], nil, err
		}

		params = append(params, p)
	}

	return nil, params, nil
}

func (v *variadic) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "\n\t%s...\t%s", v.name, v.desc)
	return int64(n), err
}

func (v *variadic) String() string {
	return fmt.Sprintf("[%s...]", v.name)
}
