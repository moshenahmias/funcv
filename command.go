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
	args   []Arg
	err    error
	params []interface{}
	desc   string
}

func (c *command) AddArg(arg Arg) Builder {
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

	return c.AddArg(&constant{text: text, insensitive: insensitive})
}

func (c *command) AddStrVar(name, desc string) Builder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	return c.AddArg(&strVar{name: name, desc: desc})
}

func (c *command) AddIntVar(name string, base int, desc string) Builder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	if base < 0 {
		c.err = fmt.Errorf("funcv: invalid base %d [arg %d]", base, len(c.args))
		return c
	}

	return c.AddArg(&intVar{name: name, desc: desc, base: base})
}

func (c *command) AddStrVarWithDefault(name, def, desc string) DefultedVariablesBuilder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	c.args = append(c.args, &defStrVar{name: name, desc: desc, def: def})
	return c
}

func (c *command) AddIntVarWithDefault(name string, def, base int, desc string) DefultedVariablesBuilder {
	if c.err != nil {
		return c
	}

	if !isValidVarName(name) {
		c.err = fmt.Errorf("funcv: invalid var name %s [arg %d]", name, len(c.args))
		return c
	}

	if base < 0 {
		c.err = fmt.Errorf("funcv: invalid base %d [arg %d]", base, len(c.args))
		return c
	}

	c.args = append(c.args, &defIntVar{name: name, desc: desc, def: def, base: base})
	return c
}

func (c *command) AddStrFlag(name, def, desc string) FlagsBuilder {
	if c.err != nil {
		return c
	}

	fb := &flagsBuilder{
		converters: make(map[string]Converter),
		values:     make(map[string]interface{}),
		founddefs:  make(map[string]interface{}),
		defaults:   make(map[string]interface{}),
		command:    c}

	return fb.AddStrFlag(name, def, desc)
}

func (c *command) AddIntFlag(name string, def, base int, desc string) FlagsBuilder {
	if c.err != nil {
		return c
	}

	fb := &flagsBuilder{
		converters: make(map[string]Converter),
		values:     make(map[string]interface{}),
		founddefs:  make(map[string]interface{}),
		defaults:   make(map[string]interface{}),
		command:    c}

	return fb.AddIntFlag(name, def, base, desc)
}

func (c *command) AddBoolFlag(name, desc string) FlagsBuilder {
	if c.err != nil {
		return c
	}

	fb := &flagsBuilder{
		converters: make(map[string]Converter),
		values:     make(map[string]interface{}),
		founddefs:  make(map[string]interface{}),
		defaults:   make(map[string]interface{}),
		command:    c}

	return fb.AddBoolFlag(name, desc)
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
		return n, ErrUnknownArgs
	}

	if fn == nil {
		return n, nil
	}

	var in []reflect.Value

	vfn := reflect.ValueOf(fn)

	if vfn.Type().Kind() != reflect.Func {
		return n, fmt.Errorf("funcv: invalid function [%v]", vfn.Type().Kind())
	}

	if vfn.Type().NumIn() != len(c.params) {
		return n, fmt.Errorf("funcv: invalid function params count [%d/%d]", vfn.Type().NumIn(), len(c.params))
	}

	for i, param := range c.params {
		v := reflect.ValueOf(param)

		if !v.Type().ConvertibleTo(vfn.Type().In(i)) {
			return n, fmt.Errorf("funcv: can't convert param %v to %v", v.Type(), vfn.Type().In(i))
		}

		in = append(in, v.Convert(vfn.Type().In(i)))
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
