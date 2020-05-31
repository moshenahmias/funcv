package funcv

import (
	"fmt"
	"io"
)

// Pair of a command and an action function
type Pair struct {
	cmd Command
	fn  interface{}
}

// Group of commands with binded action functions
type Group []Pair

// Add a command and an action function to the group
func (g *Group) Add(cmd Command, fn interface{}) *Group {
	*g = append(*g, Pair{cmd, fn})
	return g
}

// Execute tests the supplied arguments against all commands
// in the group, if some are compatible, the paired action function
// for each command is called with the extracted parameters,
// the number of called functions is returned, to be called, every
// function needs to be compatible with the command's arguments
func (g *Group) Execute(args []string) int {
	n := 0

	for _, p := range *g {
		if _, err := p.cmd.Execute(args, p.fn); err == nil {
			n++
		}
	}

	return n
}

// WriteTo will write to the writer an informative usage
// text about the commands in the group
func (g *Group) WriteTo(w io.Writer) (int64, error) {
	var written int64

	for i, p := range *g {
		if n, err := p.cmd.WriteTo(w); err == nil {
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
