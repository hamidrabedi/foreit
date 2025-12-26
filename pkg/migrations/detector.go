// Package migrations provides migration generation and management.
package migrations

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/detection"
)

// Detector detects changes between current and previous model states
// This is a compatibility wrapper around the new detection package.
type Detector struct {
	detector detection.ChangeDetector
}

// NewDetector creates a new change detector
func NewDetector() *Detector {
	return &Detector{
		detector: detection.NewDetector(),
	}
}

// DetectChanges compares current models to previous state and returns all changes
func (d *Detector) DetectChanges(current, previous []*generator.ModelDefinition) ([]Change, error) {
	changes, err := d.detector.DetectChanges(current, previous)
	if err != nil {
		return nil, err
	}
	// Change is a type alias for core.Change, so we can return directly
	return changes, nil
}
