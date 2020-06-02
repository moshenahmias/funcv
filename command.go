package funcv

import (
	"fmt"
	"io"
	"reflect"
	"regexp"
)

var (
	isValidVarName   = regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`).MatchString
	isValidConstName = regexp.MustCompile(`[^-\s]`).MatchString
)

type command struct {
	args   []Argument
	err    error
	params []interface{}
	desc   string
}

func (c *command) AddArgument(arg Argument) Builder {
	c.args = append(c.args, arg)
	return c
}

func (c *command) AddConstant(text string, insensitive bool) Builder {
	if c.err != nil {
		return c
	}

	if !isValidConstName(text) {
		c.err = fmt.Errorf("funcv: invalid constant [arg %d]", len(c.args))
		return c
	}

	return c.AddArgument(&constant{text: text, insensitive: insensitive})
}

func (c *command) AddVariable(name, desc string, conv Converter) Builder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	return c.AddArgument(&variable{name: name, desc: desc, conv: conv})
}

func (c *command) AddVariableWithDefault(name, desc string, conv Converter, def interface{}) ClosingBuilder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	if def == nil {
		c.err = fmt.Errorf("funcv: default for variable %s is nil [arg %d]", name, len(c.args))
		return c
	}

	if !conv.IsSupported(def) {
		c.err = fmt.Errorf("funcv: invalid default %v for var %s [arg %d]", def, name, len(c.args))
		return c
	}

	return c.AddArgument(&variable{name: name, desc: desc, conv: conv, def: def})
}

func (c *command) AddFlag(name, desc string, conv Converter, def interface{}) Builder {
	if c.err != nil {
		return c
	}

	fb := &flagsBuilder{
		converters: make(map[string]Converter),
		values:     make(map[string]interface{}),
		founddefs:  make(map[string]interface{}),
		defaults:   make(map[string]interface{}),
		command:    c}

	return fb.AddFlag(name, desc, conv, def)
}

func (c *command) AddParameterlessFlag(name, desc string, conv Converter, found, missing interface{}) Builder {
	if c.err != nil {
		return c
	}

	fb := &flagsBuilder{
		converters: make(map[string]Converter),
		values:     make(map[string]interface{}),
		founddefs:  make(map[string]interface{}),
		defaults:   make(map[string]interface{}),
		command:    c}

	return fb.AddParameterlessFlag(name, desc, conv, found, missing)
}

func (c *command) AddVariadic(name, desc string, conv Converter) Compiler {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	return c.AddArgument(&variadic{name: name, desc: desc, conv: conv})
}

func (c *command) Compile() (Command, error) {
	if c.err != nil {
		return nil, c.err
	}

	if len(c.args) == 0 {
		return nil, ErrNoArguments
	}

	return c, nil
}

func (c *command) MustCompile() Command {
	cmd, err := c.Compile()

	if err != nil {
		panic(err)
	}

	return cmd
}

func (c *command) ToGroup(grp *Group, fn interface{}) error {
	cmd, err := c.Compile()

	if err != nil {
		return err
	}

	grp.Add(cmd, fn)

	return nil
}

func (c *command) Execute(args []string, fn interface{}) (int, error) {
	n := 0

	if c.err != nil {
		return n, c.err
	}

	c.params = nil

	var err error
	var params []interface{}

	for _, arg := range c.args {
		l := len(args)
		args, params, err = arg.Extract(args)
		n += l - len(args)

		if err != nil {
			return n, err
		}

		c.params = append(c.params, params...)
	}

	if len(args) > 0 {
		return n, fmt.Errorf("funcv: %v (%w)", args, ErrUnknownArgs)
	}

	if fn == nil {
		return n, nil
	}

	vfn := reflect.ValueOf(fn)

	if vfn.Type().Kind() != reflect.Func {
		return n, fmt.Errorf("funcv: invalid function [%v]", vfn.Type().Kind())
	}

	if vfn.Type().NumIn() != len(c.params) {

		if !vfn.Type().IsVariadic() {
			return n, fmt.Errorf("funcv: invalid function params count [count: %d, input: %d]", vfn.Type().NumIn(), len(c.params))
		}

		if len(c.params) < vfn.Type().NumIn()-1 {
			return n, fmt.Errorf("funcv: invalid variadic function params count [count: %d..inf, input=%d]", vfn.Type().NumIn()-1, len(c.params))
		}
	}

	var in []reflect.Value

	i := 0

	for _, param := range c.params {
		v := reflect.ValueOf(param)
		t := vfn.Type().In(i)

		if i+1 == vfn.Type().NumIn() && vfn.Type().IsVariadic() {
			t = t.Elem()
		} else {
			i++
		}

		if !v.Type().ConvertibleTo(t) {
			return n, fmt.Errorf("funcv: can't convert param %v to %v", v.Type(), t)
		}

		in = append(in, v.Convert(t))
	}

	vfn.Call(in)

	return n, nil
}

func (c *command) WriteTo(w io.Writer) (int64, error) {
	var written int64

	if len(c.args) == 0 {
		return written, nil
	}

	if c.desc != "" {
		if n, err := fmt.Fprintf(w, "%s:\t", c.desc); err == nil {
			written += int64(n)
		} else {
			return written + int64(n), err
		}
	} else if n, err := fmt.Fprint(w, "\t"); err == nil {
		written += int64(n)
	} else {
		return written + int64(n), err
	}

	for i, arg := range c.args {

		if i == 0 {
			if n, err := fmt.Fprint(w, "> "); err == nil {
				written += int64(n)
			} else {
				return written + int64(n), err
			}
		}

		if n, err := fmt.Fprintf(w, "%s", arg.String()); err == nil {
			written += int64(n)
		} else {
			return written + int64(n), err
		}

		if i+1 < len(c.args) {
			if n, err := fmt.Fprint(w, " "); err == nil {
				written += int64(n)
			} else {
				return written + int64(n), err
			}
		} else if n, err := fmt.Fprint(w, "\n"); err == nil {
			written += int64(n)
		} else {
			return written + int64(n), err
		}
	}

	for _, arg := range c.args {
		if n, err := arg.WriteTo(w); err == nil {
			written += n
		} else {
			return written + n, err
		}
	}

	return written, nil
}
