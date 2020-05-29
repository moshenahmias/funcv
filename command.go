package funcv

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	isValidVarName = regexp.MustCompile(`^[0-9a-zA-Z\-_]+$`).MatchString
)

type command struct {
	args   []Arg
	err    error
	params []interface{}
	exe    string
	desc   string
}

func (c *command) AddArg(arg Arg) Builder {
	c.args = append(c.args, arg)
	return c
}

func (c *command) AddConst(text string, insensitive bool) Builder {
	if c.err != nil {
		return c
	}

	if text == "" {
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

func (c *command) AddDefStrVar(name, def, desc string) DefultedVariablesBuilder {
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

func (c *command) AddDefIntVar(name string, def, base int, desc string) DefultedVariablesBuilder {
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
		command:    c}

	return fb.AddBoolFlag(name, desc)
}

func (c *command) Create() Command {
	if c.err != nil {
		return c
	}

	if len(c.args) == 0 {
		c.err = ErrNoArguments
	}

	return c
}

func (c *command) Execute(args []string, fn interface{}) error {
	if c.err != nil {
		return c.err
	}

	c.params = nil

	var err error
	var params []interface{}

	for _, arg := range c.args {

		args, params, err = arg.Extract(args)

		if err != nil {
			return err
		}

		c.params = append(c.params, params...)
	}

	if len(args) > 0 {
		return ErrUnknownArgs
	}

	var in []reflect.Value

	vfn := reflect.ValueOf(fn)

	if vfn.Type().Kind() != reflect.Func {
		return fmt.Errorf("funcv: invalid function [%v]", vfn.Type().Kind())
	}

	if vfn.Type().NumIn() != len(c.params) {
		return fmt.Errorf("funcv: invalid function params count [%d/%d]", vfn.Type().NumIn(), len(c.params))
	}

	for i, param := range c.params {
		v := reflect.ValueOf(param)

		if !v.Type().ConvertibleTo(vfn.Type().In(i)) {
			return fmt.Errorf("funcv: can't convert param %v to %v", v.Type(), vfn.Type().In(i))
		}

		in = append(in, v.Convert(vfn.Type().In(i)))
	}

	vfn.Call(in)

	return nil
}

func (c *command) String() string {
	if len(c.args) == 0 {
		return c.exe
	}

	var sb strings.Builder

	if c.desc != "" {
		sb.WriteString(fmt.Sprintf("%s: ", c.desc))
	}

	if c.exe != "" {
		sb.WriteString(fmt.Sprintf("%s ", c.exe))
	}

	for i, arg := range c.args {
		sb.WriteString(arg.String())

		if i+1 < len(c.args) {
			sb.WriteString(" ")
		} else {
			sb.WriteString("\r\n\r\n")
		}
	}

	for i, arg := range c.args {

		desc := arg.Description()

		if desc != "" {
			sb.WriteString(fmt.Sprintf("\t%s", desc))

			if i+1 < len(c.args) {
				sb.WriteString("\r\n")
			}
		}

	}

	return sb.String()
}
