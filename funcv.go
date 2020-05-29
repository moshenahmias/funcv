package funcv

import (
	"errors"
	"fmt"
)

var (
	ErrNoArguments  = errors.New("funcv: no arguments")
	ErrArgNotFound  = errors.New("funcv: arg not found")
	ErrUnknownArgs  = errors.New("funcv: unknown arguments")
	ErrInvalidValue = errors.New("funcv: invalid value")
)

type DefultedVariablesAdder interface {
	AddStrVarWithDefault(name, def, desc string) DefultedVariablesBuilder
	AddIntVarWithDefault(name string, def, base int, desc string) DefultedVariablesBuilder
}

type ConstantAdder interface {
	AddConstant(text string, insensitive bool) Builder
}

type VariablesAdder interface {
	AddStrVar(name, desc string) Builder
	AddIntVar(name string, base int, desc string) Builder
}

type FlagsAdder interface {
	AddStrFlag(name, def, desc string) FlagsBuilder
	AddIntFlag(name string, def, base int, desc string) FlagsBuilder
	AddBoolFlag(name, desc string) FlagsBuilder
}

type ArgAdder interface {
	AddArg(arg Arg) Builder
}

type DefultedVariablesBuilder interface {
	DefultedVariablesAdder
	Compiler
}

type FlagsBuilder interface {
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
	Compiler
}

type Builder interface {
	ArgAdder
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
	Compiler
}

type Compiler interface {
	Compile() (Command, error)
	MustCompile() Command
	ToGroup(grp Group, fn interface{}) Command
}

type Command interface {
	Execute(args []string, fn interface{}) error
	fmt.Stringer
}

type Group interface {
	Add(cmd Command, fn interface{}) Group
	Execute(args []string) int
	fmt.Stringer
}

type Arg interface {
	Extract(args []string) ([]string, []interface{}, error)
	Description() string
	fmt.Stringer
}

type Converter interface {
	Convert(arg string) (interface{}, error)
}

func NewCommand(desc string) Builder {
	return &command{desc: desc}
}

func NewGroup() Group {
	return new(group)
}
