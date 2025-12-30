package admin

// GetGlobalRegistry returns the global admin registry
// This allows external packages to access the registry for HTTP routing
func GetGlobalRegistry() *Registry {
	return globalRegistry
}
