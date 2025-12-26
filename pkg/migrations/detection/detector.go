package detection

import (
	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

// ChangeDetector detects changes between current and previous model states
type ChangeDetector interface {
	// DetectChanges compares current models to previous state and returns all changes
	DetectChanges(current, previous []*generator.ModelDefinition) ([]core.Change, error)
}

// Detector is the default implementation of ChangeDetector
type Detector struct{}

// NewDetector creates a new change detector
func NewDetector() ChangeDetector {
	return &Detector{}
}

// DetectChanges compares current models to previous state and returns all changes
func (d *Detector) DetectChanges(current, previous []*generator.ModelDefinition) ([]core.Change, error) {
	var changes []core.Change

	// Build maps for easier lookup
	currentMap := make(map[string]*generator.ModelDefinition)
	previousMap := make(map[string]*generator.ModelDefinition)

	for _, def := range current {
		tableName := getTableName(def)
		currentMap[tableName] = def
	}

	for _, def := range previous {
		tableName := getTableName(def)
		previousMap[tableName] = def
	}

	// Detect table-level changes
	for tableName, currentDef := range currentMap {
		previousDef, exists := previousMap[tableName]
		if !exists {
			// New table - create table and add all indexes/foreign keys
			changes = append(changes, &core.CreateTable{Table: currentDef})
			
			// Add indexes for new table
			for _, idx := range currentDef.Meta.Indexes {
				changes = append(changes, &core.AddIndex{
					Table: tableName,
					Index: idx,
				})
			}
			
			// Add foreign keys for new table
			for _, rel := range currentDef.Relations {
				if rel.Type == "ForeignKey" || rel.Type == "OneToOne" {
					// Check if the relation column exists in fields
					hasColumn := false
					for _, field := range currentDef.Fields {
						if field.Name == rel.Name {
							hasColumn = true
							break
						}
					}
					if !hasColumn {
						continue
					}

					targetTable := findTargetTable(rel.To, current)
					if targetTable != "" {
						changes = append(changes, &core.AddForeignKey{
							Table:       tableName,
							Relation:    rel,
							TargetTable: targetTable,
						})
					}
				}
			}
		} else if previousDef != nil {
			// Existing table - detect column, index, constraint changes
			tableChanges, err := d.detectTableChanges(currentDef, previousDef, tableName, current)
			if err != nil {
				return nil, err
			}
			changes = append(changes, tableChanges...)
		}
	}

	// Detect dropped tables
	for tableName := range previousMap {
		if _, exists := currentMap[tableName]; !exists {
			changes = append(changes, &core.DropTable{Table: tableName})
		}
	}

	return changes, nil
}

// detectTableChanges detects changes within a single table
func (d *Detector) detectTableChanges(current, previous *generator.ModelDefinition, tableName string, allDefs []*generator.ModelDefinition) ([]core.Change, error) {
	var changes []core.Change

	// Detect column changes
	columnChanges, err := d.detectColumnChanges(tableName, current.Fields, previous.Fields)
	if err != nil {
		return nil, err
	}
	changes = append(changes, columnChanges...)

	// Detect index changes
	indexChanges, err := d.detectIndexChanges(tableName, current.Meta.Indexes, previous.Meta.Indexes)
	if err != nil {
		return nil, err
	}
	changes = append(changes, indexChanges...)

	// Detect constraint changes
	constraintChanges, err := d.detectConstraintChanges(tableName, current.Meta.Constraints, previous.Meta.Constraints)
	if err != nil {
		return nil, err
	}
	changes = append(changes, constraintChanges...)

	// Detect foreign key changes
	fkChanges, err := d.detectForeignKeyChanges(tableName, current.Relations, previous.Relations, current, previous, allDefs)
	if err != nil {
		return nil, err
	}
	changes = append(changes, fkChanges...)

	return changes, nil
}


