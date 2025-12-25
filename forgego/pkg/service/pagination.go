package service

import (
	"context"

	"github.com/forgego/forge/pkg/models"
)

type DefaultPagination[T any] struct{}

func (p *DefaultPagination[T]) Paginate(ctx context.Context, qs models.QuerySet[T], params ListParams) (*PaginationResult[T], error) {
	total, err := qs.Count(ctx)
	if err != nil {
		return nil, err
	}

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 200 {
		pageSize = 200
	}

	offset := (params.Page - 1) * pageSize
	results, err := qs.Limit(pageSize).Offset(offset).All(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]T, len(results))
	for i, r := range results {
		items[i] = *r
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	hasMore := params.Page < totalPages

	return &PaginationResult[T]{
		Items:   items,
		Total:   int64(total),
		HasMore: hasMore,
	}, nil
}

