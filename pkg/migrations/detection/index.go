package detection

import (
	"fmt"
	"strings"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// detectIndexChanges detects index changes
func (d *Detector) detectIndexChanges(tableName string, current, previous []generator.IndexDefinition) ([]core.Change, error) {
	var changes []core.Change

	currentMap := make(map[string]generator.IndexDefinition)
	previousMap := make(map[string]generator.IndexDefinition)

	for _, idx := range current {
		idxName := idx.Name
		if idxName == "" {
			idxName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
		}
		currentMap[idxName] = idx
	}

	for _, idx := range previous {
		idxName := idx.Name
		if idxName == "" {
			idxName = fmt.Sprintf("idx_%s_%s", tableName, strings.Join(idx.Fields, "_"))
		}
		previousMap[idxName] = idx
	}

	// Detect new indexes
	for name, idx := range currentMap {
		if _, exists := previousMap[name]; !exists {
			changes = append(changes, &core.AddIndex{
				Table: tableName,
				Index: idx,
			})
		} else {
			// Check if index changed
			prevIdx := previousMap[name]
			if d.indexChanged(idx, prevIdx) {
				changes = append(changes, &core.ModifyIndex{
					Table:   tableName,
					OldIndex: prevIdx,
					NewIndex: idx,
				})
			}
		}
	}

	// Detect dropped indexes
	for name := range previousMap {
		if _, exists := currentMap[name]; !exists {
			changes = append(changes, &core.DropIndex{
				Table:     tableName,
				IndexName: name,
			})
		}
	}

	return changes, nil
}

// indexChanged checks if an index has changed
func (d *Detector) indexChanged(current, previous generator.IndexDefinition) bool {
	if current.Unique != previous.Unique {
		return true
	}
	if len(current.Fields) != len(previous.Fields) {
		return true
	}
	for i, field := range current.Fields {
		if field != previous.Fields[i] {
			return true
		}
	}
	return false
}

