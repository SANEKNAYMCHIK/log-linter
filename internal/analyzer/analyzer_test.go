package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	a := NewAnalyzer(nil)
	analysistest.Run(t, analysistest.TestData(), a, "slog")
}
