package detection

import (
	"fmt"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// detectForeignKeyChanges detects foreign key changes
func (d *Detector) detectForeignKeyChanges(tableName string, current, previous []generator.RelationDefinition, currentDef, previousDef *generator.ModelDefinition, allDefs []*generator.ModelDefinition) ([]core.Change, error) {
	var changes []core.Change

	currentMap := make(map[string]generator.RelationDefinition)
	previousMap := make(map[string]generator.RelationDefinition)

	for _, rel := range current {
		if rel.Type == "ForeignKey" || rel.Type == "OneToOne" {
			currentMap[rel.Name] = rel
		}
	}

	for _, rel := range previous {
		if rel.Type == "ForeignKey" || rel.Type == "OneToOne" {
			previousMap[rel.Name] = rel
		}
	}

	// Find target table names using all model definitions
	targetTableMap := make(map[string]string)
	for _, rel := range current {
		if rel.Type == "ForeignKey" || rel.Type == "OneToOne" {
			targetTable := findTargetTable(rel.To, allDefs)
			if targetTable != "" {
				targetTableMap[rel.Name] = targetTable
			}
		}
	}

	// Detect new foreign keys
	for name, rel := range currentMap {
		// Validate that the relation column exists in the table's fields
		hasColumn := false
		for _, field := range currentDef.Fields {
			if field.Name == rel.Name {
				hasColumn = true
				break
			}
		}
		if !hasColumn {
			// Skip FK if column doesn't exist - this prevents broken foreign keys
			continue
		}

		if _, exists := previousMap[name]; !exists {
			targetTable := targetTableMap[name]
			if targetTable == "" {
				// Skip if target table not found
				continue
			}
			changes = append(changes, &core.AddForeignKey{
				Table:       tableName,
				Relation:    rel,
				TargetTable: targetTable,
			})
		} else {
			// Check if foreign key changed
			prevRel := previousMap[name]
			if d.relationChanged(rel, prevRel) {
				targetTable := targetTableMap[name]
				if targetTable == "" {
					// Skip if target table not found
					continue
				}
				changes = append(changes, &core.ModifyForeignKey{
					Table:       tableName,
					OldFK:       prevRel,
					NewFK:       rel,
					TargetTable: targetTable,
				})
			}
		}
	}

	// Detect dropped foreign keys
	for name := range previousMap {
		if _, exists := currentMap[name]; !exists {
			changes = append(changes, &core.DropForeignKey{
				Table:  tableName,
				FKName: fmt.Sprintf("fk_%s_%s", tableName, name),
			})
		}
	}

	return changes, nil
}

// relationChanged checks if a relation has changed
func (d *Detector) relationChanged(current, previous generator.RelationDefinition) bool {
	if current.To != previous.To {
		return true
	}
	currentOnDelete, _ := current.Options["on_delete"].(string)
	previousOnDelete, _ := previous.Options["on_delete"].(string)
	if currentOnDelete != previousOnDelete {
		return true
	}
	currentOnUpdate, _ := current.Options["on_update"].(string)
	previousOnUpdate, _ := previous.Options["on_update"].(string)
	if currentOnUpdate != previousOnUpdate {
		return true
	}
	return false
}

