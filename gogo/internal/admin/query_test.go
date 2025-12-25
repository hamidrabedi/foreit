package admin

import (
	"testing"
)

// TestQueryBuilder_ParseQueryParams tests query parameter parsing
func TestQueryBuilder_ParseQueryParams(t *testing.T) {
	meta := &ModelMeta{
		Name: "TestModel",
		Fields: []FieldMeta{
			{Name: "name", Filterable: true, Sortable: true},
			{Name: "age", Filterable: true, Sortable: true},
			{Name: "email", Filterable: true, Sortable: false},
		},
	}
	
	qb := NewQueryBuilder(meta)
	
	tests := []struct {
		name     string
		params   map[string]string
		expected *QueryParams
	}{
		{
			name: "basic pagination",
			params: map[string]string{
				"page":      "2",
				"page_size": "50",
			},
			expected: &QueryParams{
				Page:     2,
				PageSize: 50,
				SortBy:   "id",
				SortOrder: "asc",
				Filters:  make(map[string]interface{}),
			},
		},
		{
			name: "sorting",
			params: map[string]string{
				"sort_by":    "name",
				"sort_order": "desc",
			},
			expected: &QueryParams{
				Page:     1,
				PageSize: 20,
				SortBy:   "name",
				SortOrder: "desc",
				Filters:  make(map[string]interface{}),
			},
		},
		{
			name: "filters",
			params: map[string]string{
				"filter_name": "John",
				"filter_age__gt": "18",
			},
			expected: &QueryParams{
				Page:     1,
				PageSize: 20,
				SortBy:   "id",
				SortOrder: "asc",
			},
		},
		{
			name: "search",
			params: map[string]string{
				"search": "test",
				"search_fields": "name,email",
			},
			expected: &QueryParams{
				Page:     1,
				PageSize: 20,
				SortBy:   "id",
				SortOrder: "asc",
				Filters:  make(map[string]interface{}),
				Search:   "test",
				SearchFields: []string{"name", "email"},
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := qb.ParseQueryParams(tt.params)
			
			if result.Page != tt.expected.Page {
				t.Errorf("Expected page %d, got %d", tt.expected.Page, result.Page)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("Expected page_size %d, got %d", tt.expected.PageSize, result.PageSize)
			}
			if result.SortBy != tt.expected.SortBy {
				t.Errorf("Expected sort_by %s, got %s", tt.expected.SortBy, result.SortBy)
			}
			if result.SortOrder != tt.expected.SortOrder {
				t.Errorf("Expected sort_order %s, got %s", tt.expected.SortOrder, result.SortOrder)
			}
			if result.Search != tt.expected.Search {
				t.Errorf("Expected search %s, got %s", tt.expected.Search, result.Search)
			}
		})
	}
}

// TestQueryBuilder_ConvertValue tests value conversion
func TestQueryBuilder_ConvertValue(t *testing.T) {
	meta := &ModelMeta{
		Name: "TestModel",
		Fields: []FieldMeta{
			{Name: "age", Type: FieldTypeNumber},
			{Name: "active", Type: FieldTypeBoolean},
			{Name: "name", Type: FieldTypeText},
		},
	}
	
	qb := NewQueryBuilder(meta)
	
	tests := []struct {
		name      string
		fieldName string
		value     string
		expected  interface{}
		shouldErr bool
	}{
		{
			name:      "convert number",
			fieldName: "age",
			value:     "25",
			expected:  25,
			shouldErr: false,
		},
		{
			name:      "convert boolean true",
			fieldName: "active",
			value:     "true",
			expected:  true,
			shouldErr: false,
		},
		{
			name:      "convert boolean false",
			fieldName: "active",
			value:     "false",
			expected:  false,
			shouldErr: false,
		},
		{
			name:      "convert string",
			fieldName: "name",
			value:     "John",
			expected:  "John",
			shouldErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := qb.ConvertValue(tt.fieldName, tt.value)
			
			if tt.shouldErr && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if !tt.shouldErr && result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestCalculatePagination tests pagination calculation
func TestCalculatePagination(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		total      int64
		expected   PaginationResult
	}{
		{
			name:     "first page",
			page:     1,
			pageSize: 20,
			total:    100,
			expected: PaginationResult{
				Page:       1,
				PageSize:   20,
				Total:      100,
				TotalPages: 5,
			},
		},
		{
			name:     "last page with remainder",
			page:     3,
			pageSize: 20,
			total:    55,
			expected: PaginationResult{
				Page:       3,
				PageSize:   20,
				Total:      55,
				TotalPages: 3,
			},
		},
		{
			name:     "empty result",
			page:     1,
			pageSize: 20,
			total:    0,
			expected: PaginationResult{
				Page:       1,
				PageSize:   20,
				Total:      0,
				TotalPages: 0,
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePagination(tt.page, tt.pageSize, tt.total)
			
			if result.Page != tt.expected.Page {
				t.Errorf("Expected page %d, got %d", tt.expected.Page, result.Page)
			}
			if result.PageSize != tt.expected.PageSize {
				t.Errorf("Expected page_size %d, got %d", tt.expected.PageSize, result.PageSize)
			}
			if result.Total != tt.expected.Total {
				t.Errorf("Expected total %d, got %d", tt.expected.Total, result.Total)
			}
			if result.TotalPages != tt.expected.TotalPages {
				t.Errorf("Expected total_pages %d, got %d", tt.expected.TotalPages, result.TotalPages)
			}
		})
	}
}

