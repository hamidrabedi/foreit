package orm

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"

	forgeerrors "github.com/forgego/forge/errors"
)

// AggregateValuer is an optional capability that QuerySet implementations may
// provide to support AggregateValues.
type AggregateValuer interface {
	AggregateValues(context.Context, ...Aggregate) (map[string]any, error)
}

type resolvedAggregate struct {
	name      string
	function  string
	field     string
	column    string
	fieldType reflect.Type
	countStar bool
}

// AggregateValues executes ungrouped aggregates for any QuerySet implementation
// that supports them.
func AggregateValues[T any](ctx context.Context, qs QuerySet[T], aggs ...Aggregate) (map[string]any, error) {
	valuer, ok := qs.(AggregateValuer)
	if !ok {
		return nil, forgeerrors.NewNotImplementedError("AggregateValues for this QuerySet implementation")
	}
	return valuer.AggregateValues(ctx, aggs...)
}

// AggregateValues executes one ungrouped aggregate query. Ordering, limits,
// and offsets are intentionally ignored.
func (qs *BaseQuerySet[T]) AggregateValues(ctx context.Context, aggs ...Aggregate) (map[string]any, error) {
	return qs.aggregateValues(ctx, aggs...)
}

func (qs *BaseQuerySet[T]) aggregateValues(ctx context.Context, aggs ...Aggregate) (map[string]any, error) {
	if qs.err != nil && !isAggregateChainError(qs.err) {
		return nil, qs.err
	}
	resolved, err := qs.resolveAggregates(aggs)
	if err != nil {
		return nil, err
	}
	query, args, err := qs.buildAggregateSQL(resolved)
	if err != nil {
		return nil, err
	}
	database, err := qs.getDB(ctx)
	if err != nil {
		return nil, err
	}

	values := make([]any, len(resolved))
	destinations := make([]any, len(values))
	for i := range values {
		destinations[i] = &values[i]
	}
	if err := database.QueryRowContext(ctx, query, args...).Scan(destinations...); err != nil {
		return nil, fmt.Errorf("aggregate query failed: %w", err)
	}

	result := make(map[string]any, len(resolved))
	for i, aggregate := range resolved {
		value, err := convertAggregateValue(aggregate, values[i])
		if err != nil {
			return nil, err
		}
		result[aggregate.name] = value
	}
	return result, nil
}

func (qs *BaseQuerySet[T]) resolveAggregates(aggs []Aggregate) ([]resolvedAggregate, error) {
	if len(aggs) == 0 {
		return nil, forgeerrors.NewInvalidInputError("aggregates", "at least one aggregate is required")
	}
	resolved := make([]resolvedAggregate, 0, len(aggs))
	names := make(map[string]bool, len(aggs))
	for _, aggregate := range aggs {
		function := strings.ToUpper(strings.TrimSpace(aggregate.Func))
		switch function {
		case string(AggCount), string(AggSum), string(AggAvg), string(AggMin), string(AggMax):
		default:
			return nil, forgeerrors.NewNotImplementedError("aggregate " + aggregate.Func)
		}
		name := aggregate.Name
		if name == "" {
			name = strings.ToLower(function) + "_" + aggregate.Field
		}
		if names[name] {
			return nil, forgeerrors.NewInvalidInputError("aggregates", "duplicate aggregate name "+name)
		}
		names[name] = true

		item := resolvedAggregate{name: name, function: function, field: aggregate.Field}
		if function == string(AggCount) && (aggregate.Field == "" || aggregate.Field == "*") {
			item.countStar = true
		} else {
			field, _, err := qs.schema.ResolvePath(aggregate.Field)
			if err != nil {
				return nil, forgeerrors.NewInvalidInputError(aggregate.Field, "field does not resolve to a model column")
			}
			item.column = EscapeIdentifier(field.DBColumn)
			item.fieldType = field.Type
		}
		resolved = append(resolved, item)
	}
	return resolved, nil
}

