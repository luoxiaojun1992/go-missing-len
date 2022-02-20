package main

import (
	"flag"
	"fmt"
	"github.com/luoxiaojun1992/go-missing-len/pkg"
	"go/parser"
	"go/token"
)

func main() {
	var filename string
	flag.StringVar(&filename, "file", "", "")
	flag.Parse()

	file, err := parser.ParseFile(token.NewFileSet(), filename, nil, 0)
	if err != nil {
		panic(err)
	}
	linter := pkg.NewMissingLenLinter()
	linter.Check(file)
	if len(linter.Hints) > 0 {
		for _, hint := range linter.Hints {
			fmt.Printf("Pos: %d, End: %d, Category: %s, Message: %s, Suggestion: %s \n", hint.Pos, hint.End, hint.Category, hint.Message, hint.Suggestion)
		}
	}
}
