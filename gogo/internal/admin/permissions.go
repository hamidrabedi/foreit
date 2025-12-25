package admin

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// PermissionChecker checks permissions for admin operations
type PermissionChecker struct {
	registry *Registry
}

// NewPermissionChecker creates a new permission checker
func NewPermissionChecker(registry *Registry) *PermissionChecker {
	return &PermissionChecker{
		registry: registry,
	}
}

// Context represents the request context for permission checking
type Context struct {
	User      interface{} // User object (from auth)
	Request   interface{} // Request object (Fiber context)
	Model     *ModelMeta
	Action    string // "list", "view", "create", "update", "delete"
	Resource  interface{} // The resource being accessed (for view/update/delete)
}

// CheckPermission checks if a user has permission to perform an action
func (pc *PermissionChecker) CheckPermission(ctx *Context) (bool, error) {
	// Check basic permissions first
	if !pc.checkBasicPermissions(ctx) {
		return false, nil
	}
	
	// Check rule-based permissions
	if allowed, err := pc.checkRules(ctx); err != nil {
		return false, err
	} else if !allowed {
		return false, nil
	}
	
	return true, nil
}

// checkBasicPermissions checks basic RBAC permissions
func (pc *PermissionChecker) checkBasicPermissions(ctx *Context) bool {
	switch ctx.Action {
	case "list":
		return ctx.Model.Permissions.CanList
	case "view":
		return ctx.Model.Permissions.CanView
	case "create":
		return ctx.Model.Permissions.CanCreate
	case "update":
		return ctx.Model.Permissions.CanUpdate
	case "delete":
		return ctx.Model.Permissions.CanDelete
	default:
		return false
	}
}

// checkRules checks rule-based permissions (like PocketBase)
func (pc *PermissionChecker) checkRules(ctx *Context) (bool, error) {
	// Get rule for this action
	rule, exists := ctx.Model.Permissions.Rules[ctx.Action]
	if !exists || rule == "" {
		// No rule means allowed (if basic permissions pass)
		return true, nil
	}
	
	// Evaluate rule
	return pc.evaluateRule(rule, ctx)
}

// evaluateRule evaluates a rule expression
func (pc *PermissionChecker) evaluateRule(rule string, ctx *Context) (bool, error) {
	// Evaluate the rule using the simple evaluator
	// The evaluator will handle context variable replacement internally
	return pc.simpleEvaluate(rule, ctx)
}

// replaceContextVariables replaces context variables in rule expressions
func (pc *PermissionChecker) replaceContextVariables(rule string, ctx *Context) string {
	// Replace @request.auth.id
	if strings.Contains(rule, "@request.auth.id") {
		authID := pc.getAuthID(ctx)
		rule = strings.ReplaceAll(rule, "@request.auth.id", fmt.Sprintf("%q", authID))
	}
	
	// Replace @request.user
	if strings.Contains(rule, "@request.user") {
		// This would need to be replaced with actual user object access
		rule = strings.ReplaceAll(rule, "@request.user", "user")
	}
	
	return rule
}

// getAuthID extracts auth ID from context
func (pc *PermissionChecker) getAuthID(ctx *Context) string {
	if ctx.User == nil {
		return ""
	}
	
	// Try to get ID from user object using reflection
	userValue := reflect.ValueOf(ctx.User)
	if userValue.Kind() == reflect.Ptr {
		userValue = userValue.Elem()
	}
	
	// Try common ID field names
	idFields := []string{"ID", "Id", "id", "UserId", "UserID"}
	for _, fieldName := range idFields {
		field := userValue.FieldByName(fieldName)
		if field.IsValid() {
			return fmt.Sprintf("%v", field.Interface())
		}
	}
	
	return ""
}

// simpleEvaluate performs simple rule evaluation
// Supports basic expressions like:
// - Boolean literals: "true", "false"
// - Equality: "field == value", "field != value"
// - Comparisons: "field > value", "field < value", "field >= value", "field <= value"
// - Logical operators: "&&", "||", "!"
// - Parentheses for grouping
// - String literals: "value" or 'value'
// - Numbers: 123, 45.67
// - Context variables: @request.auth.id, @request.user.id
func (pc *PermissionChecker) simpleEvaluate(rule string, ctx *Context) (bool, error) {
	rule = strings.TrimSpace(rule)
	
	// Check for always true/false
	if rule == "true" || rule == "1" {
		return true, nil
	}
	if rule == "false" || rule == "0" || rule == "" {
		return false, nil
	}
	
	// Replace context variables first
	rule = pc.replaceContextVariables(rule, ctx)
	
	// Evaluate the expression
	return pc.evaluateExpression(rule, ctx)
}

