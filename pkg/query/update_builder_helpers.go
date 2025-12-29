package query

import "context"

// NewUpdateBuilderFromQuerySet creates an UpdateBuilder from a QuerySetV2
func NewUpdateBuilderFromQuerySet[T any](qs QuerySetV2[T]) (*UpdateBuilder[T], error) {
	schema, err := GetModelSchema[T]()
	if err != nil {
		return nil, err
	}
	
	return &UpdateBuilder[T]{
		qs:      &querySetWrapper[T]{qs: qs},
		schema:  schema,
		updates: make(map[string]interface{}),
	}, nil
}

// querySetWrapper wraps QuerySetV2 to work with UpdateBuilder
type querySetWrapper[T any] struct {
	qs QuerySetV2[T]
}

func (w *querySetWrapper[T]) Update(ctx context.Context, updates UpdateMap) (int64, error) {
	return w.qs.Update(ctx, updates)
}
