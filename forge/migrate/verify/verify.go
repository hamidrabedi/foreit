// Package verify re-exports the db/migrate/verify package
package verify

import (
	"github.com/forgego/forge/db/migrate/verify"
)

// Re-export types
type (
	Linter       = verify.Linter
	LintResult   = verify.LintResult
	LintRule     = verify.LintRule
	DriftDetector = verify.DriftDetector
	Drift        = verify.Drift
)

// Re-export functions
var (
	NewLinter       = verify.NewLinter
	NewDriftDetector = verify.NewDriftDetector
)
