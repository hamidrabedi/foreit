package api

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    *MetaInfo   `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type MetaInfo struct {
	Pagination *PaginationMeta `json:"pagination,omitempty"`
	Version    string          `json:"version,omitempty"`
	Timestamp  string          `json:"timestamp,omitempty"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func Success(c *fiber.Ctx, data interface{}) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
	})
}

func SuccessWithMeta(c *fiber.Ctx, data interface{}, meta *MetaInfo) error {
	return c.JSON(Response{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}

func ErrorWithDetails(c *fiber.Ctx, status int, code, message string, details interface{}) error {
	return c.Status(status).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func Paginated(c *fiber.Ctx, data interface{}, page, pageSize int, total int64) error {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	
	return SuccessWithMeta(c, data, &MetaInfo{
		Pagination: &PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