// evaluateExpression evaluates a rule expression
func (pc *PermissionChecker) evaluateExpression(expr string, ctx *Context) (bool, error) {
	expr = strings.TrimSpace(expr)
	
	// Handle logical OR (lowest precedence)
	if idx := pc.findOperator(expr, "||"); idx != -1 {
		left, err := pc.evaluateExpression(expr[:idx], ctx)
		if err != nil {
			return false, err
		}
		if left {
			return true, nil
		}
		right, err := pc.evaluateExpression(expr[idx+2:], ctx)
		if err != nil {
			return false, err
		}
		return right, nil
	}
	
	// Handle logical AND
	if idx := pc.findOperator(expr, "&&"); idx != -1 {
		left, err := pc.evaluateExpression(expr[:idx], ctx)
		if err != nil {
			return false, err
		}
		if !left {
			return false, nil
		}
		right, err := pc.evaluateExpression(expr[idx+2:], ctx)
		if err != nil {
			return false, err
		}
		return right, nil
	}
	
	// Handle logical NOT
	if strings.HasPrefix(expr, "!") {
		result, err := pc.evaluateExpression(expr[1:], ctx)
		if err != nil {
			return false, err
		}
		return !result, nil
	}
	
	// Handle parentheses
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return pc.evaluateExpression(expr[1:len(expr)-1], ctx)
	}
	
	// Handle comparison operators
	operators := []string{"!=", "==", ">=", "<=", ">", "<"}
	for _, op := range operators {
		if idx := strings.Index(expr, op); idx != -1 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+len(op):])
			
			leftVal := pc.getValue(left, ctx)
			rightVal := pc.getValue(right, ctx)
			
			return pc.compareValues(leftVal, rightVal, op), nil
		}
	}
	
	// If no operator found, treat as a boolean value
	val := pc.getValue(expr, ctx)
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	if strVal, ok := val.(string); ok {
		return strVal != "" && strVal != "false" && strVal != "0", nil
	}
	
	// Default: if value exists and is truthy, return true
	return val != nil, nil
}

// findOperator finds an operator in an expression, respecting parentheses
func (pc *PermissionChecker) findOperator(expr string, op string) int {
	depth := 0
	for i := 0; i < len(expr)-len(op)+1; i++ {
		if expr[i] == '(' {
			depth++
		} else if expr[i] == ')' {
			depth--
		} else if depth == 0 && strings.HasPrefix(expr[i:], op) {
			return i
		}
	}
	return -1
}

// getValue gets a value from an expression (field reference, literal, or context variable)
func (pc *PermissionChecker) getValue(expr string, ctx *Context) interface{} {
	expr = strings.TrimSpace(expr)
	
	// String literal
	if (strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`)) ||
		(strings.HasPrefix(expr, `'`) && strings.HasSuffix(expr, `'`)) {
		return strings.Trim(expr, `"'`)
	}
	
	// Boolean literal
	if expr == "true" {
		return true
	}
	if expr == "false" {
		return false
	}
	
	// Number
	if num, err := strconv.ParseFloat(expr, 64); err == nil {
		return num
	}
	if num, err := strconv.ParseInt(expr, 10, 64); err == nil {
		return num
	}
	
	// Field reference - get from resource
	if ctx.Resource != nil {
		val := pc.getFieldValue(ctx.Resource, expr)
		if val != nil {
			return val
		}
	}
	
	// Context variable (should have been replaced, but handle just in case)
	if strings.HasPrefix(expr, "@") {
		return pc.getContextVariable(expr, ctx)
	}
	
	// Default: return as string
	return expr
}

// getFieldValue gets a field value from a resource using reflection
func (pc *PermissionChecker) getFieldValue(resource interface{}, fieldName string) interface{} {
	val := reflect.ValueOf(resource)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}
	
	if val.Kind() != reflect.Struct {
		return nil
	}
	
	// Try exact field name
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface()
	}
	
	// Try with different capitalizations
	fieldNameLower := strings.ToLower(fieldName)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if strings.ToLower(field.Name) == fieldNameLower {
			return val.Field(i).Interface()
		}
	}
	
	return nil
}

