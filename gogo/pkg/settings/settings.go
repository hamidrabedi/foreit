package settings

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// Loader loads settings from various sources
type Loader struct {
	sources []Source
}

// Source represents a configuration source
type Source interface {
	Get(key string) (string, bool)
}

// EnvSource loads from environment variables
type EnvSource struct{}

func (e *EnvSource) Get(key string) (string, bool) {
	value := os.Getenv(key)
	return value, value != ""
}

// FileSource loads from a file
type FileSource struct {
	path   string
	values map[string]string
}

// NewFileSource creates a new file source
func NewFileSource(path string) (*FileSource, error) {
	fs := &FileSource{
		path:   path,
		values: make(map[string]string),
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			fs.values[key] = value
		}
	}
	
	return fs, nil
}

func (f *FileSource) Get(key string) (string, bool) {
	value, ok := f.values[key]
	return value, ok
}

// NewLoader creates a new settings loader
func NewLoader(sources ...Source) *Loader {
	if len(sources) == 0 {
		sources = []Source{&EnvSource{}}
	}
	return &Loader{
		sources: sources,
	}
}

// Load loads settings into a struct
func Load[T any]() (*T, error) {
	loader := NewLoader()
	return loader.Load[T]()
}

// Load loads settings into the provided type
func (l *Loader) Load[T any]() (*T, error) {
	var zero T
	config := new(T)
	
	typ := reflect.TypeOf(zero)
	val := reflect.ValueOf(config).Elem()
	
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)
		
		// Get env tag
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}
		
		// Get setting tag (alternative key)
		settingTag := field.Tag.Get("setting")
		if settingTag != "" {
			envTag = settingTag
		}
		
		// Try to get value from sources
		var value string
		var found bool
		for _, source := range l.sources {
			if v, ok := source.Get(envTag); ok {
				value = v
				found = true
				break
			}
		}
		
		if !found {
			// Check for default
			defaultTag := field.Tag.Get("default")
			if defaultTag != "" {
				value = defaultTag
				found = true
			} else {
				// Check if required
				required := field.Tag.Get("required") == "true"
				if required {
					return nil, fmt.Errorf("required setting %s is missing", envTag)
				}
				continue
			}
		}
		
		// Convert and set value
		if err := setFieldValue(fieldVal, value); err != nil {
			return nil, fmt.Errorf("failed to set %s: %w", envTag, err)
		}
	}
	
	return config, nil
}

// setFieldValue sets a field value from a string
func setFieldValue(field reflect.Value, value string) error {
	if !field.CanSet() {
		return fmt.Errorf("field cannot be set")
	}
	
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(boolVal)
	case reflect.Slice:
		// Handle slice of strings (comma-separated)
		if field.Type().Elem().Kind() == reflect.String {
			parts := strings.Split(value, ",")
			slice := reflect.MakeSlice(field.Type(), len(parts), len(parts))
			for i, part := range parts {
				slice.Index(i).SetString(strings.TrimSpace(part))
			}
			field.Set(slice)
		}
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	
	return nil
}

// Registry allows registering and accessing settings globally
type Registry struct {
	settings map[string]interface{}
}

var globalRegistry = &Registry{
	settings: make(map[string]interface{}),
}

// Register registers a setting value
func Register(key string, value interface{}) {
	globalRegistry.settings[key] = value
}

// Get retrieves a setting value
func Get(key string) (interface{}, bool) {
	value, ok := globalRegistry.settings[key]
	return value, ok
}

// GetString retrieves a string setting
func GetString(key string, defaultValue string) string {
	if value, ok := Get(key); ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt retrieves an int setting
func GetInt(key string, defaultValue int) int {
	if value, ok := Get(key); ok {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return defaultValue
}

// GetBool retrieves a bool setting
func GetBool(key string, defaultValue bool) bool {
	if value, ok := Get(key); ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}

