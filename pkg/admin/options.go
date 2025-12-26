package admin

import (
	"context"
)

// WithModelAdminConfig applies a full ModelAdminConfig to a model
func WithModelAdminConfig(config ModelAdminConfig) AdminOption {
	return func(m *AdminModel) {
		ApplyModelAdminConfig(m, config)
	}
}

// WithDateHierarchy sets the date hierarchy field
func WithDateHierarchy(field string) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["date_hierarchy"] = field
	}
}

// WithOrdering sets the default ordering
func WithOrdering(fields ...string) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["ordering"] = fields
	}
}

// WithListPerPage sets items per page
func WithListPerPage(count int) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["list_per_page"] = count
	}
}

// WithFieldsets sets fieldsets for form grouping
func WithFieldsets(fieldsets ...Fieldset) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["fieldsets"] = fieldsets
	}
}

// WithVerboseName sets the verbose name
func WithVerboseName(name string) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["verbose_name"] = name
	}
}

// WithVerboseNamePlural sets the plural verbose name
func WithVerboseNamePlural(name string) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["verbose_name_plural"] = name
	}
}

// WithSaveButtons configures save button options
func WithSaveButtons(saveOnTop, saveAs, saveAsContinue, saveAndAddAnother bool) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["save_on_top"] = saveOnTop
		m.ExtendedConfig["save_as"] = saveAs
		m.ExtendedConfig["save_as_continue"] = saveAsContinue
		m.ExtendedConfig["save_and_add_another"] = saveAndAddAnother
	}
}

// WithCustomQueryset sets a custom queryset function
func WithCustomQueryset(fn func(context.Context, interface{}) (interface{}, error)) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["get_queryset"] = fn
	}
}

// WithCustomSave sets a custom save hook
func WithCustomSave(fn func(context.Context, interface{}, map[string]interface{}, bool) error) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["save_model"] = fn
	}
}

// WithCustomDelete sets a custom delete hook
func WithCustomDelete(fn func(context.Context, interface{}) error) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["delete_model"] = fn
	}
}

// WithExportFormats sets supported export formats
func WithExportFormats(formats ...string) AdminOption {
	return func(m *AdminModel) {
		if m.ExtendedConfig == nil {
			m.ExtendedConfig = make(map[string]interface{})
		}
		m.ExtendedConfig["export_formats"] = formats
	}
}

