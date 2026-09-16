package utils

import "testing"

func TestPluralize(t *testing.T) {
	tests := map[string]string{
		"category":        "categories",
		"box":             "boxes",
		"product_variant": "product_variants",
		"address":         "addresses",
		"day":             "days",
	}
	for singular, want := range tests {
		t.Run(singular, func(t *testing.T) {
			if got := Pluralize(singular); got != want {
				t.Fatalf("Pluralize(%q) = %q, want %q", singular, got, want)
			}
		})
	}
}
