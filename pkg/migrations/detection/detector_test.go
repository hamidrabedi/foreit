package detection

import (
	"testing"

	"github.com/forgego/forge/pkg/generator"
	"github.com/forgego/forge/pkg/migrations/core"
)

func TestDetector_DetectTableChanges(t *testing.T) {
	detector := NewDetector()

	tests := []struct {
		name     string
		current  []*generator.ModelDefinition
		previous []*generator.ModelDefinition
		want     []core.ChangeType
	}{
		{
			name: "new table",
			current: []*generator.ModelDefinition{
				{
					Name: "User",
					Meta: generator.MetaDefinition{TableName: "users"},
					Fields: []generator.FieldDefinition{
						{Name: "id", Type: "Int64", PrimaryKey: true},
					},
				},
			},
			previous: []*generator.ModelDefinition{},
			want:     []core.ChangeType{core.ChangeTypeCreateTable},
		},
		{
			name: "dropped table",
			current: []*generator.ModelDefinition{},
			previous: []*generator.ModelDefinition{
				{
					Name: "User",
					Meta: generator.MetaDefinition{TableName: "users"},
				},
			},
			want: []core.ChangeType{core.ChangeTypeDropTable},
		},
		{
			name: "no changes",
			current: []*generator.ModelDefinition{
				{
					Name: "User",
					Meta: generator.MetaDefinition{TableName: "users"},
					Fields: []generator.FieldDefinition{
						{Name: "id", Type: "Int64", PrimaryKey: true},
					},
				},
			},
			previous: []*generator.ModelDefinition{
				{
					Name: "User",
					Meta: generator.MetaDefinition{TableName: "users"},
					Fields: []generator.FieldDefinition{
						{Name: "id", Type: "Int64", PrimaryKey: true},
					},
				},
			},
			want: []core.ChangeType{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes, err := detector.DetectChanges(tt.current, tt.previous)
			if err != nil {
				t.Fatalf("DetectChanges() error = %v", err)
			}

			if len(changes) != len(tt.want) {
				t.Errorf("Expected %d changes, got %d", len(tt.want), len(changes))
			}

			changeTypes := make(map[core.ChangeType]bool)
			for _, change := range changes {
				changeTypes[change.Type()] = true
			}

			for _, wantType := range tt.want {
				if !changeTypes[wantType] {
					t.Errorf("Expected change type %s not found", wantType)
				}
			}
		})
	}
}

func TestDetector_DetectColumnChanges(t *testing.T) {
	detector := NewDetector()

	current := []*generator.ModelDefinition{
		{
			Name: "User",
			Meta: generator.MetaDefinition{TableName: "users"},
			Fields: []generator.FieldDefinition{
				{Name: "id", Type: "Int64", PrimaryKey: true},
				{Name: "email", Type: "String", Required: true},
				{Name: "name", Type: "String"}, // New field
			},
		},
	}

	previous := []*generator.ModelDefinition{
		{
			Name: "User",
			Meta: generator.MetaDefinition{TableName: "users"},
			Fields: []generator.FieldDefinition{
				{Name: "id", Type: "Int64", PrimaryKey: true},
				{Name: "email", Type: "String", Required: true},
			},
		},
	}

	changes, err := detector.DetectChanges(current, previous)
	if err != nil {
		t.Fatalf("DetectChanges() error = %v", err)
	}

	foundAddColumn := false
	for _, change := range changes {
		if change.Type() == core.ChangeTypeAddColumn {
			if addCol, ok := change.(*core.AddColumn); ok && addCol.Column.Name == "name" {
				foundAddColumn = true
				break
			}
		}
	}

	if !foundAddColumn {
		t.Error("Expected AddColumn change for 'name' field")
	}
}

