package analyzer

import (
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"

	"github.com/SANEKNAYMCHIK/log-linter/internal/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type Config struct {
	SensitiveWords    []string
	SensitivePatterns []string
}

// var Analyzer = &analysis.Analyzer{
// 	Name:     "log-linter",
// 	Doc:      "Check log messages against style rules",
// 	Run:      run,
// 	Requires: []*analysis.Analyzer{inspect.Analyzer},
// }

func NewAnalyzer(cfg *Config) *analysis.Analyzer {
	a := &analysis.Analyzer{
		Name:     "loglint",
		Doc:      "Checks log messages against style rules",
		Run:      run(cfg),
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Flags:    flag.FlagSet{}, // we'll set flags in init or via constructor
	}
	// Register flags for configuration (bonus #1)
	a.Flags.StringVar(&cfgSensitiveWords, "sensitive-words", "password,secret,token,key", "comma-separated sensitive words")
	a.Flags.StringVar(&cfgSensitivePatterns, "sensitive-patterns", "", "comma-separated regex patterns for sensitive data")
	return a
}

var (
	cfgSensitiveWords    string
	cfgSensitivePatterns string
)

// run returns the main analysis function.
func run(cfg *Config) func(pass *analysis.Pass) (interface{}, error) {
	return func(pass *analysis.Pass) (interface{}, error) {
		// Build config from flags if cfg is nil (standalone mode)
		if cfg == nil {
			cfg = &Config{}
			if cfgSensitiveWords != "" {
				cfg.SensitiveWords = strings.Split(cfgSensitiveWords, ",")
			}
			if cfgSensitivePatterns != "" {
				cfg.SensitivePatterns = strings.Split(cfgSensitivePatterns, ",")
			}
		}

		inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

		nodeFilter := []ast.Node{
			(*ast.CallExpr)(nil),
		}

		inspect.Preorder(nodeFilter, func(n ast.Node) {
			call := n.(*ast.CallExpr)

			msg, ok := extractLogMessage(pass, call)
			if !ok {
				return
			}

			// Apply each rule, passing config where needed
			rules.CheckEnglish(pass, call, msg)
			rules.CheckLowercase(pass, call, msg)
			rules.CheckEmoji(pass, call, msg)
			rules.CheckSensitive(pass, call, msg, cfg.SensitiveWords, cfg.SensitivePatterns)
		})

		return nil, nil
	}
}

// extractLogMessage tries to extract the message string from a logging call.
// It supports slog and zap. Returns the message and true if found.
func extractLogMessage(pass *analysis.Pass, call *ast.CallExpr) (string, bool) {
	// We need to identify the function being called.
	// We'll check based on package and method name, and also by type (more robust).
	fun, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	// Get the type of the receiver/package to know what logger it is.
	// This handles cases like slog.Info, logger.Info, zap.L().Info, etc.
	receiverType := pass.TypesInfo.TypeOf(fun.X)
	if receiverType == nil {
		return "", false
	}

	// Helper to check if receiver is from a particular package
	isPkg := func(pkgPath string) bool {
		if named, ok := receiverType.(*types.Named); ok {
			return named.Obj().Pkg().Path() == pkgPath
		}
		// Could also be a pointer type
		if ptr, ok := receiverType.(*types.Pointer); ok {
			if named, ok := ptr.Elem().(*types.Named); ok {
				return named.Obj().Pkg().Path() == pkgPath
			}
		}
		return false
	}

	// Determine the method name
	methodName := fun.Sel.Name

	// Define which method names we consider logging levels
	logLevels := map[string]bool{
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"DPanic": true, "Panic": true, "Fatal": true,
		// slog also has Log(ctx, level, msg, ...)
		"Log": true,
	}
	if !logLevels[methodName] {
		return "", false
	}

	// Check if receiver belongs to slog or zap
	// slog: *log/slog.Logger, log/slog.Logger, or the package function slog.Info (receiverType is *types.Package?)
	// For package-level functions like slog.Info, fun.X is an *ast.Ident with type *types.PkgName.
	if pkg, ok := fun.X.(*ast.Ident); ok {
		// This could be slog.Info (pkg name "slog")
		if pkg.Name == "slog" {
			// slog package function
			return extractStringArg(pass, call, 0)
		}
		if pkg.Name == "zap" {
			// zap.L().Info, zap.Info? Actually zap doesn't have package-level Info, it's on logger.
			// But zap.L() returns a logger, so that's handled by type check below.
		}
	}

	// Check by type
	if isPkg("log/slog") {
		// slog.Logger methods
		return extractStringArg(pass, call, 0)
	}
	if isPkg("go.uber.org/zap") {
		// zap.Logger or zap.SugaredLogger methods
		return extractStringArg(pass, call, 0)
	}
	if isPkg("go.uber.org/zap/zapcore") {
		// zapcore.CheckedEntry? Not needed.
		return "", false
	}

	return "", false
}

// extractStringArg extracts a string literal or constant from the i-th argument.
func extractStringArg(pass *analysis.Pass, call *ast.CallExpr, idx int) (string, bool) {
	if idx >= len(call.Args) {
		return "", false
	}
	arg := call.Args[idx]

	// Direct string literal
	if lit, ok := arg.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		s := strings.Trim(lit.Value, `"`)
		return s, true
	}

	// Ident that refers to a constant
	if ident, ok := arg.(*ast.Ident); ok {
		if obj, ok := pass.TypesInfo.Uses[ident]; ok {
			if con, ok := obj.(*types.Const); ok {
				// Constant value may be typed as string
				if val := con.Val().String(); val != "" {
					// It's stored with quotes, so unquote
					if unquoted, err := strconv.Unquote(val); err == nil {
						return unquoted, true
					}
				}
			}
		}
	}

	// Could also be a call to fmt.Sprint? Not handling for simplicity.
	return "", false
}
