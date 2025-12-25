package admin

// FilterType represents filter type
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

// Lookup represents a lookup type
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

// Choice represents a choice option
type Choice struct {
	Value interface{}
	Label string
}