func convertAggregateValue(aggregate resolvedAggregate, value any) (any, error) {
	if aggregate.function == string(AggMin) || aggregate.function == string(AggMax) {
		if bytes, ok := value.([]byte); ok {
			if aggregate.fieldType != nil && aggregate.fieldType.Kind() == reflect.String {
				return string(bytes), nil
			}
			return bytes, nil
		}
		return value, nil
	}
	if value == nil {
		if aggregate.function == string(AggCount) {
			return int64(0), nil
		}
		return nil, nil
	}
	var text string
	switch number := value.(type) {
	case int64:
		if aggregate.function == string(AggCount) {
			return number, nil
		}
		return float64(number), nil
	case float64:
		if aggregate.function == string(AggCount) {
			return int64(number), nil
		}
		return number, nil
	case []byte:
		text = string(number)
	case string:
		text = number
	default:
		text = fmt.Sprint(value)
	}
	if aggregate.function == string(AggCount) {
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("aggregate %s: convert COUNT result: %w", aggregate.name, err)
		}
		return parsed, nil
	}
	parsed, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return nil, fmt.Errorf("aggregate %s: convert numeric result: %w", aggregate.name, err)
	}
	return parsed, nil
}

// Aggregate represents an aggregate function
type Aggregate struct {
	Name  string
	Field string
	Func  string // COUNT, SUM, AVG, MAX, MIN, etc.
}

// AggregateFunc represents an aggregate function
type AggregateFunc string

const (
	AggCount    AggregateFunc = "COUNT"
	AggSum      AggregateFunc = "SUM"
	AggAvg      AggregateFunc = "AVG"
	AggMax      AggregateFunc = "MAX"
	AggMin      AggregateFunc = "MIN"
	AggStdDev   AggregateFunc = "STDDEV"
	AggVariance AggregateFunc = "VARIANCE"
)

var (
	aggregateRegistryMu sync.RWMutex
	aggregateRegistry   = map[string]aggregateRegistration{}
)

type aggregateRegistration struct {
	funcName string
	builder  func(string) Aggregate
}

// Count creates a COUNT aggregate
func Count(field string) Aggregate {
	return Aggregate{
		Name:  "count",
		Field: field,
		Func:  string(AggCount),
	}
}

// Sum creates a SUM aggregate
func Sum(field string) Aggregate {
	return Aggregate{
		Name:  "sum",
		Field: field,
		Func:  string(AggSum),
	}
}

// Avg creates an AVG aggregate
func Avg(field string) Aggregate {
	return Aggregate{
		Name:  "avg",
		Field: field,
		Func:  string(AggAvg),
	}
}

// Max creates a MAX aggregate
func Max(field string) Aggregate {
	return Aggregate{
		Name:  "max",
		Field: field,
		Func:  string(AggMax),
	}
}

// Min creates a MIN aggregate
func Min(field string) Aggregate {
	return Aggregate{
		Name:  "min",
		Field: field,
		Func:  string(AggMin),
	}
}

// StdDev creates a STDDEV aggregate
func StdDev(field string) Aggregate {
	return Aggregate{
		Name:  "stddev",
		Field: field,
		Func:  string(AggStdDev),
	}
}

// Variance creates a VARIANCE aggregate
func Variance(field string) Aggregate {
	return Aggregate{
		Name:  "variance",
		Field: field,
		Func:  string(AggVariance),
	}
}

// RegisterAggregate registers a custom aggregate function
func RegisterAggregate(name, funcName string, builder func(string) Aggregate) {
	if builder == nil || strings.TrimSpace(name) == "" {
		return
	}

	aggregateRegistryMu.Lock()
	defer aggregateRegistryMu.Unlock()
	aggregateRegistry[normalizeRegistryName(name)] = aggregateRegistration{
		funcName: strings.TrimSpace(funcName),
		builder:  builder,
	}
}

// BuildAggregate builds an aggregate from a registered custom aggregate name.
func BuildAggregate(name, field string) (Aggregate, bool) {
	aggregateRegistryMu.RLock()
	registration, ok := aggregateRegistry[normalizeRegistryName(name)]
	aggregateRegistryMu.RUnlock()
	if !ok {
		return Aggregate{}, false
	}

	aggregate := registration.builder(field)
	if aggregate.Name == "" {
		aggregate.Name = normalizeRegistryName(name)
	}
	if aggregate.Field == "" {
		aggregate.Field = field
	}
	if aggregate.Func == "" {
		aggregate.Func = registration.funcName
	}
	return aggregate, true
}
