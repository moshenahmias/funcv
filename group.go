package funcv

import (
	"strings"
)

type pair struct {
	cmd Command
	fn  interface{}
}

type group struct {
	commands []pair
}

func (g *group) Add(cmd Command, fn interface{}) Group {
	g.commands = append(g.commands, pair{cmd, fn})
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

func (g *group) String() string {
	if len(g.commands) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, p := range g.commands {
		sb.WriteString(p.cmd.String())

		if i+1 < len(g.commands) {
			sb.WriteString(newline)
			sb.WriteString(newline)
		}
	}

	return sb.String()
}
