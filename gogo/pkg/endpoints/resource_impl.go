package endpoints

import (
	"context"
	"strconv"
)

func (r *BaseResource[T, Q]) Index(ctx *Context) ([]*T, error) {
	query := r.Repo.Query()
	query = ApplyQuery(query, ctx.Ctx)
	
	results, err := r.Repo.All(ctx.Request.Context(), query)
	if err != nil {
		return nil, err
	}
	
	return results, nil
}

func (r *BaseResource[T, Q]) Show(ctx *Context) (*T, error) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, NewError(400, "invalid_id", "Invalid ID format")
	}
	
	result, err := r.Repo.GetByID(ctx.Request.Context(), idInt)
	if err != nil {
		return nil, &Error{Code: "not_found", Message: "Resource not found", Status: 404}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Create(ctx *Context) (*T, error) {
	var data T
	if err := ctx.Bind(&data); err != nil {
		return nil, NewError(400, "invalid_data", "Invalid request data")
	}
	
	result, err := r.Repo.Create(ctx.Request.Context(), &data)
	if err != nil {
		return nil, &Error{Code: "create_failed", Message: err.Error(), Status: 500}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Update(ctx *Context) (*T, error) {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, NewError(400, "invalid_id", "Invalid ID format")
	}
	
	var data T
	if err := ctx.Bind(&data); err != nil {
		return nil, NewError(400, "invalid_data", "Invalid request data")
	}
	
	result, err := r.Repo.Update(ctx.Request.Context(), idInt, &data)
	if err != nil {
		return nil, &Error{Code: "update_failed", Message: err.Error(), Status: 500}
	}
	
	return result, nil
}

func (r *BaseResource[T, Q]) Destroy(ctx *Context) error {
	id := ctx.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return &Error{Code: "invalid_id", Message: "Invalid ID format", Status: 400}
	}
	
	if err := r.Repo.Delete(ctx.Request.Context(), idInt); err != nil {
		return &Error{Code: "delete_failed", Message: err.Error(), Status: 500}
	}
	
	return nil
}

