package rest

type PaginatedResponse[T any] struct {
	Data       []*T            `json:"data"`
	Pagination PaginationResult `json:"pagination"`
}

func NewPaginatedResponse[T any](data []*T, page, pageSize int, total int64) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{
		Data:       data,
		Pagination: CalculatePagination(page, pageSize, total),
	}
}

func (r *BaseResource[T, Q]) IndexPaginated(ctx *Context) (*PaginatedResponse[T], error) {
	query := r.Repo.Query()
	
	query = ApplyQuery(query, ctx.Ctx)
	
	page, pageSize := ParsePagination(ctx.Ctx)
	
	total, err := r.Repo.Count(ctx.Context(), query)
	if err != nil {
		return nil, err
	}
	
	if paginatedQuery, ok := any(query).(interface {
		Offset(int) interface{}
		Limit(int) interface{}
	}); ok {
		offset := (page - 1) * pageSize
		query = paginatedQuery.Offset(offset).(Q)
		query = paginatedQuery.Limit(pageSize).(Q)
	}
	
	results, err := r.Repo.All(ctx.Context(), query)
	if err != nil {
		return nil, err
	}
	
	return NewPaginatedResponse(results, page, pageSize, int64(total)), nil
}

