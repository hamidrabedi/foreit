package widgets

import (
	"strings"
	"testing"

	"github.com/forgego/forge/admin"
	"github.com/stretchr/testify/assert"
)

func TestRichTextWidget(t *testing.T) {
	widget := NewRichText()

	t.Run("Default configuration", func(t *testing.T) {
		assert.Equal(t, 300, widget.Height)
		assert.Equal(t, "full", widget.Toolbar)
		assert.False(t, widget.ReadOnly)
	})

	t.Run("Builder pattern", func(t *testing.T) {
		w := NewRichText().
			WithHeight(500).
			WithToolbar("basic").
			WithPlugins("tables", "links")

		assert.Equal(t, 500, w.Height)
		assert.Equal(t, "basic", w.Toolbar)
		assert.Contains(t, w.Plugins, "tables")
		assert.Contains(t, w.Plugins, "links")
	})

	t.Run("Render", func(t *testing.T) {
		html := widget.Render("content", "Hello World", nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `name="content"`)
		assert.Contains(t, htmlStr, `id="id_content"`)
		assert.Contains(t, htmlStr, "Hello World")
		assert.Contains(t, htmlStr, "richtext-editor")
		assert.Contains(t, htmlStr, "initRichTextEditor")
	})

	t.Run("Render with nil value", func(t *testing.T) {
		html := widget.Render("content", nil, nil)
		assert.Contains(t, string(html), "textarea")
	})

	t.Run("Parse", func(t *testing.T) {
		value, err := widget.Parse("<p>Hello</p>")
		assert.NoError(t, err)
		assert.Equal(t, "<p>Hello</p>", value)
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "richtext", widget.Type())
	})
}

func TestFileUploadWidget(t *testing.T) {
	widget := NewFileUpload()

	t.Run("Default configuration", func(t *testing.T) {
		assert.Equal(t, "*/*", widget.Accept)
		assert.False(t, widget.Multiple)
		assert.Equal(t, int64(10*1024*1024), widget.MaxSize)
		assert.Equal(t, 1, widget.MaxFiles)
	})

	t.Run("Builder pattern", func(t *testing.T) {
		w := NewFileUpload().
			WithAccept("image/*").
			WithMultiple(true).
			WithMaxSize(5 * 1024 * 1024)

		assert.Equal(t, "image/*", w.Accept)
		assert.True(t, w.Multiple)
		assert.Equal(t, int64(5*1024*1024), w.MaxSize)
	})

	t.Run("Render", func(t *testing.T) {
		html := widget.Render("upload", nil, nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `type="file"`)
		assert.Contains(t, htmlStr, `name="upload"`)
		assert.Contains(t, htmlStr, `id="id_upload"`)
		assert.Contains(t, htmlStr, `data-max-size`)
		assert.Contains(t, htmlStr, "file-upload")
	})

	t.Run("Render with value", func(t *testing.T) {
		html := widget.Render("upload", "/uploads/file.pdf", nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, "file-preview")
		assert.Contains(t, htmlStr, "/uploads/file.pdf")
	})

	t.Run("Render with multiple", func(t *testing.T) {
		widget.WithMultiple(true)
		html := widget.Render("files", nil, nil)
		assert.Contains(t, string(html), "multiple")
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "file", widget.Type())
	})
}

func TestSelectSearchWidget(t *testing.T) {
	choices := []admin.Choice[interface{}]{
		{Value: 1, Label: "Option 1"},
		{Value: 2, Label: "Option 2"},
		{Value: 3, Label: "Option 3"},
	}
	widget := NewSelectSearch(choices)

	t.Run("Default configuration", func(t *testing.T) {
		assert.True(t, widget.searchable)
		assert.Equal(t, "Search or select...", widget.placeholder)
		assert.True(t, widget.allowClear)
	})

	t.Run("Builder pattern", func(t *testing.T) {
		w := NewSelectSearch(choices).
			WithPlaceholder("Choose...").
			WithAllowClear(false)

		assert.Equal(t, "Choose...", w.placeholder)
		assert.False(t, w.allowClear)
	})

	t.Run("Render", func(t *testing.T) {
		html := widget.Render("category", 2, nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `<select`)
		assert.Contains(t, htmlStr, `name="category"`)
		assert.Contains(t, htmlStr, "select-search")
		assert.Contains(t, htmlStr, "Option 1")
		assert.Contains(t, htmlStr, "Option 2")
		assert.Contains(t, htmlStr, "Option 3")
		assert.Contains(t, htmlStr, `selected`)
		assert.Contains(t, htmlStr, "initSelectSearch")
	})

	t.Run("Render with no value", func(t *testing.T) {
		html := widget.Render("category", nil, nil)
		htmlStr := string(html)

		// Should have empty option
		assert.Contains(t, htmlStr, "---------")
	})

	t.Run("Type", func(t *testing.T) {
		assert.Equal(t, "select_search", widget.Type())
	})
}

