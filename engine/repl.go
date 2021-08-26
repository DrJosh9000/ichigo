package engine

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

// REPL runs a read-evaluate-print-loop. Commands are taken from src and output
// is written to dst. assets is needed for commands like reload.
func (g *Game) REPL(src io.Reader, dst io.Writer, assets fs.FS) error {
	fmt.Fprint(dst, "game>")
	sc := bufio.NewScanner(src)
	for sc.Scan() {
		tok := strings.Split(sc.Text(), " ")
		if len(tok) == 0 {
			continue
		}
		switch tok[0] {
		case "quit":
			os.Exit(0)
		case "pause":
			g.Disable()
		case "resume", "unpause":
			g.Enable()
		case "save":
			if len(tok) != 2 {
				fmt.Fprintln(dst, "Usage: save ID")
				break
			}
			id := tok[1]
			c := g.Component(id)
			if c == nil {
				fmt.Fprintf(dst, "Component %q not found\n", id)
				break
			}
			s, ok := c.(Saver)
			if !ok {
				fmt.Fprintf(dst, "Component %q not a Saver (type %T)\n", id, c)
				break
			}
			if err := s.Save(); err != nil {
				fmt.Fprintf(dst, "Couldn't save %q: %v\n", id, err)
			}
		case "reload":
			g.Disable()
			g.Hide()
			if err := g.Load(assets); err != nil {
				fmt.Fprintf(dst, "Couldn't load: %v\n", err)
				break
			}
			g.Prepare()
			g.Enable()
			g.Show()
		case "tree":
			var c interface{} = g
			if len(tok) == 2 {
				// subtree
				id := tok[1]
				c = g.Component(id)
				if c == nil {
					fmt.Fprintf(dst, "Component %q not found\n", id)
					break
				}
			}
			type item struct {
				c     interface{}
				depth int
			}
			stack := []item{{c, 0}}
			for len(stack) > 0 {
				tail := len(stack) - 1
				x := stack[tail]
				stack = stack[:tail]
				c := x.c

				indent := ""
				if x.depth > 0 {
					indent = strings.Repeat("  ", x.depth-1) + "↳ "
				}
				i, ok := c.(Identifier)
				if ok {
					fmt.Fprintf(dst, "%s%T %s\n", indent, c, i.Ident())
				} else {
					fmt.Fprintf(dst, "%s%T\n", indent, c)
				}

				if s, ok := c.(Scanner); ok {
					for _, y := range s.Scan() {
						stack = append(stack, item{y, x.depth + 1})
					}
				}
			}
		}
		fmt.Fprint(dst, "game>")
	}
	return sc.Err()
}
