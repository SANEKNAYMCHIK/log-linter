package rules

import (
	"go/ast"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// CheckEmoji reports if the log message contains emoji or other special symbols.
// We consider Emoji and Symbol categories.
func CheckEmoji(pass *analysis.Pass, call *ast.CallExpr, msg string) {
	for _, r := range msg {
		if unicode.In(r, unicode.Letter, unicode.Digit, unicode.Space) {
			continue
		}
		pass.Reportf(call.Pos(), "log message contains disallowed character or emoji: %q", r)
		return
	}
}
