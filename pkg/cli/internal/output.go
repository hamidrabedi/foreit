package internal

import "fmt"

// PrintSuccess prints a success message with a checkmark
func PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Printf("❌ %s\n", message)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Printf("⚠️  %s\n", message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}

