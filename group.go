package funcv

import (
	"fmt"
	"io"
	"reflect"
)

// Pair of a command and an action function
type Pair struct {
	Cmd Command
	Fn  interface{}
}

// Group of commands with binded action functions
type Group []Pair

// Add a command and an action function to the group
func (g *Group) Add(cmd Command, fn interface{}) *Group {
	*g = append(*g, Pair{cmd, fn})
	return g
}

// Call the i's action function with the given parameters
// invalid input will panic
func (g *Group) Call(i int, params ...interface{}) {
	var in []reflect.Value

	for _, param := range params {
		in = append(in, reflect.ValueOf(param))
	}

	reflect.ValueOf((*g)[i].Fn).Call(in)
}

// ExecuteAll tests the supplied arguments against all commands
// in the group, if a suitable command found, the paired action
// function is called with the extracted parameters, the number of
// called functions is returned
func (g *Group) ExecuteAll(args []string) (n int) {
	for _, p := range *g {
		if _, err := p.Cmd.Execute(args, p.Fn); err == nil {
			n++
		}
	}

	return
}

// ExecuteFirst tests the supplied arguments against the commands
// in the group, if a suitable command found, the paired action
// function is called with the extracted parameters and the method
// returns immediately the command's index, without testing other
// commands, if no suitable command found, the method returns a
// negative value
func (g *Group) ExecuteFirst(args []string) (i int) {
	var p Pair

	for i, p = range *g {
		if _, err := p.Cmd.Execute(args, p.Fn); err == nil {
			return
		}
	}

	return -1
}

// WriteTo will write to the writer an informative usage
// text about the commands in the group
func (g *Group) WriteTo(w io.Writer) (int64, error) {
	var written int64

	for i, p := range *g {
		if n, err := p.Cmd.WriteTo(w); err == nil {
			written += n
		} else {
			return written + n, err
		}

		if i+1 < len(*g) {
			if n, err := fmt.Fprint(w, "\n\n"); err == nil {
				written += int64(n)
			} else {
				return written + int64(n), err
			}
		}
	}

	return written, nil
}
