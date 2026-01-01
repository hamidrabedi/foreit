package orm

import (
	"fmt"
)

// IsPreloaded checks if a relation is preloaded (for N+1 prevention)
func (qs *BaseQuerySet[T]) IsPreloaded(relationName string) bool {
	return qs.preloaded[relationName]
}

// CheckPreloaded validates that a relation is preloaded
// Returns error if relation is not preloaded (N+1 prevention)
func (qs *BaseQuerySet[T]) CheckPreloaded(relationName string) error {
	if !qs.IsPreloaded(relationName) {
		return fmt.Errorf("relation '%s' accessed but not preloaded - use Preload() or PrefetchRelated() to prevent N+1 queries", relationName)
	}
	return nil
}

// PreloadDynamic is a dynamic API for preloading relations
// Usage: qs.PreloadDynamic("author", "comments")
func (qs *BaseQuerySet[T]) PreloadDynamic(relations ...string) QuerySet[T] {
	clone := qs.clone()
	for _, relation := range relations {
		clone.prefetchRelated = append(clone.prefetchRelated, relation)
		clone.preloaded[relation] = true
	}
	return clone
}

// ErrRelationNotLoaded is returned when accessing a non-preloaded relation
type ErrRelationNotLoaded struct {
	Relation string
}

func (e ErrRelationNotLoaded) Error() string {
	return fmt.Sprintf("relation '%s' accessed but not preloaded - use Preload() or PrefetchRelated() to prevent N+1 queries", e.Relation)
}
