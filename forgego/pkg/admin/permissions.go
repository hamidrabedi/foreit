package admin

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type PermissionChecker struct {
	registry *Registry
}

func NewPermissionChecker(registry *Registry) *PermissionChecker {
	return &PermissionChecker{
		registry: registry,
	}
}

type Context struct {
	User      interface{}
	Request   interface{}
	Model     *ModelMeta
	Action    string
	Resource  interface{}
	Field     string
}

func (pc *PermissionChecker) CheckPermission(ctx *Context) (bool, error) {
	if !pc.checkBasicPermissions(ctx) {
		return false, nil
	}
	
	if allowed, err := pc.checkRules(ctx); err != nil {
		return false, err
	} else if !allowed {
		return false, nil
	}
	
	return true, nil
}
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

func (pc *PermissionChecker) checkRules(ctx *Context) (bool, error) {
	rule, exists := ctx.Model.Permissions.Rules[ctx.Action]
	if !exists || rule == "" {
		return true, nil
	}
	
	return pc.evaluateRule(rule, ctx)
}

func (pc *PermissionChecker) evaluateRule(rule string, ctx *Context) (bool, error) {
	return pc.simpleEvaluate(rule, ctx)
}

func (pc *PermissionChecker) replaceContextVariables(rule string, ctx *Context) string {
	if strings.Contains(rule, "@request.auth.id") {
		authID := pc.getAuthID(ctx)
		rule = strings.ReplaceAll(rule, "@request.auth.id", fmt.Sprintf("%q", authID))
	}
	
	if strings.Contains(rule, "@request.user") {
		rule = strings.ReplaceAll(rule, "@request.user", "user")
	}
	
	return rule
}

func (pc *PermissionChecker) getAuthID(ctx *Context) string {
	if ctx.User == nil {
		return ""
	}
	
	userValue := reflect.ValueOf(ctx.User)
	if userValue.Kind() == reflect.Ptr {
		userValue = userValue.Elem()
	}
	
	idFields := []string{"ID", "Id", "id", "UserId", "UserID"}
	for _, fieldName := range idFields {
		field := userValue.FieldByName(fieldName)
		if field.IsValid() {
			return fmt.Sprintf("%v", field.Interface())
		}
	}
	
	return ""
}

func (pc *PermissionChecker) simpleEvaluate(rule string, ctx *Context) (bool, error) {
	rule = strings.TrimSpace(rule)
	
	if rule == "true" || rule == "1" {
		return true, nil
	}
	if rule == "false" || rule == "0" || rule == "" {
		return false, nil
	}
	
	rule = pc.replaceContextVariables(rule, ctx)
	return pc.evaluateExpression(rule, ctx)
}

func (pc *PermissionChecker) evaluateExpression(expr string, ctx *Context) (bool, error) {
	expr = strings.TrimSpace(expr)
	
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
	
	if strings.HasPrefix(expr, "!") {
		result, err := pc.evaluateExpression(expr[1:], ctx)
		if err != nil {
			return false, err
		}
		return !result, nil
	}
	
	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return pc.evaluateExpression(expr[1:len(expr)-1], ctx)
	}
	
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
	
	val := pc.getValue(expr, ctx)
	if boolVal, ok := val.(bool); ok {
		return boolVal, nil
	}
	if strVal, ok := val.(string); ok {
		return strVal != "" && strVal != "false" && strVal != "0", nil
	}
	
	return val != nil, nil
}

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

func (pc *PermissionChecker) getValue(expr string, ctx *Context) interface{} {
	expr = strings.TrimSpace(expr)
	
	if (strings.HasPrefix(expr, `"`) && strings.HasSuffix(expr, `"`)) ||
		(strings.HasPrefix(expr, `'`) && strings.HasSuffix(expr, `'`)) {
		return strings.Trim(expr, `"'`)
	}
	
	if expr == "true" {
		return true
	}
	if expr == "false" {
		return false
	}
	
	if num, err := strconv.ParseFloat(expr, 64); err == nil {
		return num
	}
	if num, err := strconv.ParseInt(expr, 10, 64); err == nil {
		return num
	}
	
	if ctx.Resource != nil {
		val := pc.getFieldValue(ctx.Resource, expr)
		if val != nil {
			return val
		}
	}
	
	if strings.HasPrefix(expr, "@") {
		return pc.getContextVariable(expr, ctx)
	}
	
	return expr
}

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
	
	field := val.FieldByName(fieldName)
	if field.IsValid() {
		return field.Interface()
	}
	
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

func (pc *PermissionChecker) getContextVariable(varName string, ctx *Context) interface{} {
	if varName == "@request.auth.id" {
		return pc.getAuthID(ctx)
	}
	
	if strings.HasPrefix(varName, "@request.user.") {
		fieldName := strings.TrimPrefix(varName, "@request.user.")
		return pc.getFieldValue(ctx.User, fieldName)
	}
	
	if strings.HasPrefix(varName, "@request.auth.") {
		fieldName := strings.TrimPrefix(varName, "@request.auth.")
		return pc.getFieldValue(ctx.User, fieldName)
	}
	
	return nil
}

func (pc *PermissionChecker) compareValues(left, right interface{}, op string) bool {
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

func (pc *PermissionChecker) normalizeValue(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	
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

func (pc *PermissionChecker) CheckFieldPermission(ctx *Context, fieldName string) bool {
	for _, excluded := range ctx.Model.Options.ExcludeFields {
		if excluded == fieldName {
			return false
		}
	}
	
	if ctx.Action == "update" || ctx.Action == "create" {
		for _, readOnly := range ctx.Model.Options.ReadOnlyFields {
			if readOnly == fieldName {
				return false
			}
		}
	}
	
	return true
}

func (pc *PermissionChecker) FilterFields(ctx *Context, fields []FieldMeta) []FieldMeta {
	filtered := make([]FieldMeta, 0, len(fields))
	
	for _, field := range fields {
		if pc.CheckFieldPermission(ctx, field.Name) {
			filtered = append(filtered, field)
		}
	}
	
	return filtered
}


