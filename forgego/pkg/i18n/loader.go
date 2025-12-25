package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type JSONLoader struct {
	path string
}

func NewJSONLoader(path string) *JSONLoader {
	return &JSONLoader{
		path: path,
	}
}

func (jl *JSONLoader) Load() (map[string]*Bundle, error) {
	bundles := make(map[string]*Bundle)
	
	err := filepath.Walk(jl.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}
		
		locale := strings.TrimSuffix(info.Name(), ".json")
		bundle := NewBundle(locale)
		
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		var messages map[string]interface{}
		if err := json.Unmarshal(data, &messages); err != nil {
			return err
		}
		
		for key, value := range messages {
			switch v := value.(type) {
			case string:
				bundle.AddMessage(key, v)
			case map[string]interface{}:
				if singular, ok := v["one"].(string); ok {
					plural := ""
					if p, ok := v["other"].(string); ok {
						plural = p
					} else if p, ok := v["many"].(string); ok {
						plural = p
					}
					bundle.AddPluralMessage(key, singular, plural)
				}
			}
		}
		
		bundles[locale] = bundle
		return nil
	})
	
	return bundles, err
}

