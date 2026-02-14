package rules

import (
	"go/ast"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

// CheckLowercase reports if the log message starts with an uppercase letter,
// and provides a SuggestedFix to lowercase it.
func CheckLowercase(pass *analysis.Pass, call *ast.CallExpr, msg string) {
	if len(msg) == 0 {
		return
	}
	first, size := utf8.DecodeRuneInString(msg)
	if first == utf8.RuneError {
		return
	}
	if !unicode.IsUpper(first) {
		return
	}
	// Prepare diagnostic with suggested fix
	msgArg := call.Args[0] // we know from extractLogMessage that first arg is message
	pos := msgArg.Pos()
	end := msgArg.End()

	// New message with first letter lowercased
	fixedMsg := string(unicode.ToLower(first)) + msg[size:]

	pass.Report(analysis.Diagnostic{
		Pos:     pos,
		End:     end,
		Message: "log message should not start with an uppercase letter",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "Lowercase first letter",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     end,
						NewText: []byte(`"` + fixedMsg + `"`),
					},
				},
			},
		},
	})
}
