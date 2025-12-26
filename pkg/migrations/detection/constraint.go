package detection

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// detectConstraintChanges detects constraint changes
func (d *Detector) detectConstraintChanges(tableName string, current, previous []generator.ConstraintDefinition) ([]core.Change, error) {
	var changes []core.Change

	currentMap := make(map[string]generator.ConstraintDefinition)
	previousMap := make(map[string]generator.ConstraintDefinition)

	for _, constr := range current {
		currentMap[constr.Name] = constr
	}

	for _, constr := range previous {
		previousMap[constr.Name] = constr
	}

	// Detect new constraints
	for name, constr := range currentMap {
		if _, exists := previousMap[name]; !exists {
			changes = append(changes, &core.AddConstraint{
				Table:      tableName,
				Constraint: constr,
			})
		}
	}

	// Detect dropped constraints
	for name := range previousMap {
		if _, exists := currentMap[name]; !exists {
			changes = append(changes, &core.DropConstraint{
				Table:          tableName,
				ConstraintName: name,
			})
		}
	}

	return changes, nil
}

