package funcv

import (
	"fmt"
	"io"
	"regexp"
	"strings"
)

var (
	flagRegex       = regexp.MustCompile(`^-([a-zA-Z])$|^--([a-zA-Z][a-zA-Z]+)$`)
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
	founddefs  map[string]interface{}
	defaults   map[string]interface{}
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

		if def, found := b.founddefs[name]; found {
			b.values[name] = def
		} else {
			delete(b.values, name)
		}

		args = args[1:]

		var v string

		conv, found := b.converters[name]

		if !found {
			return args, nil, fmt.Errorf("funcv: missing converter for %s", name)
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

func (b *flagsBuilder) AddFlag(name, desc string, conv Converter, def interface{}) Builder {
	if b.command.err != nil {
		return b
	}

	if !isValidFlagName(name) {
		b.command.err = fmt.Errorf("funcv: invalid flag name %s [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	if def == nil {
		b.command.err = fmt.Errorf("funcv: default for flag %s is nil [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	if !conv.IsSupported(def) {
		b.command.err = fmt.Errorf("funcv: invalid default %v for flag %s [arg %d]", def, name, len(b.command.args)+len(b.flags))
		return b
	}

	b.desc = append(b.desc, desc)
	b.flags = append(b.flags, name)
	b.values[name] = def
	b.defaults[name] = def
	b.converters[name] = conv
	return b
}

func (b *flagsBuilder) AddParameterlessFlag(name, desc string, conv Converter, found, missing interface{}) Builder {
	if b.command.err != nil {
		return b
	}

	if !isValidFlagName(name) {
		b.command.err = fmt.Errorf("funcv: invalid flag name %s [arg %d]", name, len(b.command.args)+len(b.flags))
		return b
	}

	if !conv.IsSupported(found) {
		b.command.err = fmt.Errorf("funcv: invalid default %v for flag %s [arg %d]", found, name, len(b.command.args)+len(b.flags))
		return b
	}

	if !conv.IsSupported(missing) {
		b.command.err = fmt.Errorf("funcv: invalid default %v for flag %s [arg %d]", missing, name, len(b.command.args)+len(b.flags))
		return b
	}

	b.desc = append(b.desc, desc)
	b.flags = append(b.flags, name)
	b.values[name] = missing
	b.defaults[name] = missing
	b.founddefs[name] = found
	b.converters[name] = conv
	return b
}

func (b *flagsBuilder) AddConstant(text string, insensitive bool) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddConstant(text, insensitive)
}

func (b *flagsBuilder) AddVariable(name, desc string, conv Converter) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddVariable(name, desc, conv)
}

func (b *flagsBuilder) AddVariableWithDefault(name, desc string, conv Converter, def interface{}) ClosingBuilder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddVariableWithDefault(name, desc, conv, def)
}

func (b *flagsBuilder) AddArgument(arg Argument) Builder {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddArgument(arg)
}

func (b *flagsBuilder) AddVariadic(name, desc string, conv Converter) Compiler {
	if b.command.err != nil {
		return b
	}

	b.command.args = append(b.command.args, b)
	return b.command.AddVariadic(name, desc, conv)
}

func (b *flagsBuilder) Compile() (Command, error) {
	b.command.args = append(b.command.args, b)
	return b.command.Compile()
}

func (b *flagsBuilder) MustCompile() Command {
	b.command.args = append(b.command.args, b)
	return b.command.MustCompile()
}

func (b *flagsBuilder) ToGroup(grp *Group, fn interface{}) error {
	b.command.args = append(b.command.args, b)
	return b.command.ToGroup(grp, fn)
}

func (b *flagsBuilder) WriteTo(w io.Writer) (int64, error) {
	var written int64

	for i, name := range b.flags {
		def, _ := b.defaults[name]

		if n, err := fmt.Fprintf(w, "\n\t%s\t%s (default: %v)", toFlag(name), b.desc[i], def); err == nil {
			written += int64(n)
		} else {
			return written + int64(n), err
		}
	}

	return written, nil
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
