package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/luoxiaojun1992/go-missing-len/pkg"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "")

	var resultFormat string
	flag.StringVar(&resultFormat, "format", "", "")

	var showAst bool
	flag.BoolVar(&showAst, "ast", false, "")

	flag.Parse()

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, filename, nil, 0)
	if err != nil {
		panic(err)
	}

	if showAst {
		fmt.Println("Ast:")
		fmt.Println()
		if err := ast.Print(fileSet, file); err != nil {
			panic(err)
		}
		fmt.Println()
	}

	linter := pkg.NewMissingLenLinter()
	linter.Check(file)

	fmt.Println("Result:")
	fmt.Println()
	pkg.SerializeHints(linter.Hints, resultFormat)
	fmt.Println()

	fmt.Println("Suggested code:")
	fmt.Println()
	rawCode := bytes.NewBuffer(nil)
	err = format.Node(rawCode, token.NewFileSet(), file)
	if err != nil {
		panic(err)
	}
	fmt.Println(rawCode.String())
	fmt.Println()
}
