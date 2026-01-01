package generator

import (
	_ "embed"
)

//go:embed templates/model.tmpl
var modelTemplate string

//go:embed templates/fields.tmpl
var fieldsTemplate string

//go:embed templates/manager.tmpl
var managerTemplate string

//go:embed templates/queryset.tmpl
var querysetTemplate string

//go:embed templates/relations.tmpl
var relationTemplate string
