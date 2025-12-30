package components

import (
	"fmt"
	"html/template"
	"math"
)

// Pagination represents pagination component
type Pagination struct {
	CurrentPage  int
	TotalPages   int
	TotalCount   int64
	PageSize     int
	BaseURL      string
	ShowFirstLast bool
	ShowPrevNext  bool
	MaxVisible    int // Maximum number of page links to show
}

// NewPagination creates a new pagination component
func NewPagination(currentPage, totalPages, pageSize int, totalCount int64, baseURL string) *Pagination {
	return &Pagination{
		CurrentPage:    currentPage,
		TotalPages:     totalPages,
		TotalCount:     totalCount,
		PageSize:       pageSize,
		BaseURL:        baseURL,
		ShowFirstLast:  true,
		ShowPrevNext:   true,
		MaxVisible:     10,
	}
}

// Render renders the pagination as HTML
func (p *Pagination) Render() template.HTML {
	if p.TotalPages <= 1 {
		return template.HTML("")
	}

	html := `<nav class="pagination" aria-label="Page navigation">`
	html += `<ul class="pagination-list">`

	// First page
	if p.ShowFirstLast && p.CurrentPage > 1 {
		html += `<li class="page-item"><a class="page-link" href="` + p.pageURL(1) + `">First</a></li>`
	}

	// Previous page
	if p.ShowPrevNext && p.CurrentPage > 1 {
		html += `<li class="page-item"><a class="page-link" href="` + p.pageURL(p.CurrentPage-1) + `">Previous</a></li>`
	}

	// Page numbers
	startPage, endPage := p.getVisiblePageRange()
	for i := startPage; i <= endPage; i++ {
		if i == p.CurrentPage {
			html += `<li class="page-item active"><span class="page-link">` + fmt.Sprintf("%d", i) + `</span></li>`
		} else {
			html += `<li class="page-item"><a class="page-link" href="` + p.pageURL(i) + `">` + fmt.Sprintf("%d", i) + `</a></li>`
		}
	}

	// Next page
	if p.ShowPrevNext && p.CurrentPage < p.TotalPages {
		html += `<li class="page-item"><a class="page-link" href="` + p.pageURL(p.CurrentPage+1) + `">Next</a></li>`
	}

	// Last page
	if p.ShowFirstLast && p.CurrentPage < p.TotalPages {
		html += `<li class="page-item"><a class="page-link" href="` + p.pageURL(p.TotalPages) + `">Last</a></li>`
	}

	html += `</ul>`
	html += `<div class="pagination-info">`
	html += fmt.Sprintf(`Showing %d to %d of %d results`,
		p.getStartIndex(), p.getEndIndex(), p.TotalCount)
	html += `</div>`
	html += `</nav>`

	return template.HTML(html)
}

// getVisiblePageRange calculates the range of page numbers to display
func (p *Pagination) getVisiblePageRange() (start, end int) {
	halfVisible := p.MaxVisible / 2
	start = p.CurrentPage - halfVisible
	end = p.CurrentPage + halfVisible

	if start < 1 {
		start = 1
		end = int(math.Min(float64(p.MaxVisible), float64(p.TotalPages)))
	}

	if end > p.TotalPages {
		end = p.TotalPages
		start = int(math.Max(1, float64(end-p.MaxVisible+1)))
	}

	return start, end
}

// pageURL generates URL for a specific page
func (p *Pagination) pageURL(page int) string {
	// Simple implementation - would need to preserve query params
	return fmt.Sprintf("%s?page=%d", p.BaseURL, page)
}

// getStartIndex gets the starting index of current page
func (p *Pagination) getStartIndex() int {
	return (p.CurrentPage-1)*p.PageSize + 1
}

// getEndIndex gets the ending index of current page
func (p *Pagination) getEndIndex() int {
	end := p.CurrentPage * p.PageSize
	if int64(end) > p.TotalCount {
		return int(p.TotalCount)
	}
	return end
}
