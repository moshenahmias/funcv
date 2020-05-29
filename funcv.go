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
	AddDefStrVar(name, def, desc string) DefultedVariablesBuilder
	AddDefIntVar(name string, def, base int, desc string) DefultedVariablesBuilder
}

type ConstantAdder interface {
	AddConst(text string, insensitive bool) Builder
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

type DefultedVariablesBuilder interface {
	DefultedVariablesAdder
	CommandCreator
}

type FlagsBuilder interface {
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
	CommandCreator
}

type Builder interface {
	AddArg(arg Arg) Builder
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
	CommandCreator
}

type CommandCreator interface {
	Create() Command
}

type Command interface {
	Execute(args []string, fn interface{}) error
	fmt.Stringer
}

func Build(exe, desc string) Builder {
	return &command{exe: exe, desc: desc}
}

type Arg interface {
	Extract(args []string) ([]string, []interface{}, error)
	Description() string
	fmt.Stringer
}

type Converter interface {
	Convert(arg string) (interface{}, error)
}
