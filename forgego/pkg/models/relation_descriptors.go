package models

// RelationDescriptor is the interface for relation definitions.
type RelationDescriptor interface {
	GetName() string
	GetType() RelationType
}

// Note: RelationType is defined in model.go
// We use the existing RelationType constants

