package migrations

import (
	"github.com/forgego/forge/pkg/migrations/linting"
)

// Linter is an alias for the linting.Linter type
type Linter = linting.Linter

// LintResult is an alias for the linting.LintResult type
type LintResult = linting.LintResult

// NewLinter creates a new migration linter
func NewLinter() *Linter {
	return linting.NewLinter()
}

