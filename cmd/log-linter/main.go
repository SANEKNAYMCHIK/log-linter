package main

import (
	"github.com/SANEKNAYMCHIK/log-linter/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	a := analyzer.NewAnalyzer(nil)
	singlechecker.Main(a)
}
