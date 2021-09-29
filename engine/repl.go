package engine

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"reflect"
	"strings"
)

// REPL runs a read-evaluate-print-loop. Commands are taken from src and output
// is written to dst. assets is needed for commands like reload.
func (g *Game) REPL(src io.Reader, dst io.Writer, assets fs.FS) error {
	const prompt = "game> "
	fmt.Fprint(dst, prompt)
	sc := bufio.NewScanner(src) // TODO: use a repl library?
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
		case "query":
			g.cmdQuery(dst, argv)
		case "hide":
			g.cmdHide(dst, argv)
		case "show":
			g.cmdShow(dst, argv)
		case "print":
			g.cmdPrint(dst, argv)
		}
		fmt.Fprint(dst, prompt)
	}
	return sc.Err()
}

func (g *Game) cmdSave(dst io.Writer, argv []string) {
	c := g.cmdutilComponentArg1(dst, argv)
	if c == nil {
		return
	}
	s, ok := c.(Saver)
	if !ok {
		fmt.Fprintf(dst, "Component not saveable (type %T)\n", c)
		return
	}
	if err := s.Save(); err != nil {
		fmt.Fprintf(dst, "Couldn't save: %v\n", err)
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
		return
	}
	c := interface{}(g)
	if len(argv) == 2 { // subtree
		id := argv[1]
		c = g.Component(id)
		if c == nil {
			fmt.Fprintf(dst, "Component %q not found\n", id)
			return
		}
	}
	g.printTreeRecursive(dst, 0, c)
}

func (g *Game) printTreeRecursive(dst io.Writer, depth int, c interface{}) {
	indent := ""
	if depth > 0 {
		indent = strings.Repeat("  ", depth-1) + "↳ "
	}
	i, ok := c.(Identifier)
	if ok {
		fmt.Fprintf(dst, "%s%v %q\n", indent, c, i.Ident())
	} else {
		fmt.Fprintf(dst, "%s%v\n", indent, c)
	}
	if sc, ok := c.(Scanner); ok {
		sc.Scan(func(x interface{}) error {
			g.printTreeRecursive(dst, depth+1, x)
			return nil
		})
	}
}

func (g *Game) cmdQuery(dst io.Writer, argv []string) {
	if len(argv) < 2 || len(argv) > 3 {
		fmt.Fprintln(dst, "Usage: query BEHAVIOUR [ANCESTOR_ID]")
		fmt.Fprint(dst, "Behaviours:")
		for _, b := range Behaviours {
			fmt.Fprintf(dst, " %s", b.Name())
		}
		return
	}

	var behaviour reflect.Type
	for _, b := range Behaviours {
		if b.Name() == argv[1] {
			behaviour = b
		}
	}
	if behaviour == nil {
		fmt.Fprintf(dst, "Unknown behaviour %q\n", argv[1])
		return
	}

	var ancestor interface{} = g
	if len(argv) == 3 {
		c := g.Component(argv[2])
		if c == nil {
			fmt.Fprintf(dst, "Component %q not found\n", argv[2])
			return
		}
		ancestor = c
	}

	noResults := true
	g.Query(ancestor, behaviour, func(c interface{}) error {
		noResults = false
		i, ok := c.(Identifier)
		if ok {
			fmt.Fprintf(dst, "%T %q\n", c, i.Ident())
		} else {
			fmt.Fprintf(dst, "%T\n", c)
		}
		return nil
	}, nil)
	if noResults {
		fmt.Fprintln(dst, "No results")
	}
}

func (g *Game) cmdutilComponentArg1(dst io.Writer, argv []string) interface{} {
	if len(argv) != 2 {
		fmt.Fprintln(dst, "Usage: hide ID")
		return nil
	}
	id := argv[1]
	c := g.Component(id)
	if c == nil {
		fmt.Fprintf(dst, "Component %q not found\n", id)
	}
	return c
}

func (g *Game) cmdHide(dst io.Writer, argv []string) {
	c := g.cmdutilComponentArg1(dst, argv)
	if c == nil {
		return
	}
	h, ok := c.(Hider)
	if !ok {
		fmt.Fprintf(dst, "Component not hidable (type %T)\n", c)
		return
	}
	h.Hide()
}

func (g *Game) cmdShow(dst io.Writer, argv []string) {
	c := g.cmdutilComponentArg1(dst, argv)
	if c == nil {
		return
	}
	h, ok := c.(Hider)
	if !ok {
		fmt.Fprintf(dst, "Component not showable (type %T)\n", c)
		return
	}
	h.Show()
}

func (g *Game) cmdPrint(dst io.Writer, argv []string) {
	c := g.cmdutilComponentArg1(dst, argv)
	if c == nil {
		return
	}
	fmt.Fprintf(dst, "%#v\n", c)
}
