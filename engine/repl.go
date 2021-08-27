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
		argv := strings.Split(sc.Text(), " ")
		if len(argv) == 0 {
			continue
		}
		switch argv[0] {
		case "quit":
			os.Exit(0)
		case "pause":
			g.Disable()
		case "resume", "unpause":
			g.Enable()
		case "save":
			g.cmdSave(dst, argv)
		case "reload":
			g.cmdReload(dst, assets)
		case "tree":
			g.cmdTree(dst, argv)
		}
		fmt.Fprint(dst, "game>")
	}
	return sc.Err()
}

func (g *Game) cmdSave(dst io.Writer, argv []string) {
	if len(argv) != 2 {
		fmt.Fprintln(dst, "Usage: save ID")
		return
	}
	id := argv[1]
	c := g.Component(id)
	if c == nil {
		fmt.Fprintf(dst, "Component %q not found\n", id)
		return
	}
	s, ok := c.(Saver)
	if !ok {
		fmt.Fprintf(dst, "Component %q not a Saver (type %T)\n", id, c)
		return
	}
	if err := s.Save(); err != nil {
		fmt.Fprintf(dst, "Couldn't save %q: %v\n", id, err)
	}
}

func (g *Game) cmdReload(dst io.Writer, assets fs.FS) {
	g.Disable()
	g.Hide()
	if err := g.LoadAndPrepare(assets); err != nil {
		fmt.Fprintf(dst, "Couldn't load: %v\n", err)
		return
	}
	g.Enable()
	g.Show()
}

func (g *Game) cmdTree(dst io.Writer, argv []string) {
	if len(argv) < 1 || len(argv) > 2 {
		fmt.Println(dst, "Usage: tree [ID]")
	}
	var c interface{} = g
	if len(argv) == 2 { // subtree
		id := argv[1]
		c = g.Component(id)
		if c == nil {
			fmt.Fprintf(dst, "Component %q not found\n", id)
			return
		}
	}
	Walk(c, func(c interface{}, p []interface{}) error {
		indent := ""
		if len(p) > 0 {
			indent = strings.Repeat("  ", len(p)-1) + "↳ "
		}
		i, ok := c.(Identifier)
		if ok {
			fmt.Fprintf(dst, "%s%T %s\n", indent, c, i.Ident())
		} else {
			fmt.Fprintf(dst, "%s%T\n", indent, c)
		}
		return nil
	})
}
