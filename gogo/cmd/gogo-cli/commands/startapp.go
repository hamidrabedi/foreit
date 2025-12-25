package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var startappCmd = &cobra.Command{
	Use:   "startapp [name]",
	Short: "Create a new app within a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		createApp(appName)
	},
}

func createApp(name string) {
	// Create app directory structure
	appDir := filepath.Join("apps", name)
	
	dirs := []string{
		appDir,
		filepath.Join(appDir, "models"),
		filepath.Join(appDir, "resources"),
		filepath.Join(appDir, "consoles"),
		filepath.Join(appDir, "policies"),
		filepath.Join(appDir, "serializers"),
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}
	
	// Create models/README.md
	modelsReadme := `# Models

Define your Ent schemas here.
`
	writeFile(filepath.Join(appDir, "models", "README.md"), modelsReadme)
	
	// Create resources/README.md
	resourcesReadme := `# Resources

Define your API resource handlers here.
`
	writeFile(filepath.Join(appDir, "resources", "README.md"), resourcesReadme)
	
	// Create consoles/README.md
	consolesReadme := `# Consoles

Define your admin console interfaces here.
`
	writeFile(filepath.Join(appDir, "consoles", "README.md"), consolesReadme)
	
	// Create policies/README.md
	policiesReadme := `# Policies

Define your authorization policies here.
`
	writeFile(filepath.Join(appDir, "policies", "README.md"), policiesReadme)
	
	fmt.Printf("App '%s' created successfully in apps/%s/\n", name, name)
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Printf("Error writing file %s: %v\n", path, err)
		os.Exit(1)
	}
}

// RegisterStartAppCommand registers the startapp command
func RegisterStartAppCommand(rootCmd *cobra.Command) {
	rootCmd.AddCommand(startappCmd)
}

