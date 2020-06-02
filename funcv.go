// funcv helps you create CLI tools with Go. It offers a different
// approach for dealing with command line arguments and flags.
// funcv supplies an easy to use command builder, you use that builder
// to build your set of commands, each such command can be tested against
// a slice of string arguments, if the arguments are compatible with the
// command, a given action function is called, the parameters for that
// function are the extracted and parsed variables and flags input values.

package funcv

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNoArguments to compile
	ErrNoArguments = errors.New("funcv: no arguments")

	// ErrArgNotFound in supplied arguments
	ErrArgNotFound = errors.New("funcv: arg not found")

	// ErrUnknownArgs in supplied arguments
	ErrUnknownArgs = errors.New("funcv: unknown arguments")

	// ErrInvalidValue in supplied arguments
	ErrInvalidValue = errors.New("funcv: invalid value")
)

// ConstantAdder is used to add a constant to a command, constants
// are a command unique identifiers and can be added anywhere
type ConstantAdder interface {
	AddConstant(text string, insensitive bool) Builder
}

// VariableAdder adds a variable to the command
type VariableAdder interface {
	AddVariable(name, desc string, conv Converter) Builder
}

// DefaultVariableAdder adds a variable with a default value to the command
type DefaultVariableAdder interface {
	AddVariableWithDefault(name, desc string, conv Converter, def interface{}) ClosingBuilder
}

// VariadicAdder adds a variadic parameter to the command
type VariadicAdder interface {
	AddVariadic(name, desc string, conv Converter) Compiler
}

// FlagAdder adds a flag to the command
type FlagAdder interface {
	// AddFlag adds a flag that require a parameter
	AddFlag(name, desc string, conv Converter, def interface{}) Builder
	// AddParameterlessFlag adds a flag that doesn't require a parameter (like boolean flags)
	AddParameterlessFlag(name, desc string, conv Converter, found, missing interface{}) Builder
}

// ArgumentAdder can be used to add any custom argument
// to the command
type ArgumentAdder interface {
	AddArgument(arg Argument) Builder
}

// ClosingBuilder is used for adding optional variables
// or variadic arguments
type ClosingBuilder interface {
	DefaultVariableAdder
	VariadicAdder
	Compiler
}

// Builder is used for building a new Command
type Builder interface {
	ArgumentAdder
	FlagAdder
	ConstantAdder
	VariableAdder
	DefaultVariableAdder
	VariadicAdder
	Compiler
}

// Compiler is used to compile and return
// a new Command
type Compiler interface {
	// Compile and return a Command or an
	// error if the compilation failed
	Compile() (Command, error)
	// MustCompile is the same as Compile
	// but will panic if the compilation failed
	MustCompile() Command
	// ToGroup compiles and adds the command and
	// the given action function to a group, returns
	// an error if the compilation failed
	ToGroup(grp *Group, fn interface{}) error
}

// Command represents a textual command that can be later
// tested against a list of text arguments with an action
// function to run
type Command interface {
	// Execute tests the supplied arguments against the
	// command, if they are all valid, the action function is
	// called with the extracted parameters and (len(args), nil) is
	// returned, else, non-nil error is returned with the number of
	// valid arguments. the action function arguments
	// need to be compatible with the command's arguments or else
	// a non-nil error is returned
	Execute(args []string, fn interface{}) (int, error)
	io.WriterTo
}

// Argument represents a command argument
type Argument interface {
	// Extract an argument(s) from the list of arguments
	// and return the rest of the arguments and extracted
	// parameters, returns a non-nil error if the extraction
	// failed for any reason
	Extract(args []string) ([]string, []interface{}, error)
	io.WriterTo
	fmt.Stringer
}

// Converter of arguments to specific typed values
type Converter interface {
	// Convert the given text argument into a
	// specific type
	Convert(arg string) (interface{}, error)
	// IsSupported returns true if the given value
	// is convertable to the types the converter
	// produces
	IsSupported(v interface{}) bool
}

// NewCommand returns a builder that is used for
// building a new command
func NewCommand(desc string) Builder {
	return &command{desc: desc}
}
