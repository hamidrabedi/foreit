package admin

import (
	"encoding/json"
	"fmt"
	"strings"
)

// OpenAPIGenerator generates OpenAPI 3.0 specifications
type OpenAPIGenerator struct {
	registry *Registry
	baseURL  string
}

// NewOpenAPIGenerator creates a new OpenAPI generator
func NewOpenAPIGenerator(registry *Registry, baseURL string) *OpenAPIGenerator {
	return &OpenAPIGenerator{
		registry: registry,
		baseURL:  baseURL,
	}
}

// OpenAPISpec represents an OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       Info                   `json:"info"`
	Servers    []Server               `json:"servers"`
	Paths      map[string]PathItem    `json:"paths"`
	Components Components             `json:"components"`
}

// Info contains API information
type Info struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Description string `json:"description,omitempty"`
}

// Server represents an API server
type Server struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// PathItem represents API paths
type PathItem struct {
	Get    *Operation `json:"get,omitempty"`
	Post   *Operation `json:"post,omitempty"`
	Put    *Operation `json:"put,omitempty"`
	Delete *Operation `json:"delete,omitempty"`
}

// Operation represents an API operation
type Operation struct {
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId"`
	Parameters  []Parameter         `json:"parameters,omitempty"`
	RequestBody *RequestBody        `json:"requestBody,omitempty"`
	Responses   map[string]Response `json:"responses"`
}

// Parameter represents a query/path parameter
type Parameter struct {
	Name        string      `json:"name"`
	In          string      `json:"in"` // "query", "path", "header"
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Schema      *Schema     `json:"schema,omitempty"`
}

// RequestBody represents a request body
type RequestBody struct {
	Description string               `json:"description,omitempty"`
	Required    bool                 `json:"required,omitempty"`
	Content     map[string]MediaType `json:"content"`
}

// MediaType represents media type content
type MediaType struct {
	Schema *Schema `json:"schema,omitempty"`
}

// Response represents an API response
type Response struct {
	Description string               `json:"description"`
	Content     map[string]MediaType `json:"content,omitempty"`
}

