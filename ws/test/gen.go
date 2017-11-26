// +build ignore

package main

import (
	"fmt"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// references:
// - https://github.com/gopherjs/gopherjs/blob/master/tool.go
// - https://github.com/kisielk/gotool/blob/master/match.go
// - https://golang.org/src/cmd/go/internal/work/build.go
// - https://golang.org/src/go/importer/importer.go
// - https://golang.org/pkg/go/build/#Package

func main() {
	println("Default Usage: go run ./ws/test/gen.go github.com/bign8/msg/ws")

	// print("importer")
	// out, err := importer.Default().Import("github.com/bign8/msg/ws")
	// if err != nil {
	// 	panic(err)
	// }
	// println(out.Name(), out.Path())
	// for i, name := range out.Scope().Names() {
	// 	x := out.Scope().Lookup(name)
	// 	println(i, name, x.Exported())
	// }

	println("building path")
	path := filepath.Join(build.Default.GOPATH, "src", os.Args[1])
	println(path)

	println("parser")
	fset := token.NewFileSet()
	f, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	fmt.Printf("FSET: %#v", f)

	println("builder")
	bout, err := build.Default.Import(os.Args[1], "", 0)
	if err != nil {
		panic(err)
	}
	for _, name := range bout.TestGoFiles {
		println(name)
		f, err := parser.ParseFile(fset, filepath.Join(path, name), nil, parser.ParseComments)
		// pkg, err := importer.For("source", nil).Import("github.com/bign8/msg/ws/" + name)
		if err != nil {
			panic(err)
		}
		println(f.Scope.String())
	}

}
