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

// DefultedVariablesAdder is used to add variables with default to a command
type DefultedVariablesAdder interface {
	// AddStrVarWithDefault adds a named string variable with
	// a default value and a description
	AddStrVarWithDefault(name, def, desc string) DefultedVariablesBuilder

	// AddIntVarWithDefault adds a named integer variable with
	// a default value and a description
	AddIntVarWithDefault(name string, def, base int, desc string) DefultedVariablesBuilder
}

// ConstantAdder is used to add a constant to a command
type ConstantAdder interface {
	// AddConstant argument, constants are a command unique
	// identifiers and can be added anywhere
	AddConstant(text string, insensitive bool) Builder
}

// VariablesAdder is used to add variables to a command
type VariablesAdder interface {
	// AddStrVar adds a named string variable with
	// a description
	AddStrVar(name, desc string) Builder
	// AddIntVar adds a named interger variable with
	// a description
	AddIntVar(name string, base int, desc string) Builder
}

// FlagsAdder is used to add flags to a command
type FlagsAdder interface {
	// AddStrFlag adds a string flag with
	// a default value and a description
	// ex: -s abcd, --str abcd
	AddStrFlag(name, def, desc string) FlagsBuilder
	// AddIntFlag adds an integer flag with
	// a default value, base and a description
	// ex: -i 123, --num 456
	AddIntFlag(name string, def, base int, desc string) FlagsBuilder
	// AddBoolFlag adds a boolean flag with
	// a description, the default value is false (flag not found)
	// ex: -b, -b false, -b true
	AddBoolFlag(name, desc string) FlagsBuilder
}

// ArgAdder is used to add custom arguments to a command
type ArgAdder interface {
	// AddArg can be used for adding a custom argument
	// to the command
	AddArg(arg Arg) Builder
}

// DefultedVariablesBuilder is a subset builder
// for variable arguments with defaults
type DefultedVariablesBuilder interface {
	DefultedVariablesAdder
	Compiler
}

// FlagsBuilder is a subset builder for flag arguments
type FlagsBuilder interface {
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
	Compiler
}

// Builder is used for building a new Command
type Builder interface {
	ArgAdder
	FlagsAdder
	ConstantAdder
	VariablesAdder
	DefultedVariablesAdder
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
	// command, if they are all compatible, the action function is
	// called with the extracted parameters and (len(args), nil) is
	// returned, else, non-nil error is returned with the number of
	// compatible arguments. the action function arguments
	// need to be compatible with the command's arguments or else
	// a non-nil error is returned
	Execute(args []string, fn interface{}) (int, error)
	io.WriterTo
}

/*
// Group ties together multiple commands
type Group interface {
	// Add a command and an action function
	// to the group
	Add(cmd Command, fn interface{}) Group
	// Execute tests the supplied arguments against all commands
	// in the group, if some are compatible, the paired action function
	// for each command is called with the extracted parameters,
	// the number of called functions is returned, to be called, every
	// function needs to be compatible with the command's arguments
	Execute(args []string) int
	io.WriterTo
}
*/

// Arg represents a command argument
type Arg interface {
	// Extract an argument(s) from the list of arguments
	// and return the rest of the arguments and extracted
	// parameters, returns a non-nil error if the extraction
	// failed for any reason
	Extract(args []string) ([]string, []interface{}, error)
	io.WriterTo
	fmt.Stringer
}

// Converter of text argument to typed value
type Converter interface {
	Convert(arg string) (interface{}, error)
}

// NewCommand returns a builder that is used for
// building a new command
func NewCommand(desc string) Builder {
	return &command{desc: desc}
}
