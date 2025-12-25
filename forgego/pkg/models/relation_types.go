package models

// RelationType represents the type of a model relationship.
type RelationType string

const (
	RelationTypeForeignKey RelationType = "foreignkey"
	RelationTypeOneToOne   RelationType = "onetoone"
	RelationTypeOneToMany  RelationType = "onetomany"
	RelationTypeManyToMany RelationType = "manytomany"
)

