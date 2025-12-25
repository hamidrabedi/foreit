package mypackage

import "log"

// MyPackage represents a third-party package that integrates with Forge
type MyPackage struct {
	config *MyPackageConfig
}

// New creates a new MyPackage instance
func New(config *MyPackageConfig) *MyPackage {
	return &MyPackage{
		config: config,
	}
}

// Start initializes the package
func (p *MyPackage) Start() error {
	if !p.config.Enable {
		log.Println("MyPackage is disabled")
		return nil
	}
	
	log.Printf("MyPackage started at path: %s", p.config.Path)
	log.Printf("MyPackage API Key: %s", p.config.APIKey)
	return nil
}

