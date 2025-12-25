package admin

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type FilterHandler interface {
	GetChoices() []Choice
	GetLookups() []Lookup
	GetType() FilterType
	ApplyFilter(query interface{}, value interface{}, lookup Lookup) (interface{}, error)
}

type FilterType string

const (
	FilterTypeList      FilterType = "list"
	FilterTypeDate      FilterType = "date"
	FilterTypeDateTime  FilterType = "datetime"
	FilterTypeBoolean   FilterType = "boolean"
	FilterTypeChoice    FilterType = "choice"
	FilterTypeRelated   FilterType = "related"
	FilterTypeRelatedOnly FilterType = "related_only"
)

type Lookup string

const (
	LookupExact      Lookup = "exact"
	LookupIExact     Lookup = "iexact"
	LookupContains   Lookup = "contains"
	LookupIContains  Lookup = "icontains"
	LookupIn         Lookup = "in"
	LookupGt         Lookup = "gt"
	LookupGte        Lookup = "gte"
	LookupLt         Lookup = "lt"
	LookupLte        Lookup = "lte"
	LookupStartswith Lookup = "startswith"
	LookupIStartswith Lookup = "istartswith"
	LookupEndswith   Lookup = "endswith"
	LookupIEndswith  Lookup = "iendswith"
	LookupRange      Lookup = "range"
	LookupDate       Lookup = "date"
	LookupYear       Lookup = "year"
	LookupMonth      Lookup = "month"
	LookupDay        Lookup = "day"
	LookupWeek       Lookup = "week"
	LookupWeekDay    Lookup = "week_day"
	LookupTime       Lookup = "time"
	LookupHour       Lookup = "hour"
	LookupMinute     Lookup = "minute"
	LookupSecond     Lookup = "second"
	LookupIsNull     Lookup = "isnull"
	LookupRegex      Lookup = "regex"
	LookupIRegex      Lookup = "iregex"
)

type DateRangeFilter struct {
	Field string
	Label string
}

func NewDateRangeFilter(field, label string) *DateRangeFilter {
	return &DateRangeFilter{
		Field: field,
		Label: label,
	}
}

func (drf *DateRangeFilter) GetChoices() []Choice {
	return []Choice{}
}

func (drf *DateRangeFilter) GetLookups() []Lookup {
	return []Lookup{
		LookupExact,
		LookupGte,
		LookupLte,
		LookupRange,
		LookupYear,
		LookupMonth,
		LookupDay,
	}
}

func (drf *DateRangeFilter) GetType() FilterType {
	return FilterTypeDate
}

func (drf *DateRangeFilter) ApplyFilter(query interface{}, value interface{}, lookup Lookup) (interface{}, error) {
	return query, nil
}

type ChoiceFilter struct {
	Field   string
	Label   string
	Choices []Choice
}

func NewChoiceFilter(field, label string, choices []Choice) *ChoiceFilter {
	return &ChoiceFilter{
		Field:   field,
		Label:   label,
		Choices: choices,
	}
}

func (cf *ChoiceFilter) GetChoices() []Choice {
	return cf.Choices
}

func (cf *ChoiceFilter) GetLookups() []Lookup {
	return []Lookup{
		LookupExact,
		LookupIn,
	}
}

func (cf *ChoiceFilter) GetType() FilterType {
	return FilterTypeChoice
}

func (cf *ChoiceFilter) ApplyFilter(query interface{}, value interface{}, lookup Lookup) (interface{}, error) {
	return query, nil
}

type BooleanFilter struct {
	Field string
	Label string
}

func NewBooleanFilter(field, label string) *BooleanFilter {
	return &BooleanFilter{
		Field: field,
		Label: label,
	}
}

func (bf *BooleanFilter) GetChoices() []Choice {
	return []Choice{
		{Value: "true", Label: "Yes"},
		{Value: "false", Label: "No"},
	}
}

func (bf *BooleanFilter) GetLookups() []Lookup {
	return []Lookup{LookupExact}
}

func (bf *BooleanFilter) GetType() FilterType {
	return FilterTypeBoolean
}

func (bf *BooleanFilter) ApplyFilter(query interface{}, value interface{}, lookup Lookup) (interface{}, error) {
	return query, nil
}

type RelatedFilter struct {
	Field       string
	Label       string
	RelatedModel string
	DisplayField string
}

func NewRelatedFilter(field, label, relatedModel, displayField string) *RelatedFilter {
	return &RelatedFilter{
		Field:        field,
		Label:        label,
		RelatedModel: relatedModel,
		DisplayField: displayField,
	}
}

func (rf *RelatedFilter) GetChoices() []Choice {
	return []Choice{}
}

func (rf *RelatedFilter) GetLookups() []Lookup {
	return []Lookup{
		LookupExact,
		LookupIn,
		LookupIsNull,
	}
}

func (rf *RelatedFilter) GetType() FilterType {
	return FilterTypeRelated
}

func (rf *RelatedFilter) ApplyFilter(query interface{}, value interface{}, lookup Lookup) (interface{}, error) {
	return query, nil
}

type FilterRegistry struct {
	filters map[string]FilterHandler
}

func NewFilterRegistry() *FilterRegistry {
	return &FilterRegistry{
		filters: make(map[string]FilterHandler),
	}
}

func (fr *FilterRegistry) Register(field string, handler FilterHandler) {
	fr.filters[field] = handler
}

func (fr *FilterRegistry) Get(field string) (FilterHandler, bool) {
	handler, ok := fr.filters[field]
	return handler, ok
}

func ParseDateRange(value string) (time.Time, time.Time, error) {
	parts := strings.Split(value, ",")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range format")
	}
	
	start, err := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	end, err := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	
	return start, end, nil
}

func ParseFilterValue(value string, filterType FilterType) (interface{}, error) {
	switch filterType {
	case FilterTypeBoolean:
		return strconv.ParseBool(value)
	case FilterTypeDate, FilterTypeDateTime:
		return time.Parse("2006-01-02", value)
	case FilterTypeChoice:
		return value, nil
	default:
		return value, nil
	}
}

func BuildFilterSpec(field string, handler FilterHandler) map[string]interface{} {
	return map[string]interface{}{
		"field":    field,
		"type":     handler.GetType(),
		"lookups":  handler.GetLookups(),
		"choices":  handler.GetChoices(),
		"label":    field,
	}
}

