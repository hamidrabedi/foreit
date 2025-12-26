package detection

import (
	"reflect"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// detectColumnChanges detects column-level changes
func (d *Detector) detectColumnChanges(tableName string, current, previous []generator.FieldDefinition) ([]core.Change, error) {
	var changes []core.Change

	currentMap := make(map[string]generator.FieldDefinition)
	previousMap := make(map[string]generator.FieldDefinition)

	for _, field := range current {
		currentMap[field.Name] = field
	}

	for _, field := range previous {
		previousMap[field.Name] = field
	}

	// Detect new columns
	for name, field := range currentMap {
		if _, exists := previousMap[name]; !exists {
			changes = append(changes, &core.AddColumn{
				Table:  tableName,
				Column: field,
			})
		}
	}

	// Detect dropped columns
	for name := range previousMap {
		if _, exists := currentMap[name]; !exists {
			changes = append(changes, &core.DropColumn{
				Table:      tableName,
				ColumnName: name,
			})
		}
	}

	// Detect modified columns
	for name, currentField := range currentMap {
		if previousField, exists := previousMap[name]; exists {
			if d.fieldChanged(currentField, previousField) {
				changes = append(changes, &core.ModifyColumn{
					Table:     tableName,
					OldColumn: previousField,
					NewColumn: currentField,
				})
			}
		}
	}

	return changes, nil
}

// fieldChanged checks if a field has changed
func (d *Detector) fieldChanged(current, previous generator.FieldDefinition) bool {
	if current.Type != previous.Type {
		return true
	}
	if current.Required != previous.Required {
		return true
	}
	if current.PrimaryKey != previous.PrimaryKey {
		return true
	}
	if current.AutoIncrement != previous.AutoIncrement {
		return true
	}
	// Use reflect.DeepEqual for default value comparison
	if !reflect.DeepEqual(current.Default, previous.Default) {
		return true
	}
	// Use reflect.DeepEqual for options map comparison
	if !reflect.DeepEqual(current.Options, previous.Options) {
		return true
	}
	return false
}

