package rules

import (
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// CheckEnglish reports if the log message contains non-Latin characters.
func CheckEnglish(pass *analysis.Pass, call *ast.CallExpr, msg string) {
	for _, r := range msg {
		// Allow Latin characters, common punctuation, and spaces.
		// We use unicode.In to check ranges.
		if !unicode.In(r, unicode.Latin, unicode.Common, unicode.Inherited) {
			pass.Reportf(call.Pos(), "log message contains non-Latin character: %q", r)
			return // report only first violation to avoid spam
		}
	}
}
