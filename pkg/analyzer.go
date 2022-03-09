package pkg

import (
	"golang.org/x/tools/go/analysis"
)

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "missinglen",
		Doc:  "finds slice or map without init len",
		Run:  RunAnalyzer,
	}
}

func RunAnalyzer(pass *analysis.Pass) (interface{}, error) {
	linter := NewMissingLenLinter()
	for _, file := range pass.Files {
		linter.Check(file)
		for _, hint := range linter.Hints {
			pass.Report(analysis.Diagnostic{
				Pos:      hint.Pos,
				End:      hint.End,
				Category: hint.Category,
				Message:  hint.Message,
				SuggestedFixes: []analysis.SuggestedFix{{
					Message: hint.Suggestion,
				}},
			})
		}
		linter.Reset()
	}

	return nil, nil
}
