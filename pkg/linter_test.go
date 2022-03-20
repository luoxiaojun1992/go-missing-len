package pkg_test

import (
	"github.com/luoxiaojun1992/go-missing-len/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go/parser"
	"go/token"
)

var _ = Describe("Linter", func() {
	Describe("Check", func() {
		Context("Check", func() {
			It("Check", func() {
				fileSet := token.NewFileSet()
				file, err := parser.ParseFile(fileSet, "./../testdata/sample.go", nil, 0)
				if err != nil {
					panic(err)
				}
				linter := pkg.NewMissingLenLinter()
				linter.Check(file)
				Expect(len(linter.Hints), Equal(4))
			})
		})
	})
})
