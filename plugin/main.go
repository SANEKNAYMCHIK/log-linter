//go:build ignore

package main

import (
	"github.com/SANEKNAYMCHIK/log-linter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

var AnalyzerPlugin *analysis.Analyzer = analyzer.NewAnalyzer(nil)
