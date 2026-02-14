package rules

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// CheckSensitive reports if the log message contains sensitive words or matches patterns.
func CheckSensitive(pass *analysis.Pass, call *ast.CallExpr, msg string, words []string, patterns []string) {
	lowerMsg := strings.ToLower(msg)

	// Check simple words
	for _, w := range words {
		if strings.Contains(lowerMsg, strings.ToLower(w)) {
			pass.Reportf(call.Pos(), "log message may contain sensitive data (word %q)", w)
			return
		}
	}

	// Check regex patterns
	for _, pat := range patterns {
		if pat == "" {
			continue
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			// invalid pattern, skip? Could report as internal error, but ignore.
			continue
		}
		if re.MatchString(msg) {
			pass.Reportf(call.Pos(), "log message may contain sensitive data (pattern %q)", pat)
			return
		}
	}
}
