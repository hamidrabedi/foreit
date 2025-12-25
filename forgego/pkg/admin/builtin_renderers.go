package admin

import (
	"fmt"
	"html"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type StringFieldRenderer struct{}

func (r *StringFieldRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "string",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *StringFieldRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	value := ""
	if ctx.Value != nil {
		value = html.EscapeString(fmt.Sprintf("%v", ctx.Value))
	}
	html := fmt.Sprintf(`<input type="text" name="%s" value="%s" %s />`,
		ctx.Field.Name, value, attrs)
	return template.HTML(html), nil
}

func (r *StringFieldRenderer) Validate(value interface{}) error {
	if str, ok := value.(string); ok && len(str) > 1000 {
		return fmt.Errorf("string too long (max 1000 characters)")
	}
	return nil
}

type IntFieldRenderer struct{}

func (r *IntFieldRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "int",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *IntFieldRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	value := ""
	if ctx.Value != nil {
		value = fmt.Sprintf("%v", ctx.Value)
	}
	html := fmt.Sprintf(`<input type="number" name="%s" value="%s" %s />`,
		ctx.Field.Name, value, attrs)
	return template.HTML(html), nil
}

func (r *IntFieldRenderer) Validate(value interface{}) error {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return nil
	case string:
		if _, err := strconv.Atoi(v); err != nil {
			return fmt.Errorf("invalid integer: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid type for integer field: %T", value)
	}
}

type FloatFieldRenderer struct{}

func (r *FloatFieldRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "float",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *FloatFieldRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	value := ""
	if ctx.Value != nil {
		value = fmt.Sprintf("%v", ctx.Value)
	}
	html := fmt.Sprintf(`<input type="number" step="0.01" name="%s" value="%s" %s />`,
		ctx.Field.Name, value, attrs)
	return template.HTML(html), nil
}

func (r *FloatFieldRenderer) Validate(value interface{}) error {
	switch v := value.(type) {
	case float32, float64:
		return nil
	case string:
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return fmt.Errorf("invalid float: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid type for float field: %T", value)
	}
}

type BoolFieldRenderer struct{}

func (r *BoolFieldRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "bool",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *BoolFieldRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	checked := ""
	if ctx.Value != nil {
		if b, ok := ctx.Value.(bool); ok && b {
			checked = "checked"
		} else if s, ok := ctx.Value.(string); ok && (s == "true" || s == "1" || s == "on") {
			checked = "checked"
		}
	}
	html := fmt.Sprintf(`<input type="checkbox" name="%s" %s %s />`,
		ctx.Field.Name, checked, attrs)
	return template.HTML(html), nil
}

func (r *BoolFieldRenderer) Validate(value interface{}) error {
	switch value.(type) {
	case bool:
		return nil
	case string:
		return nil
	default:
		return fmt.Errorf("invalid type for boolean field: %T", value)
	}
}

type DateTimeFieldRenderer struct{}

func (r *DateTimeFieldRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "datetime",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *DateTimeFieldRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	value := ""
	if ctx.Value != nil {
		switch v := ctx.Value.(type) {
		case time.Time:
			value = v.Format("2006-01-02T15:04:05")
		case string:
			value = v
		default:
			value = fmt.Sprintf("%v", v)
		}
	}
	html := fmt.Sprintf(`<input type="datetime-local" name="%s" value="%s" %s />`,
		ctx.Field.Name, html.EscapeString(value), attrs)
	return template.HTML(html), nil
}

func (r *DateTimeFieldRenderer) Validate(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		return nil
	case string:
		if _, err := time.Parse("2006-01-02T15:04:05", v); err != nil {
			if _, err2 := time.Parse(time.RFC3339, v); err2 != nil {
				return fmt.Errorf("invalid datetime format: %w", err)
			}
		}
		return nil
	default:
		return fmt.Errorf("invalid type for datetime field: %T", value)
	}
}

type ForeignKeyRenderer struct{}

func (r *ForeignKeyRenderer) Info() RendererInfo {
	return RendererInfo{
		Name:    "foreignkey",
		Version: "1.0",
		Author:  "gogo",
	}
}

func (r *ForeignKeyRenderer) RenderHTML(ctx RenderContext) (template.HTML, error) {
	attrs := buildHTMXAttrs(ctx.HTMXAttrs)
	value := ""
	if ctx.Value != nil {
		value = html.EscapeString(fmt.Sprintf("%v", ctx.Value))
	}

	if len(ctx.Field.Choices) > 0 {
		var options []string
		for _, choice := range ctx.Field.Choices {
			selected := ""
			if choice.Value == value {
				selected = "selected"
			}
			options = append(options, fmt.Sprintf(`<option value="%s" %s>%s</option>`,
				html.EscapeString(choice.Value), selected, html.EscapeString(choice.Label)))
		}
		html := fmt.Sprintf(`<select name="%s" %s>%s</select>`,
			ctx.Field.Name, attrs, strings.Join(options, ""))
		return template.HTML(html), nil
	}

	html := fmt.Sprintf(`<input type="text" name="%s" value="%s" %s />`,
		ctx.Field.Name, value, attrs)
	return template.HTML(html), nil
}

func (r *ForeignKeyRenderer) Validate(value interface{}) error {
	if value == nil || value == "" {
		return nil
	}
	return nil
}

func buildHTMXAttrs(attrs map[string]string) string {
	if len(attrs) == 0 {
		return ""
	}
	var parts []string
	for k, v := range attrs {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, k, html.EscapeString(v)))
	}
	return strings.Join(parts, " ")
}

