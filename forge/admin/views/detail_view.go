package views

import (
	"context"
	"fmt"
	"net/http"

	"github.com/forgego/forge/admin"
)

// DetailView represents a type-safe detail view for viewing a single object
type DetailView[T any] struct {
	*BaseView[T]
}

// NewDetailView creates a new detail view
func NewDetailView[T any](admin *admin.Admin[T]) *DetailView[T] {
	return &DetailView[T]{
		BaseView: NewBaseView(admin),
	}
}

// DetailData contains data for rendering the detail view
type DetailData[T any] struct {
	Instance            *T
	Fields              []FieldDisplayData
	HasChangePermission bool
	HasDeletePermission bool
	ViewOnSiteURL       string
}

// FieldDisplayData contains data for displaying a field
type FieldDisplayData struct {
	Name        string
	Label       string
	Value       interface{}
	Formatted   string
	HelpText    string
	IsReadOnly  bool
}

// Render renders the detail view and returns the data
func (dv *DetailView[T]) Render(ctx context.Context, r *http.Request, user interface{}, instance *T) (*DetailData[T], error) {
	// Get all fields from schema
	fields := dv.admin.Fields()

	// Build field display data
	fieldData := make([]FieldDisplayData, 0, len(fields))
	for _, field := range fields {
		// Get field value from instance (would use ORM field accessor)
		var value interface{} = nil
		formatted := fmt.Sprintf("%v", value)

		fieldData = append(fieldData, FieldDisplayData{
			Name:       field.Name,
			Label:      field.VerboseName,
			Value:      value,
			Formatted:  formatted,
			HelpText:   field.HelpText,
			IsReadOnly: field.ReadOnly,
		})
	}

	// Get view on site URL if configured
	viewOnSiteURL := ""
	config := dv.admin.Config()
	if config != nil && config.ViewOnSite != nil {
		viewOnSiteURL = config.ViewOnSite(ctx, instance)
	}

	return &DetailData[T]{
		Instance:            instance,
		Fields:              fieldData,
		HasChangePermission: dv.admin.HasChangePermission(ctx, user, instance),
		HasDeletePermission: dv.admin.HasDeletePermission(ctx, user, instance),
		ViewOnSiteURL:       viewOnSiteURL,
	}, nil
}
