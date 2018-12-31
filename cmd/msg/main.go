package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"
)

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	path, err := os.Getwd()
	chk(err)
	println("parsing package:" + path)

	target := os.Args[1]
	println("Searching for: " + target)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, func(fi os.FileInfo) bool {
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.ParseComments)
	chk(err)
	for _, v := range pkgs {
		processPackage(fset, v, map[string]bool{
			target: true,
		})
	}
}

type walker struct {
	fset *token.FileSet             // active fileset
	imps map[string]string          // import package name to full gopath
	extz map[string]map[string]bool // external types needed
	file map[string]bool            // external packages we reference
	need map[string]bool            // things we need
	have map[string]ast.Node        // desired result nodes
	seen map[string]ast.Node        // nodes walked
}

func (w *walker) Visit(n ast.Node) ast.Visitor {
	switch p := n.(type) {

	case *ast.ImportSpec: // find imports in case we see external packages
		clean := strings.Trim(p.Path.Value, "\"")
		var name string
		if p.Name == nil {
			_, name = path.Split(clean)
		} else {
			name = p.Name.Name
		}
		w.imps[name] = clean
		w.extz[name] = make(map[string]bool)

	case *ast.TypeSpec: // record tupes as we run across them
		name := p.Name.Name
		if w.need[name] {
			w.have[name] = n
			fmt.Printf("Start: %q\n", name)
			p.Name = nil // don't care about names
			return w
		}
		w.seen[name] = n

	case *ast.Field:
		p.Names = nil // don't care about names
		return w

	case *ast.SelectorExpr:
		// TODO: add file for parser to process
		if p.Sel.IsExported() {
			fmt.Printf("field-ext: \"%s.%s\"\n", p.X, p.Sel.Name)
			w.extz[fmt.Sprint(p.X)][p.Sel.Name] = true
		}
		return nil

	case *ast.Ident:
		if p == nil {
			return nil // skipped by parent (usually a name)
		}
		if p.IsExported() {
			fmt.Printf("field: %q\n", p.Name)

			// Found a new type that needs to be serialized
			w.need[p.Name] = true
			if o, ok := w.seen[p.Name]; ok {
				w.have[p.Name] = o
			}
		}
		return nil

	}
	return w
}

func processPackage(fset *token.FileSet, pkg *ast.Package, need map[string]bool) {
	w := &walker{
		fset: fset,
		extz: make(map[string]map[string]bool),
		imps: make(map[string]string),
		need: need,
		have: make(map[string]ast.Node),
	}
	println(pkg.Name)
	ast.Walk(w, pkg)
	fmt.Println()

	fmt.Printf("%#v\n", w.imps)
	fmt.Println()

	// External imports
	for k, v := range w.extz {
		fmt.Printf("import: %q : %#v\n", k, v)
	}
	fmt.Println()

	// ast.Print(fset, pkg)
	for k, v := range w.have {
		fmt.Printf("found %q\n", k)
		_ = v
		// ast.Print(fset, v)
	}
}