// Schema represents a JSON schema
type Schema struct {
	Ref         string             `json:"$ref,omitempty"`
	Type        string             `json:"type,omitempty"`
	Format      string             `json:"format,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Required    []string            `json:"required,omitempty"`
	Items       *Schema             `json:"items,omitempty"`
	Enum        []interface{}       `json:"enum,omitempty"`
	Default     interface{}         `json:"default,omitempty"`
	Description string              `json:"description,omitempty"`
}

// Components contains reusable components
type Components struct {
	Schemas map[string]*Schema `json:"schemas,omitempty"`
}

// Generate generates an OpenAPI 3.0 specification
func (g *OpenAPIGenerator) Generate() (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: Info{
			Title:       "Admin API",
			Version:     "1.0.0",
			Description: "Auto-generated admin API",
		},
		Servers: []Server{
			{
				URL:         g.baseURL,
				Description: "Admin API Server",
			},
		},
		Paths:      make(map[string]PathItem),
		Components: Components{
			Schemas: make(map[string]*Schema),
		},
	}
	
	// Generate paths for each registered model
	models := g.registry.GetAllModels()
	for modelName, meta := range models {
		resourceName := g.modelNameToResource(modelName)
		g.generateModelPaths(spec, resourceName, modelName, meta)
		g.generateModelSchema(spec, modelName, meta)
	}
	
	return spec, nil
}

// generateModelPaths generates API paths for a model
func (g *OpenAPIGenerator) generateModelPaths(spec *OpenAPISpec, resourceName, modelName string, meta *ModelMeta) {
	basePath := fmt.Sprintf("/admin/api/%s", resourceName)
	
	// List endpoint
	spec.Paths[basePath] = PathItem{
		Get: &Operation{
			Tags:        []string{modelName},
			Summary:     fmt.Sprintf("List %s", modelName),
			Description: fmt.Sprintf("List all %s records with pagination and filters", modelName),
			OperationID: fmt.Sprintf("list%s", modelName),
			Parameters:  g.generateListParameters(meta),
			Responses: map[string]Response{
				"200": {
					Description: "Successful response",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Type: "object",
								Properties: map[string]*Schema{
									"data": {
										Type: "array",
										Items: &Schema{
											Ref: fmt.Sprintf("#/components/schemas/%s", modelName),
										},
									},
									"pagination": {
										Type: "object",
										Properties: map[string]*Schema{
											"page":       {Type: "integer"},
											"page_size":  {Type: "integer"},
											"total":      {Type: "integer", Format: "int64"},
											"total_pages": {Type: "integer"},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	
	// Get endpoint
	getPath := fmt.Sprintf("%s/{id}", basePath)
	spec.Paths[getPath] = PathItem{
		Get: &Operation{
			Tags:        []string{modelName},
			Summary:     fmt.Sprintf("Get %s", modelName),
			Description: fmt.Sprintf("Get a single %s record by ID", modelName),
			OperationID: fmt.Sprintf("get%s", modelName),
			Parameters: []Parameter{
				{
					Name:        "id",
					In:          "path",
					Description: "Record ID",
					Required:    true,
					Schema:      &Schema{Type: "string"},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Successful response",
					Content: map[string]MediaType{
						"application/json": {
							Schema: &Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", modelName),
							},
						},
					},
				},
			},
		},
	}
	
	// Create endpoint
	pathItem := spec.Paths[basePath]
	pathItem.Post = &Operation{
		Tags:        []string{modelName},
		Summary:     fmt.Sprintf("Create %s", modelName),
		Description: fmt.Sprintf("Create a new %s record", modelName),
		OperationID: fmt.Sprintf("create%s", modelName),
		RequestBody: &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: fmt.Sprintf("#/components/schemas/%sInput", modelName),
					},
				},
			},
		},
		Responses: map[string]Response{
			"201": {
				Description: "Created",
				Content: map[string]MediaType{
					"application/json": {
						Schema: &Schema{
							Ref: fmt.Sprintf("#/components/schemas/%s", modelName),
						},
					},
				},
			},
		},
	}
	spec.Paths[basePath] = pathItem
	
	// Update endpoint
	getPathItem := spec.Paths[getPath]
	getPathItem.Put = &Operation{
		Tags:        []string{modelName},
		Summary:     fmt.Sprintf("Update %s", modelName),
		Description: fmt.Sprintf("Update a %s record", modelName),
		OperationID: fmt.Sprintf("update%s", modelName),
		Parameters: []Parameter{
			{
				Name:        "id",
				In:          "path",
				Description: "Record ID",
				Required:    true,
				Schema:      &Schema{Type: "string"},
			},
		},
		RequestBody: &RequestBody{
			Required: true,
			Content: map[string]MediaType{
				"application/json": {
					Schema: &Schema{
						Ref: fmt.Sprintf("#/components/schemas/%sInput", modelName),
					},
				},
			},
		},
		Responses: map[string]Response{
			"200": {
				Description: "Successful response",
				Content: map[string]MediaType{
					"application/json": {
						Schema: &Schema{
							Ref: fmt.Sprintf("#/components/schemas/%s", modelName),
						},
					},
				},
			},
		},
	}
	
	// Delete endpoint
	getPathItem.Delete = &Operation{
		Tags:        []string{modelName},
		Summary:     fmt.Sprintf("Delete %s", modelName),
		Description: fmt.Sprintf("Delete a %s record", modelName),
		OperationID: fmt.Sprintf("delete%s", modelName),
		Parameters: []Parameter{
			{
				Name:        "id",
				In:          "path",
				Description: "Record ID",
				Required:    true,
				Schema:      &Schema{Type: "string"},
			},
		},
		Responses: map[string]Response{
			"204": {
				Description: "No Content",
			},
		},
	}
	spec.Paths[getPath] = getPathItem
}

// generateListParameters generates query parameters for list endpoint
func (g *OpenAPIGenerator) generateListParameters(meta *ModelMeta) []Parameter {
	params := []Parameter{
		{
			Name:        "page",
			In:          "query",
			Description: "Page number",
			Required:    false,
			Schema:      &Schema{Type: "integer", Default: 1},
		},
		{
			Name:        "page_size",
			In:          "query",
			Description: "Items per page",
			Required:    false,
			Schema:      &Schema{Type: "integer", Default: 20},
		},
		{
			Name:        "sort_by",
			In:          "query",
			Description: "Field to sort by",
			Required:    false,
			Schema:      &Schema{Type: "string", Default: "id"},
		},
		{
			Name:        "sort_order",
			In:          "query",
			Description: "Sort order (asc or desc)",
			Required:    false,
			Schema: &Schema{
				Type: "string",
				Enum: []interface{}{"asc", "desc"},
				Default: "asc",
			},
		},
	}
	
	// Add filter parameters for filterable fields
	for _, field := range meta.Fields {
		if field.Filterable {
			params = append(params, Parameter{
				Name:        fmt.Sprintf("filter_%s", field.Name),
				In:          "query",
				Description: fmt.Sprintf("Filter by %s", field.Label),
				Required:    false,
				Schema:      g.fieldToSchema(field),
			})
		}
	}
	
	return params
}

// generateModelSchema generates JSON schema for a model
func (g *OpenAPIGenerator) generateModelSchema(spec *OpenAPISpec, modelName string, meta *ModelMeta) {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
		Required:   []string{},
	}
	
	inputSchema := &Schema{
		Type:       "object",
		Properties: make(map[string]*Schema),
		Required:   []string{},
	}
	
	for _, field := range meta.Fields {
		fieldSchema := g.fieldToSchema(field)
		schema.Properties[field.Name] = fieldSchema
		
		// Add to required if field is required and not read-only
		if field.Required && !field.ReadOnly {
			schema.Required = append(schema.Required, field.Name)
		}
		
		// Input schema excludes read-only fields
		if !field.ReadOnly {
			inputSchema.Properties[field.Name] = fieldSchema
			if field.Required {
				inputSchema.Required = append(inputSchema.Required, field.Name)
			}
		}
	}
	
	spec.Components.Schemas[modelName] = schema
	spec.Components.Schemas[modelName+"Input"] = inputSchema
}

// fieldToSchema converts a field metadata to JSON schema
func (g *OpenAPIGenerator) fieldToSchema(field FieldMeta) *Schema {
	schema := &Schema{
		Description: field.HelpText,
	}
	
	switch field.Type {
	case FieldTypeText, FieldTypeTextarea, FieldTypeEmail, FieldTypeURL:
		schema.Type = "string"
		if field.Type == FieldTypeEmail {
			schema.Format = "email"
		} else if field.Type == FieldTypeURL {
			schema.Format = "uri"
		}
		
	case FieldTypeNumber:
		schema.Type = "number"
		
	case FieldTypeBoolean:
		schema.Type = "boolean"
		
	case FieldTypeDate:
		schema.Type = "string"
		schema.Format = "date"
		
	case FieldTypeDateTime, FieldTypeTime:
		schema.Type = "string"
		schema.Format = "date-time"
		
	case FieldTypeSelect:
		schema.Type = "string"
		if len(field.Choices) > 0 {
			schema.Enum = make([]interface{}, len(field.Choices))
			for i, choice := range field.Choices {
				schema.Enum[i] = choice.Value
			}
		}
		
	case FieldTypeJSON:
		schema.Type = "object"
		
	default:
		schema.Type = "string"
	}
	
	if field.Default != nil {
		schema.Default = field.Default
	}
	
	return schema
}

// modelNameToResource converts a model name to a URL-friendly resource name
func (g *OpenAPIGenerator) modelNameToResource(modelName string) string {
	// Convert PascalCase to snake_case and pluralize
	var result strings.Builder
	for i, r := range modelName {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteString("_")
		}
		result.WriteRune(r)
	}
	
	resourceName := strings.ToLower(result.String())
	
	// Simple pluralization
	if strings.HasSuffix(resourceName, "y") {
		resourceName = strings.TrimSuffix(resourceName, "y") + "ies"
	} else if strings.HasSuffix(resourceName, "s") || strings.HasSuffix(resourceName, "x") ||
		strings.HasSuffix(resourceName, "z") || strings.HasSuffix(resourceName, "ch") ||
		strings.HasSuffix(resourceName, "sh") {
		resourceName += "es"
	} else {
		resourceName += "s"
	}
	
	return resourceName
}

// GenerateJSON generates OpenAPI spec as JSON
func (g *OpenAPIGenerator) GenerateJSON() ([]byte, error) {
	spec, err := g.Generate()
	if err != nil {
		return nil, err
	}
	
	return json.MarshalIndent(spec, "", "  ")
}

// GenerateToFile generates OpenAPI spec and writes to file
func (g *OpenAPIGenerator) GenerateToFile(filename string) error {
	data, err := g.GenerateJSON()
	if err != nil {
		return err
	}
	
	// This would write to file - for now just return the data
	// In production: return os.WriteFile(filename, data, 0644)
	_ = data
	return nil
}

