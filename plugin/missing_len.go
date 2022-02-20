package main

import (
	"github.com/luoxiaojun1992/go-missing-len/pkg"
	"golang.org/x/tools/go/analysis"
)

type analyzerPlugin struct{}

// This must be implemented
func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		pkg.NewAnalyzer(),
	}
}

// This must be defined and named 'AnalyzerPlugin'
var AnalyzerPlugin analyzerPlugin