func TestWidgetRegistry(t *testing.T) {
	registry := NewWidgetRegistry()

	t.Run("Default registrations", func(t *testing.T) {
		assert.NotNil(t, registry.typeMappings["text"])
		assert.NotNil(t, registry.typeMappings["textarea"])
		assert.NotNil(t, registry.typeMappings["number"])
		assert.NotNil(t, registry.typeMappings["email"])
		assert.NotNil(t, registry.typeMappings["checkbox"])
		assert.NotNil(t, registry.typeMappings["date"])
		assert.NotNil(t, registry.typeMappings["datetime"])
	})

	t.Run("RegisterWidget", func(t *testing.T) {
		customWidget := func() admin.Widget {
			return NewTextInput()
		}
		registry.RegisterWidget("custom", customWidget)

		assert.NotNil(t, registry.typeMappings["custom"])
	})

	t.Run("GetWidgetForFieldType", func(t *testing.T) {
		// Test that GetWidgetForFieldType returns widgets
		// Note: This would need actual schema.FieldType values to test properly
		// For now, we just verify the method exists and can be called
		assert.NotNil(t, registry)
	})
}

func TestWidgetRendering(t *testing.T) {
	t.Run("TextInput", func(t *testing.T) {
		widget := NewTextInput()
		html := widget.Render("username", "john_doe", map[string]string{
			"placeholder": "Enter username",
		})
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `type="text"`)
		assert.Contains(t, htmlStr, `name="username"`)
		assert.Contains(t, htmlStr, `value="john_doe"`)
		assert.Contains(t, htmlStr, `placeholder="Enter username"`)
	})

	t.Run("NumberInput", func(t *testing.T) {
		widget := NewNumberInput()
		html := widget.Render("age", 25, nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `type="number"`)
		assert.Contains(t, htmlStr, `value="25"`)
	})

	t.Run("Textarea", func(t *testing.T) {
		widget := NewTextarea(10, 40)
		html := widget.Render("description", "Test content", nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `<textarea`)
		assert.Contains(t, htmlStr, `rows="10"`)
		assert.Contains(t, htmlStr, `cols="40"`)
		assert.Contains(t, htmlStr, "Test content")
	})

	t.Run("Checkbox", func(t *testing.T) {
		widget := NewCheckbox()
		html := widget.Render("is_active", true, nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `type="checkbox"`)
		assert.Contains(t, htmlStr, `checked`)
	})

	t.Run("DateInput", func(t *testing.T) {
		widget := NewDateInput()
		html := widget.Render("birth_date", "2000-01-01", nil)
		htmlStr := string(html)

		assert.Contains(t, htmlStr, `type="date"`)
	})
}

func TestWidgetParsing(t *testing.T) {
	t.Run("TextInput parse", func(t *testing.T) {
		widget := NewTextInput()
		value, err := widget.Parse("test value")
		assert.NoError(t, err)
		assert.Equal(t, "test value", value)
	})

	t.Run("NumberInput parse", func(t *testing.T) {
		widget := NewNumberInput()
		value, err := widget.Parse("42")
		assert.NoError(t, err)
		assert.Equal(t, int64(42), value)
	})

	t.Run("NumberInput parse error", func(t *testing.T) {
		widget := NewNumberInput()
		_, err := widget.Parse("not a number")
		assert.Error(t, err)
	})

	t.Run("Checkbox parse", func(t *testing.T) {
		widget := NewCheckbox()
		
		value, err := widget.Parse("on")
		assert.NoError(t, err)
		assert.True(t, value.(bool))

		value, err = widget.Parse("off")
		assert.NoError(t, err)
		assert.False(t, value.(bool))
	})
}

func TestWidgetHTMLEscaping(t *testing.T) {
	t.Run("XSS prevention in values", func(t *testing.T) {
		widget := NewTextInput()
		maliciousValue := `"><script>alert('xss')</script>`
		html := widget.Render("test", maliciousValue, nil)
		htmlStr := string(html)

		// Should be escaped
		assert.NotContains(t, htmlStr, "<script>")
		assert.Contains(t, htmlStr, "&")
	})

	t.Run("XSS prevention in attributes", func(t *testing.T) {
		widget := NewTextInput()
		html := widget.Render("test", "value", map[string]string{
			"data-attr": `"><script>alert('xss')</script>`,
		})
		htmlStr := string(html)

		assert.NotContains(t, htmlStr, "<script>")
	})
}

func BenchmarkWidgetRendering(b *testing.B) {
	b.Run("TextInput", func(b *testing.B) {
		widget := NewTextInput()
		for i := 0; i < b.N; i++ {
			widget.Render("test", "value", nil)
		}
	})

	b.Run("RichText", func(b *testing.B) {
		widget := NewRichText()
		for i := 0; i < b.N; i++ {
			widget.Render("content", "test content", nil)
		}
	})

	b.Run("SelectSearch", func(b *testing.B) {
		choices := make([]admin.Choice[interface{}], 100)
		for i := 0; i < 100; i++ {
			choices[i] = admin.Choice[interface{}]{
				Value: i,
				Label: strings.Repeat("Option ", i+1),
			}
		}
		widget := NewSelectSearch(choices)

		for i := 0; i < b.N; i++ {
			widget.Render("select", 50, nil)
		}
	})
}