// getContextVariable gets a value from context variables
func (pc *PermissionChecker) getContextVariable(varName string, ctx *Context) interface{} {
	// @request.auth.id
	if varName == "@request.auth.id" {
		return pc.getAuthID(ctx)
	}
	
	// @request.user.*
	if strings.HasPrefix(varName, "@request.user.") {
		fieldName := strings.TrimPrefix(varName, "@request.user.")
		return pc.getFieldValue(ctx.User, fieldName)
	}
	
	// @request.auth.* (for other auth fields)
	if strings.HasPrefix(varName, "@request.auth.") {
		fieldName := strings.TrimPrefix(varName, "@request.auth.")
		return pc.getFieldValue(ctx.User, fieldName)
	}
	
	return nil
}

// compareValues compares two values using an operator
func (pc *PermissionChecker) compareValues(left, right interface{}, op string) bool {
	// Convert to comparable types
	leftVal := pc.normalizeValue(left)
	rightVal := pc.normalizeValue(right)
	
	switch op {
	case "==":
		return reflect.DeepEqual(leftVal, rightVal)
	case "!=":
		return !reflect.DeepEqual(leftVal, rightVal)
	case ">", ">=", "<", "<=":
		return pc.compareNumeric(leftVal, rightVal, op)
	default:
		return false
	}
}

// normalizeValue normalizes a value for comparison
func (pc *PermissionChecker) normalizeValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	
	// Convert string numbers to numbers
	if str, ok := val.(string); ok {
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return num
		}
		if num, err := strconv.ParseInt(str, 10, 64); err == nil {
			return num
		}
		if str == "true" {
			return true
		}
		if str == "false" {
			return false
		}
		return str
	}
	
	return val
}

// compareNumeric compares two numeric values
func (pc *PermissionChecker) compareNumeric(left, right interface{}, op string) bool {
	leftNum := pc.toFloat(left)
	rightNum := pc.toFloat(right)
	
	switch op {
	case ">":
		return leftNum > rightNum
	case ">=":
		return leftNum >= rightNum
	case "<":
		return leftNum < rightNum
	case "<=":
		return leftNum <= rightNum
	default:
		return false
	}
}

// toFloat converts a value to float64
func (pc *PermissionChecker) toFloat(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num
		}
	}
	return 0
}

// CheckFieldPermission checks if a user can access a specific field
func (pc *PermissionChecker) CheckFieldPermission(ctx *Context, fieldName string) bool {
	// Check if field is excluded
	for _, excluded := range ctx.Model.Options.ExcludeFields {
		if excluded == fieldName {
			return false
		}
	}
	
	// Check if field is read-only and action is update
	if ctx.Action == "update" || ctx.Action == "create" {
		for _, readOnly := range ctx.Model.Options.ReadOnlyFields {
			if readOnly == fieldName {
				return false
			}
		}
	}
	
	return true
}

// FilterFields filters fields based on permissions
func (pc *PermissionChecker) FilterFields(ctx *Context, fields []FieldMeta) []FieldMeta {
	filtered := make([]FieldMeta, 0, len(fields))
	
	for _, field := range fields {
		if pc.CheckFieldPermission(ctx, field.Name) {
			filtered = append(filtered, field)
		}
	}
	
	return filtered
}

// RuleEvaluator is an interface for rule evaluation engines
type RuleEvaluator interface {
	Evaluate(rule string, context map[string]interface{}) (bool, error)
}

// SimpleRuleEvaluator is a basic rule evaluator
type SimpleRuleEvaluator struct{}

// Evaluate evaluates a rule with the given context
func (e *SimpleRuleEvaluator) Evaluate(rule string, context map[string]interface{}) (bool, error) {
	// Placeholder implementation
	// In production, integrate with a proper expression evaluator like:
	// - github.com/antonmedv/expr
	// - github.com/Knetic/govaluate
	return true, nil
}

// GetUserFromContext extracts user from Fiber context
func GetUserFromContext(fiberCtx interface{}) interface{} {
	// This will be implemented to extract user from Fiber context
	// Typically from JWT token, session, etc.
	return nil
}

// GetRequestInfo extracts request information for rule evaluation
func GetRequestInfo(fiberCtx interface{}) map[string]interface{} {
	info := make(map[string]interface{})
	
	// Extract common request info
	// This will be populated from Fiber context
	
	return info
}

