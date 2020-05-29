package funcv

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	flagRegex       = regexp.MustCompile("^-([a-zA-Z])$|^--([a-zA-Z][a-zA-Z]+)$")
	isValidFlagName = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
)

func toFlag(name string) string {
	if !isValidFlagName(name) {
		return ""
	}

	if len(name) == 1 {
		return "-" + name
	}

	return "--" + name
}

func extractFlagName(arg string) string {
	m := flagRegex.FindStringSubmatch(arg)

	if len(m) != 3 {
		return ""
	}

	if m[2] == "" {
		return m[1]
	}

	return m[2]
}

type flagsBuilder struct {
	converters map[string]Converter
	values     map[string]interface{}
	flags      []string
	desc       []string
	command    *command
}

func (b *flagsBuilder) toParams() ([]interface{}, error) {

	var params []interface{}

	for _, name := range b.flags {

		v, found := b.values[name]

		if !found {
			return nil, fmt.Errorf("funcv: flag %s not found", name)
		}

		params = append(params, v)
	}

	return params, nil
}

func (b *flagsBuilder) Extract(args []string) ([]string, []interface{}, error) {

	for len(args) > 0 {

		name := extractFlagName(args[0])

		if name == "" {
			params, err := b.toParams()
			return args, params, err
		}

		if _, found := b.values[name]; !found {
			params, err := b.toParams()
			return args, params, err
		}

		args = args[1:]

		var v string

		conv, found := b.converters[name]

		if !found {
			return nil, nil, fmt.Errorf("funcv: missing converter for %s", name)
		}

		var i int

		if len(args) > 0 {
			next := extractFlagName(args[0])

			if next == "" {
				v = args[0]
				i = 1
			}
		}

		if conval, err := conv.Convert(v); err == nil {
			b.values[name] = conval
		} else {
			i = 0
		}

		args = args[i:]
	}

	params, err := b.toParams()
	return args, params, err
}

func (b *flagsBuilder) AddStrFlag(name, def, desc string) FlagsBuilder {
	if b.command.err != nil {
		return b
	}

	if !isValidFlagName(name) {
		b.command.err = fmt.Errorf("funcv: invalid flag name %s [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	b.desc = append(b.desc, desc)
	b.flags = append(b.flags, name)
	b.values[name] = def
	b.converters[name] = new(StringConverter)
	return b
}

func (b *flagsBuilder) AddIntFlag(name string, def, base int, desc string) FlagsBuilder {
	if b.command.err != nil {
		return b
	}

	if !isValidFlagName(name) {
		b.command.err = fmt.Errorf("funcv: invalid flag name %s [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	if base < 0 {
		b.command.err = fmt.Errorf("funcv: invalid base %d [arg %d]", base, len(b.command.args)+len(b.flags))
		return b
	}

	b.desc = append(b.desc, desc)
	b.flags = append(b.flags, name)
	b.values[name] = def
	b.converters[name] = &IntegerConverter{base}
	return b
}

func (b *flagsBuilder) AddBoolFlag(name, desc string) FlagsBuilder {
	if b.command.err != nil {
		return b
	}

	if !isValidFlagName(name) {
		b.command.err = fmt.Errorf("funcv: invalid flag name %s [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	b.desc = append(b.desc, desc)
	b.flags = append(b.flags, name)
	b.values[name] = false
	b.converters[name] = &BoolConverter{true}
	return b
}

func (b *flagsBuilder) AddConst(text string, insensitive bool) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddConst(text, insensitive)
}

func (b *flagsBuilder) AddStrVar(name, desc string) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddStrVar(name, desc)
}

func (b *flagsBuilder) AddIntVar(name string, base int, desc string) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddIntVar(name, base, desc)
}

func (b *flagsBuilder) AddDefStrVar(name, def, desc string) DefultedVariablesBuilder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddDefStrVar(name, def, desc)
}

func (b *flagsBuilder) AddDefIntVar(name string, def, base int, desc string) DefultedVariablesBuilder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddDefIntVar(name, def, base, desc)
}

func (b *flagsBuilder) AddArg(arg Arg) Builder {
	if b.command.err != nil {
		return b
	}

	return b.command.AddArg(arg)
}

func (b *flagsBuilder) Create() Command {
	b.command.args = append(b.command.args, b)
	return b.command.Create()
}

func (b *flagsBuilder) Description() string {
	var sb strings.Builder

	for i, name := range b.flags {
		def, _ := b.values[name]
		sb.WriteString(fmt.Sprintf("%s\t%s (default: %v)", toFlag(name), b.desc[i], def))

		if i+1 < len(b.flags) {
			sb.WriteString("\r\n\t")
		}
	}

	return sb.String()
}

func (b *flagsBuilder) String() string {
	var sb strings.Builder

	for i, name := range b.flags {
		sb.WriteString(fmt.Sprintf("[%s]", toFlag(name)))

		if i+1 < len(b.flags) {
			sb.WriteString(" ")
		}
	}

	return sb.String()
}
