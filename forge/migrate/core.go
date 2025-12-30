// Package migrate re-exports the core sub-package
package migrate

import (
	"github.com/forgego/forge/db/migrate/core"
)

// Re-export core package types (already exported via migrate.go, but keeping for completeness)
type (
	// Core types are already re-exported in migrate.go
	_ = core.Change // Ensure core package is imported
)
