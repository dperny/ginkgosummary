package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "requires a file")
		os.Exit(1)
	}
	fset := token.NewFileSet()

	tree, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	ast.Walk(Visitor{0}, tree)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing go file: %v\n", err)
		os.Exit(1)
	}

}

type Visitor struct {
	depth int
}

func (v Visitor) Visit(node ast.Node) ast.Visitor {
	tabs := ""
	for i := 0; i < v.depth; i++ {
		tabs = tabs + "    "
	}
	inBlock := false
	switch n := node.(type) {
	case *ast.AssignStmt:
		for _, right := range n.Rhs {
			if call, ok := right.(*ast.CallExpr); ok {
				if ident, ok := call.Fun.(*ast.Ident); ok {
					switch ident.Name {
					case "Describe", "Context":
						if lit, ok := call.Args[0].(*ast.BasicLit); ok {
							fmt.Printf("%s%s\n", tabs, strings.Trim(lit.Value, "\""))
							inBlock = true
						}
					}
				}
			}
		}
	case *ast.CallExpr:
		if ident, ok := n.Fun.(*ast.Ident); ok {
			dash := ""
			switch ident.Name {
			case "It":
				dash = "- "
				fallthrough
			case "Describe", "Context":
				if lit, ok := n.Args[0].(*ast.BasicLit); ok {
					fmt.Printf("%s%s%s\n", tabs, dash, strings.Trim(lit.Value, "\""))
					inBlock = true
				}
			}
		}
	}
	d := v.depth
	if inBlock {
		d = d + 1
	}
	return Visitor{depth: d}
}
