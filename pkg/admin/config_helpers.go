package admin

// getBoolConfig gets a boolean config value with default
func getBoolConfig(model *AdminModel, key string, defaultValue bool) bool {
	if model.ExtendedConfig == nil {
		return defaultValue
	}
	if val, ok := model.ExtendedConfig[key].(bool); ok {
		return val
	}
	return defaultValue
}

// getStringConfig gets a string config value with default
func getStringConfig(model *AdminModel, key string, defaultValue string) string {
	if model.ExtendedConfig == nil {
		return defaultValue
	}
	if val, ok := model.ExtendedConfig[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

