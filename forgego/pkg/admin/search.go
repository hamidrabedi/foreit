package admin

import (
	"fmt"
	"strings"
)

type SearchEngine interface {
	Search(query string, fields []string, modelMeta *ModelMeta) (interface{}, error)
	FullTextSearch(query string, fields []string, modelMeta *ModelMeta) (interface{}, error)
}

type BasicSearchEngine struct {
	db interface{}
}

func NewBasicSearchEngine(db interface{}) *BasicSearchEngine {
	return &BasicSearchEngine{db: db}
}

func (bse *BasicSearchEngine) Search(query string, fields []string, modelMeta *ModelMeta) (interface{}, error) {
	if query == "" || len(fields) == 0 {
		return nil, nil
	}
	
	searchTerms := strings.Fields(query)
	conditions := make([]string, 0)
	args := make([]interface{}, 0)
	
	for _, field := range fields {
		if !isFieldSearchable(field, modelMeta) {
			continue
		}
		
		fieldConditions := make([]string, 0)
		for _, term := range searchTerms {
			fieldConditions = append(fieldConditions, fmt.Sprintf("%s LIKE ?", field))
			args = append(args, "%"+term+"%")
		}
		
		if len(fieldConditions) > 0 {
			conditions = append(conditions, "("+strings.Join(fieldConditions, " AND ")+")")
		}
	}
	
	if len(conditions) == 0 {
		return nil, nil
	}
	
	whereClause := strings.Join(conditions, " OR ")
	return map[string]interface{}{
		"where": whereClause,
		"args":  args,
	}, nil
}

func (bse *BasicSearchEngine) FullTextSearch(query string, fields []string, modelMeta *ModelMeta) (interface{}, error) {
	if query == "" || len(fields) == 0 {
		return nil, nil
	}
	
	searchTerms := strings.Fields(query)
	searchQuery := strings.Join(searchTerms, " & ")
	
	tsvectorFields := make([]string, 0)
	for _, field := range fields {
		if !isFieldSearchable(field, modelMeta) {
			continue
		}
		tsvectorFields = append(tsvectorFields, fmt.Sprintf("to_tsvector('english', %s)", field))
	}
	
	if len(tsvectorFields) == 0 {
		return nil, nil
	}
	
	tsvector := strings.Join(tsvectorFields, " || ")
	
	return map[string]interface{}{
		"where": fmt.Sprintf("%s @@ to_tsquery('english', ?)", tsvector),
		"args":  []interface{}{searchQuery},
	}, nil
}

type PostgreSQLSearchEngine struct {
	db interface{}
}

func NewPostgreSQLSearchEngine(db interface{}) *PostgreSQLSearchEngine {
	return &PostgreSQLSearchEngine{db: db}
}

func (pse *PostgreSQLSearchEngine) Search(query string, fields []string, modelMeta *ModelMeta) (interface{}, error) {
	return pse.FullTextSearch(query, fields, modelMeta)
}

func (pse *PostgreSQLSearchEngine) FullTextSearch(query string, fields []string, modelMeta *ModelMeta) (interface{}, error) {
	if query == "" || len(fields) == 0 {
		return nil, nil
	}
	
	searchTerms := strings.Fields(query)
	searchQuery := strings.Join(searchTerms, " & ")
	
	tsvectorFields := make([]string, 0)
	for _, field := range fields {
		if !isFieldSearchable(field, modelMeta) {
			continue
		}
		tsvectorFields = append(tsvectorFields, fmt.Sprintf("to_tsvector('english', COALESCE(%s::text, ''))", field))
	}
	
	if len(tsvectorFields) == 0 {
		return nil, nil
	}
	
	tsvector := strings.Join(tsvectorFields, " || ")
	
	return map[string]interface{}{
		"where": fmt.Sprintf("%s @@ to_tsquery('english', ?)", tsvector),
		"args":  []interface{}{searchQuery},
	}, nil
}

func isFieldSearchable(fieldName string, modelMeta *ModelMeta) bool {
	if len(modelMeta.Options.SearchFields) > 0 {
		for _, f := range modelMeta.Options.SearchFields {
			if f == fieldName {
				return true
			}
		}
		return false
	}
	
	for _, field := range modelMeta.Fields {
		if field.Name == fieldName && field.Searchable {
			return true
		}
	}
	
	return false
}

func BuildSearchQuery(query string, fields []string, modelMeta *ModelMeta, engine SearchEngine) (interface{}, error) {
	if engine == nil {
		engine = NewBasicSearchEngine(nil)
	}
	
	return engine.Search(query, fields, modelMeta)
}

func BuildFullTextSearchQuery(query string, fields []string, modelMeta *ModelMeta, engine SearchEngine) (interface{}, error) {
	if engine == nil {
		engine = NewPostgreSQLSearchEngine(nil)
	}
	
	return engine.FullTextSearch(query, fields, modelMeta)
}

