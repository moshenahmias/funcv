package funcv

import (
	"fmt"
	"io"
)

type bindedCommand struct {
	cmd Command
	fn  interface{}
}

type group struct {
	commands []bindedCommand
}

func (g *group) Add(cmd Command, fn interface{}) Group {
	g.commands = append(g.commands, bindedCommand{cmd, fn})
	return g
}

func (g *group) Execute(args []string) int {
	n := 0

	for _, p := range g.commands {
		if err := p.cmd.Execute(args, p.fn); err == nil {
			n++
		}
	}

	return n
}

func (g *group) WriteTo(w io.Writer) (int64, error) {
	var written int64

	for i, p := range g.commands {
		if n, err := p.cmd.WriteTo(w); err == nil {
			written += n
		} else {
			return written + n, err
		}

		if i+1 < len(g.commands) {
			if n, err := fmt.Fprint(w, "\n\n"); err == nil {
				written += int64(n)
			} else {
				return written + int64(n), err
			}
		}
	}

	return written, nil
}
