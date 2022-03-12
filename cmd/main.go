package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/luoxiaojun1992/go-missing-len/pkg"
	"go/format"
	"go/parser"
	"go/token"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "")

	var resultFormat string
	flag.StringVar(&resultFormat, "format", "", "")

	flag.Parse()

	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		panic(err)
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
