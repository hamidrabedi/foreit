package schema

import (
	"github.com/forgego/forge/schema"
)

// RelationMapper maps schema relations to admin relation configurations
type RelationMapper struct{}

// NewRelationMapper creates a new relation mapper
func NewRelationMapper() *RelationMapper {
	return &RelationMapper{}
}

// MapRelationType maps a schema relation type to admin relation type
func (rm *RelationMapper) MapRelationType(relType schema.RelationType) string {
	switch relType {
	case schema.RelationForeignKey:
		return "foreignkey"
	case schema.RelationOneToOne:
		return "onetoone"
	case schema.RelationManyToMany:
		return "manytomany"
	default:
		return "unknown"
	}
}

// ShouldDisplayInForm determines if a relation should be displayed in form by default
func (rm *RelationMapper) ShouldDisplayInForm(relation schema.Relation) bool {
	// All relations can be displayed, but optional ones might be hidden
	return true
}

// GetDefaultWidgetType returns the default widget type for a relation
func (rm *RelationMapper) GetDefaultWidgetType(relation schema.Relation) string {
	switch relation.Type {
	case schema.RelationForeignKey, schema.RelationOneToOne:
		return "select"
	case schema.RelationManyToMany:
		return "multiselect"
	default:
		return "select"
	}
}

// IsRequired determines if a relation is required
func (rm *RelationMapper) IsRequired(relation schema.Relation) bool {
	// If LimitChoicesTo is set, it might indicate required
	// For now, we'll check if it's explicitly marked as required
	// This would need to be enhanced based on schema relation options
	return relation.LimitChoicesTo != nil
}
