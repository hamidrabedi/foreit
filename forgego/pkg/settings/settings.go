package settings

// Load is now implemented in loader.go using Viper
// This file maintains backward compatibility with the Registry pattern

type Registry struct {
	settings map[string]interface{}
}

var globalRegistry = &Registry{
	settings: make(map[string]interface{}),
}

func Register(key string, value interface{}) {
	globalRegistry.settings[key] = value
}

func Get(key string) (interface{}, bool) {
	value, ok := globalRegistry.settings[key]
	return value, ok
}

func GetString(key string, defaultValue string) string {
	if value, ok := Get(key); ok {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	if value, ok := Get(key); ok {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return defaultValue
}

func GetBool(key string, defaultValue bool) bool {
	if value, ok := Get(key); ok {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}
