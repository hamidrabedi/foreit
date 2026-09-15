package api

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/forgego/forge/api/authentication"
	"github.com/forgego/forge/api/docs"
	apierrors "github.com/forgego/forge/api/errors"
	"github.com/forgego/forge/api/exceptions"
	"github.com/forgego/forge/api/permissions"
	"github.com/forgego/forge/api/throttling"
	"github.com/forgego/forge/orm"
	"github.com/forgego/forge/schema"
	forgehttp "github.com/forgego/forge/server"
)

// ViewSet is the base interface for all viewsets
type ViewSet interface {
	// List handles GET /resource/
	List(w http.ResponseWriter, r *http.Request)
	// Create handles POST /resource/
	Create(w http.ResponseWriter, r *http.Request)
	// Retrieve handles GET /resource/{id}/
	Retrieve(w http.ResponseWriter, r *http.Request)
	// Update handles PUT /resource/{id}/
	Update(w http.ResponseWriter, r *http.Request)
	// PartialUpdate handles PATCH /resource/{id}/
	PartialUpdate(w http.ResponseWriter, r *http.Request)
	// Destroy handles DELETE /resource/{id}/
	Destroy(w http.ResponseWriter, r *http.Request)
}

// QuerySetInterface defines the common methods needed from QuerySet for viewset operations.
// This interface allows type-safe operations without reflection when the QuerySet implements it.
type QuerySetInterface interface {
	Count(ctx context.Context) (int64, error)
	All(ctx context.Context) (interface{}, error)
	Filter(expr interface{}) interface{}
	OrderBy(fields ...interface{}) interface{}
	Limit(limit int) interface{}
	Offset(offset int) interface{}
}

// ManagerInterface defines the common methods needed from Manager for viewset operations.
type ManagerInterface interface {
	Get(ctx context.Context, id int64) (interface{}, error)
	Create(ctx context.Context, model interface{}) error
	Update(ctx context.Context, model interface{}) error
	Delete(ctx context.Context, model interface{}) error
}

// BaseViewSet provides common viewset functionality
type BaseViewSet struct {
	Serializer func() Serializer
	Queryset   interface{} // This would be a QuerySet in real implementation
	Model      interface{}
	// ExcludeResponseFields holds response keys removed from serialized output; nil keeps every field.
	ExcludeResponseFields []string
	// Authentication uses the current defaults when nil; a non-nil empty slice disables authentication.
	Authentication []authentication.Authentication
	// Permissions uses the current defaults when nil; a non-nil empty slice disables permission checks.
	Permissions []permissions.Permission
	// Throttles uses the current defaults when nil; a non-nil empty slice disables throttling.
	Throttles   []throttling.Throttle
	ErrorWriter func(http.ResponseWriter, *http.Request, error)

	actionMu sync.RWMutex
	action   string
}

type actionContextKeyType struct{}

var actionContextKey = actionContextKeyType{}

// ActionFromContext returns the action name stored in ctx, or an empty string if none is set.
func ActionFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if action, ok := ctx.Value(actionContextKey).(string); ok {
		return action
	}
	return ""
}

// GetActionFromRequest returns the action name stored in r's context, or an empty string if none is set.
func GetActionFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	return ActionFromContext(r.Context())
}

func withAction(r *http.Request, action string) *http.Request {
	if r == nil {
		return nil
	}
	return r.WithContext(context.WithValue(r.Context(), actionContextKey, action))
}

// NewBaseViewSet creates a new base viewset
func NewBaseViewSet(serializer func() Serializer, queryset, model interface{}) *BaseViewSet {
	return &BaseViewSet{
		Serializer: serializer,
		Queryset:   queryset,
		Model:      model,
	}
}

// getManager gets the manager for operations
func (vs *BaseViewSet) getManager() reflect.Value {
	// If Queryset is set and looks like a manager (has Create method), use it
	if vs.Queryset != nil {
		qsValue := reflect.ValueOf(vs.Queryset)
		qsType := qsValue.Type()

		// Use cached method lookup instead of MethodByName
		if _, ok := globalCache.GetMethod(qsType, "Create"); ok {
			return qsValue
		}
	}

	return reflect.Value{}
}

// GetAction returns the dispatched action. Inside permission and throttle checks
// the view is request-scoped, so it returns that request's action; on the shared
// viewset it returns the most recently dispatched action.
func (vs *BaseViewSet) GetAction() string {
	vs.actionMu.RLock()
	defer vs.actionMu.RUnlock()
	return vs.action
}

// SetAction sets the fallback action name on BaseViewSet for backward compatibility.
// During normal HTTP request processing, actions are carried in the request context.
func (vs *BaseViewSet) SetAction(action string) {
	vs.actionMu.Lock()
	defer vs.actionMu.Unlock()
	vs.action = action
}

func (vs *BaseViewSet) authenticateRequest(r *http.Request) error {
	authClasses := vs.Authentication
	if authClasses == nil {
		authClasses = GetDefaultAuthentication()
	}
	result, err := authentication.AuthenticateRequest(r, authClasses)
	if err != nil {
		return exceptions.NewAuthenticationFailed(err.Error())
	}
	if result == nil {
		return nil
	}
	authentication.SetUserOnRequest(r, result.User)
	authentication.SetAuthOnRequest(r, result.Auth)
	return nil
}

// viewForRequest returns a request-scoped view carrying the request's action.
// When the request carries no action it returns vs unchanged.
func (vs *BaseViewSet) viewForRequest(r *http.Request) *BaseViewSet {
	action := GetActionFromRequest(r)
	if action == "" {
		return vs
	}
	return &BaseViewSet{
		Serializer:            vs.Serializer,
		Queryset:              vs.Queryset,
		Model:                 vs.Model,
		ExcludeResponseFields: vs.ExcludeResponseFields,
		Authentication:        vs.Authentication,
		Permissions:           vs.Permissions,
		Throttles:             vs.Throttles,
		ErrorWriter:           vs.ErrorWriter,
		action:                action,
	}
}

func (vs *BaseViewSet) checkPermissions(r *http.Request) error {
	perms := vs.Permissions
	if perms == nil {
		perms = GetDefaultPermissions()
	}
	reqView := vs.viewForRequest(r)
	if permissions.CheckPermissions(r, reqView, perms) {
		return nil
	}
	for _, permission := range perms {
		if !permission.HasPermission(r, reqView) {
			return exceptions.NewPermissionDenied(permission.GetMessage())
		}
	}
	return exceptions.NewPermissionDenied("Permission denied")
}

func (vs *BaseViewSet) checkObjectPermissions(r *http.Request, object interface{}) error {
	perms := vs.Permissions
	if perms == nil {
		perms = GetDefaultPermissions()
	}
	reqView := vs.viewForRequest(r)
	if permissions.CheckObjectPermissions(r, reqView, object, perms) {
		return nil
	}
	for _, permission := range perms {
		if !permission.HasObjectPermission(r, reqView, object) {
			return exceptions.NewPermissionDenied(permission.GetMessage())
		}
	}
	return exceptions.NewPermissionDenied("Permission denied")
}

func (vs *BaseViewSet) checkThrottles(r *http.Request) error {
	throttles := vs.Throttles
	if throttles == nil {
		throttles = GetDefaultThrottles()
	}
	err := throttling.CheckThrottles(r, vs.viewForRequest(r), throttles)
	if err == nil {
		return nil
	}
	throttled, ok := err.(*throttling.ThrottledError)
	if !ok {
		return err
	}
	return exceptions.NewThrottled("Request was throttled", throttled.WaitDuration)
}

func (vs *BaseViewSet) checkRequest(w http.ResponseWriter, r *http.Request, action string) bool {
	vs.SetAction(action)
	if GetActionFromRequest(r) == "" {
		r = withAction(r, action)
	}
	checks := []func(*http.Request) error{
		vs.authenticateRequest,
		vs.checkPermissions,
		vs.checkThrottles,
	}
	for _, check := range checks {
		if err := check(r); err != nil {
			vs.handleException(w, r, err)
			return false
		}
	}
	return true
}

func (vs *BaseViewSet) handleException(w http.ResponseWriter, r *http.Request, err error) {
	if vs.ErrorWriter != nil {
		vs.ErrorWriter(w, r, err)
		return
	}
	apierrors.WriteError(w, r, err)
}

func (vs *BaseViewSet) allowObject(w http.ResponseWriter, r *http.Request, object interface{}) bool {
	if err := vs.checkObjectPermissions(r, object); err != nil {
		vs.handleException(w, r, err)
		return false
	}
	return true
}

// List handles GET /resource/
func (vs *BaseViewSet) List(w http.ResponseWriter, r *http.Request) {
	r = withAction(r, "list")
	if !vs.checkRequest(w, r, "list") {
		return
	}
	ctx := r.Context()

	// Get pagination parameters, defaulting to the configured page size.
	defaultPageSize := 20
	if s := GetSettings(); s != nil && s.PageSize > 0 {
		defaultPageSize = s.PageSize
	}
	page, pageSize, _ := ParsePaginationParams(r, defaultPageSize)

	// Get queryset using reflection. A bare manager that cannot paginate is
	// converted to a queryset via its QuerySet constructor when available, so
	// Offset/Limit are applied through the queryset instead of loading all rows.
	querysetValue := reflect.ValueOf(vs.Queryset)
	if !querysetValue.IsValid() {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Queryset not set")
		return
	}
	querysetValue = paginatableQueryset(querysetValue)

	// Apply filtering from query params
	qs := applyFilters(querysetValue, r)

	// Apply ordering
	qs = applyOrdering(qs, r)

	// Get total count using cached method lookup
	qsType := qs.Type()
	countMethod, ok := globalCache.GetMethod(qsType, "Count")
	if !ok {
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Count method not found")
		return
	}

	var totalCount int64
	results := countMethod.Func.Call([]reflect.Value{qs, reflect.ValueOf(ctx)})
	if len(results) >= 2 {
		if !results[1].IsNil() {
			if err, ok := results[1].Interface().(error); ok {
				_ = forgehttp.SendError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if results[0].CanInterface() {
			if count, ok := results[0].Interface().(int64); ok {
				totalCount = count
			}
		}
	}

	// Apply pagination using cached method lookup
	offset := (page - 1) * pageSize

	offsetMethod, ok := globalCache.GetMethod(qsType, "Offset")
	if ok {
		offsetResults := offsetMethod.Func.Call([]reflect.Value{qs, reflect.ValueOf(offset)})
		if len(offsetResults) > 0 {
			if newQS := offsetResults[0].Interface(); newQS != nil {
				qs = reflect.ValueOf(newQS)
				qsType = qs.Type()
			}
		}
	}

	limitMethod, ok := globalCache.GetMethod(qsType, "Limit")
	if ok {
		limitResults := limitMethod.Func.Call([]reflect.Value{qs, reflect.ValueOf(pageSize)})
		if len(limitResults) > 0 {
			if newQS := limitResults[0].Interface(); newQS != nil {
				qs = reflect.ValueOf(newQS)
				qsType = qs.Type()
			}
		}
	}

	// Execute query using cached method lookup
	allMethod, ok := globalCache.GetMethod(qsType, "All")
	if !ok {
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "All method not found")
		return
	}

	var resultList []interface{}
	allResults := allMethod.Func.Call([]reflect.Value{qs, reflect.ValueOf(ctx)})
	if len(allResults) >= 2 {
		// Check for error first
		if errVal := allResults[1]; !errVal.IsNil() {
			if err, ok := errVal.Interface().(error); ok {
				_ = forgehttp.SendError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		if allResults[0].CanInterface() {
			if objects, ok := allResults[0].Interface().([]interface{}); ok {
				resultList = objects
			} else if allResults[0].Kind() == reflect.Slice {
				// Try to convert slice of pointers
				sliceValue := allResults[0]
				for i := 0; i < sliceValue.Len(); i++ {
					resultList = append(resultList, sliceValue.Index(i).Interface())
				}
			}
		}
	}

	// Serialize results
	serialized := vs.stripExcludedFieldsMany(SerializeMany(resultList))

	// Send paginated response
	// nolint:errcheck // HTTP response errors can't be handled meaningfully
	_ = SendPaginatedResponse(w, r, serialized, int(totalCount), page, pageSize)
}

// Create handles POST /resource/
func (vs *BaseViewSet) Create(w http.ResponseWriter, r *http.Request) {
	r = withAction(r, "create")
	if !vs.checkRequest(w, r, "create") {
		return
	}
	ctx := r.Context()

	var data map[string]interface{}
	if err := forgehttp.GetJSON(r, &data); err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	serializer := vs.Serializer()
	serializer.SetData(data)

	if err := serializer.Validate(); err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": serializer.Errors(),
		})
		return
	}

	// Create model instance from serializer data
	modelValue := reflect.New(reflect.TypeOf(vs.Model).Elem())
	instance := modelValue.Interface()

	// Populate instance from data
	populateFromMap(instance, data)

	// Get manager and call Create
	manager := vs.getManager()
	if !manager.IsValid() {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Manager not found")
		return
	}

	managerType := manager.Type()
	createMethod, ok := globalCache.GetMethod(managerType, "Create")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Create method not found")
		return
	}

	results := createMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(results) > 0 && !results[0].IsNil() {
		if err, ok := results[0].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	// Serialize and return created instance
	serialized := vs.stripExcludedFields(SerializeModel(instance))
	// nolint:errcheck // HTTP response errors can't be handled meaningfully
	_ = forgehttp.SendJSON(w, http.StatusCreated, serialized)
}

// Retrieve handles GET /resource/{id}/
func (vs *BaseViewSet) Retrieve(w http.ResponseWriter, r *http.Request) {
	r = withAction(r, "retrieve")
	if !vs.checkRequest(w, r, "retrieve") {
		return
	}
	ctx := r.Context()

	// Get ID from URL
	idStr := forgehttp.GetParam(r, "id")
	if idStr == "" {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "ID is required")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Get manager and call Get
	manager := vs.getManager()
	if !manager.IsValid() {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Manager not found")
		return
	}

	managerType := manager.Type()
	getMethod, ok := globalCache.GetMethod(managerType, "Get")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Get method not found")
		return
	}

	results := getMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(results) < 2 {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Invalid Get method")
		return
	}

	if !results[1].IsNil() {
		if err, ok := results[1].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	instance := results[0].Interface()
	if !vs.allowObject(w, r, instance) {
		return
	}
	serialized := vs.stripExcludedFields(SerializeModel(instance))
	// nolint:errcheck // HTTP response errors can't be handled meaningfully
	_ = forgehttp.SendJSON(w, http.StatusOK, serialized)
}

// Update handles PUT /resource/{id}/
func (vs *BaseViewSet) Update(w http.ResponseWriter, r *http.Request) {
	vs.update(w, r, "update")
}

func (vs *BaseViewSet) update(w http.ResponseWriter, r *http.Request, action string) {
	r = withAction(r, action)
	if !vs.checkRequest(w, r, action) {
		return
	}
	ctx := r.Context()

	idStr := forgehttp.GetParam(r, "id")
	if idStr == "" {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "ID is required")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	var data map[string]interface{}
	if err := forgehttp.GetJSON(r, &data); err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	serializer := vs.Serializer()
	serializer.SetData(data)

	if err := serializer.Validate(); err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendJSON(w, http.StatusBadRequest, map[string]interface{}{
			"errors": serializer.Errors(),
		})
		return
	}

	// Get existing instance
	manager := vs.getManager()
	if !manager.IsValid() {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Manager not found")
		return
	}

	managerType := manager.Type()
	getMethod, ok := globalCache.GetMethod(managerType, "Get")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Get method not found")
		return
	}

	getResults := getMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(getResults) < 2 || !getResults[1].IsNil() {
		if err, ok := getResults[1].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	instance := getResults[0].Interface()
	if !vs.allowObject(w, r, instance) {
		return
	}

	// Get schema primary key field info
	pkGoName, pkDBName, pkJSONName := vs.getSchemaPrimaryKeyField()

	// Capture primary key from existing instance before populating
	origPK := getPrimaryKeyValue(instance, pkGoName, pkDBName, pkJSONName)

	// Collect primary-key fields and read-only fields to ignore from body
	ignoredKeys := []string{"id", "ID", "Id"}
	if pkJSONName != "" {
		ignoredKeys = append(ignoredKeys, pkJSONName)
	}
	if pkGoName != "" {
		ignoredKeys = append(ignoredKeys, pkGoName)
	}
	if pkDBName != "" && pkDBName != pkGoName && pkDBName != pkJSONName {
		ignoredKeys = append(ignoredKeys, pkDBName)
	}
	if ro, ok := serializer.(interface{ ReadOnlyFields() []string }); ok {
		ignoredKeys = append(ignoredKeys, ro.ReadOnlyFields()...)
	}

	// Populate from data, ignoring primary-key fields
	populateFromMap(instance, data, ignoredKeys...)

	// Restore primary key from the object loaded via the URL after populating
	restorePrimaryKey(instance, origPK, id, pkGoName, pkDBName, pkJSONName)

	// Update
	updateMethod, ok := globalCache.GetMethod(managerType, "Update")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Update method not found")
		return
	}

	updateResults := updateMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(updateResults) > 0 && !updateResults[0].IsNil() {
		if err, ok := updateResults[0].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	serialized := vs.stripExcludedFields(SerializeModel(instance))
	forgehttp.SendJSON(w, http.StatusOK, serialized)
}

// PartialUpdate handles PATCH /resource/{id}/
func (vs *BaseViewSet) PartialUpdate(w http.ResponseWriter, r *http.Request) {
	vs.update(w, r, "partial_update")
}

// Destroy handles DELETE /resource/{id}/
func (vs *BaseViewSet) Destroy(w http.ResponseWriter, r *http.Request) {
	r = withAction(r, "destroy")
	if !vs.checkRequest(w, r, "destroy") {
		return
	}
	ctx := r.Context()

	idStr := forgehttp.GetParam(r, "id")
	if idStr == "" {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "ID is required")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	// Get instance first
	manager := vs.getManager()
	if !manager.IsValid() {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Manager not found")
		return
	}

	managerType := manager.Type()
	getMethod, ok := globalCache.GetMethod(managerType, "Get")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Get method not found")
		return
	}

	getResults := getMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(id),
	})

	if len(getResults) < 2 || !getResults[1].IsNil() {
		if err, ok := getResults[1].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusNotFound, err.Error())
			return
		}
	}

	instance := getResults[0].Interface()
	if !vs.allowObject(w, r, instance) {
		return
	}

	// Delete
	deleteMethod, ok := globalCache.GetMethod(managerType, "Delete")
	if !ok {
		// nolint:errcheck // HTTP response errors can't be handled meaningfully
		_ = forgehttp.SendError(w, http.StatusInternalServerError, "Delete method not found")
		return
	}

	deleteResults := deleteMethod.Func.Call([]reflect.Value{
		manager,
		reflect.ValueOf(ctx),
		reflect.ValueOf(instance),
	})

	if len(deleteResults) > 0 && !deleteResults[0].IsNil() {
		if err, ok := deleteResults[0].Interface().(error); ok && err != nil {
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewSetHandler creates an HTTP handler from a viewset
func ViewSetHandler(vs ViewSet, actions []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Check if it's a list or retrieve
			if forgehttp.GetParam(r, "id") != "" {
				vs.Retrieve(w, r)
			} else {
				vs.List(w, r)
			}
		case http.MethodPost:
			vs.Create(w, r)
		case http.MethodPut:
			vs.Update(w, r)
		case http.MethodPatch:
			vs.PartialUpdate(w, r)
		case http.MethodDelete:
			vs.Destroy(w, r)
		default:
			// nolint:errcheck // HTTP response errors can't be handled meaningfully
			_ = forgehttp.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	}
}

type customRoute struct {
	method  string
	path    string
	handler http.HandlerFunc
}

// ActionConfig describes an extra endpoint on a registered resource.
type ActionConfig struct {
	Methods []string // e.g. http.MethodPost; defaults to GET
	Detail  bool     // true: /{resource}/{id}/{path}; false: /{resource}/{path}
	URLPath string   // defaults to the action name
}

type actionRoute struct {
	resource string
	name     string
	cfg      ActionConfig
	handler  http.HandlerFunc
}

func isValidHTTPMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost,
		http.MethodPut, http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

// Router is a router for viewsets (DRF-like)
type Router struct {
	prefix       string
	routes       map[string]ViewSet
	customRoutes []customRoute
	actions      []actionRoute
}

// NewRouter creates a new API router
func NewRouter(prefix string) *Router {
	return &Router{
		prefix:       prefix,
		routes:       make(map[string]ViewSet),
		customRoutes: make([]customRoute, 0),
		actions:      make([]actionRoute, 0),
	}
}

// Register registers a viewset with a resource name
func (r *Router) Register(resource string, vs ViewSet) {
	r.routes[resource] = vs
}

// Action registers an extra endpoint on a resource.
func (r *Router) Action(resource, name string, cfg ActionConfig, handler http.HandlerFunc) {
	if len(cfg.Methods) == 0 {
		cfg.Methods = []string{http.MethodGet}
	} else {
		for _, method := range cfg.Methods {
			if !isValidHTTPMethod(method) {
				panic(fmt.Sprintf("unknown HTTP method: %s", method))
			}
		}
	}
	if cfg.URLPath == "" {
		cfg.URLPath = name
	}
	r.actions = append(r.actions, actionRoute{
		resource: resource,
		name:     name,
		cfg:      cfg,
		handler:  handler,
	})
}

// Get registers a custom GET route on the API router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.customRoutes = append(r.customRoutes, customRoute{
		method:  http.MethodGet,
		path:    path,
		handler: handler,
	})
}

// Post registers a custom POST route on the API router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.customRoutes = append(r.customRoutes, customRoute{
		method:  http.MethodPost,
		path:    path,
		handler: handler,
	})
}

// Handle registers a custom route with a specific method on the API router
func (r *Router) Handle(method, path string, handler http.HandlerFunc) {
	r.customRoutes = append(r.customRoutes, customRoute{
		method:  method,
		path:    path,
		handler: handler,
	})
}

func (r *Router) registerResourceRoutes(router *forgehttp.Router, resource string, vs ViewSet) {
	path := r.prefix + "/" + resource
	router.Route(path, func(sub *forgehttp.Router) {
		sub.Get("/", vs.List)
		sub.Post("/", vs.Create)
		sub.Get("/{id}", vs.Retrieve)
		sub.Put("/{id}", vs.Update)
		sub.Patch("/{id}", vs.PartialUpdate)
		sub.Delete("/{id}", vs.Destroy)
		sub.Options("/", func(w http.ResponseWriter, r *http.Request) {
			docs.OptionsHandler(w, r, vs, nil)
		})
	})
}

func (r *Router) registerActions(router *forgehttp.Router) {
	cleanPrefix := strings.TrimSuffix(r.prefix, "/")
	for _, act := range r.actions {
		cleanResource := strings.Trim(act.resource, "/")
		cleanActionPath := strings.Trim(act.cfg.URLPath, "/")
		var fullPath string
		if act.cfg.Detail {
			fullPath = cleanPrefix + "/" + cleanResource + "/{id}/" + cleanActionPath
		} else {
			fullPath = cleanPrefix + "/" + cleanResource + "/" + cleanActionPath
		}
		for _, method := range act.cfg.Methods {
			router.Method(method, fullPath, act.handler)
		}
	}
}

// RegisterRoutes registers all routes on a chi router
func (r *Router) RegisterRoutes(router *forgehttp.Router) {
	for resource, vs := range r.routes {
		r.registerResourceRoutes(router, resource, vs)
	}
	r.registerActions(router)
	for _, cr := range r.customRoutes {
		cleanPath := "/" + strings.TrimPrefix(cr.path, "/")
		fullPath := r.prefix + cleanPath
		router.Method(cr.method, fullPath, cr.handler)
	}
}

// Helper functions for viewset operations

// paginatableQueryset returns a queryset capable of Offset/Limit. If qs already
// supports Offset and Limit it is returned unchanged. Otherwise, when it exposes
// a zero-argument QuerySet constructor (like orm managers and the generated
// API managers), that constructor is invoked and its result is used, so
// pagination is applied through the queryset instead of loading all rows.
func paginatableQueryset(qs reflect.Value) reflect.Value {
	if !qs.IsValid() {
		return qs
	}
	qsType := qs.Type()
	if _, ok := globalCache.GetMethod(qsType, "Offset"); !ok {
		return maybeConvertToQueryset(qs)
	}
	if _, ok := globalCache.GetMethod(qsType, "Limit"); !ok {
		return maybeConvertToQueryset(qs)
	}
	return qs
}

func maybeConvertToQueryset(qs reflect.Value) reflect.Value {
	qsType := qs.Type()
	qsMethod, ok := globalCache.GetMethod(qsType, "QuerySet")
	if !ok {
		return qs
	}
	if !qsMethod.Func.IsValid() || qsMethod.Func.Type().NumIn() != 1 {
		return qs
	}
	results := qsMethod.Func.Call([]reflect.Value{qs})
	if len(results) == 0 {
		return qs
	}
	if len(results) >= 2 {
		if errVal := results[1]; errVal.IsValid() && errVal.CanInterface() {
			if err, ok := errVal.Interface().(error); ok && err != nil {
				return qs
			}
		}
	}
	first := results[0]
	if !first.IsValid() || (first.Kind() == reflect.Ptr && first.IsNil()) {
		return qs
	}
	if !first.CanInterface() || first.Interface() == nil {
		return qs
	}
	if first.Kind() == reflect.Interface && !first.IsNil() {
		first = first.Elem()
	}
	return first
}

// applyFilters applies query parameter filters to queryset
func applyFilters(qs reflect.Value, r *http.Request) reflect.Value {
	// Get filter parameters from query string
	query := r.URL.Query()
	qsType := qs.Type()

	// Get cached Filter method
	filterMethod, hasFilter := globalCache.GetMethod(qsType, "Filter")

	for key, values := range query {
		if len(values) == 0 || key == "page" || key == "page_size" || key == "ordering" || key == "search" {
			continue
		}

		if hasFilter {
			value := values[0]

			expr := buildFilterExpr(key, value)
			if expr != nil {
				results := filterMethod.Func.Call([]reflect.Value{qs, reflect.ValueOf(expr)})
				if len(results) > 0 {
					if newQS := results[0].Interface(); newQS != nil {
						qs = reflect.ValueOf(newQS)
						qsType = qs.Type()
						// Update cached method for new queryset type
						filterMethod, hasFilter = globalCache.GetMethod(qsType, "Filter")
					}
				}
			}
		}
	}

	return qs
}

func buildFilterExpr(rawKey string, rawValue string) orm.Expression {
	field, lookup := parseLookup(rawKey)
	parsed := parseFilterValue(rawValue)

	f := orm.F(field)
	switch lookup {
	case "exact":
		return f.Eq(parsed)
	case "ne":
		return f.Ne(parsed)
	case "gt":
		return f.Gt(parsed)
	case "gte":
		return f.Gte(parsed)
	case "lt":
		return f.Lt(parsed)
	case "lte":
		return f.Lte(parsed)
	case "contains":
		return f.Contains(rawValue)
	case "icontains":
		return f.IContains(rawValue)
	case "startswith":
		return f.StartsWith(rawValue)
	case "endswith":
		return f.EndsWith(rawValue)
	case "in":
		parts := strings.Split(rawValue, ",")
		args := make([]interface{}, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			args = append(args, parseFilterValue(part))
		}
		if len(args) == 0 {
			return f.Eq(parsed)
		}
		return f.In(args...)
	case "isnull":
		lower := strings.ToLower(rawValue)
		if lower == "true" || lower == "1" {
			return f.IsNull()
		}
		return f.IsNotNull()
	default:
		return f.Eq(parsed)
	}
}

func parseLookup(key string) (string, string) {
	if strings.Contains(key, "__") {
		parts := strings.SplitN(key, "__", 2)
		return parts[0], parts[1]
	}
	return key, "exact"
}

func parseFilterValue(raw string) interface{} {
	if raw == "" {
		return raw
	}
	if boolVal, err := strconv.ParseBool(raw); err == nil {
		return boolVal
	}
	if intVal, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return intVal
	}
	if floatVal, err := strconv.ParseFloat(raw, 64); err == nil {
		return floatVal
	}
	return raw
}

// applyOrdering applies ordering from query parameters
func applyOrdering(qs reflect.Value, r *http.Request) reflect.Value {
	ordering := r.URL.Query().Get("ordering")
	if ordering == "" {
		return qs
	}

	qsType := qs.Type()
	orderByMethod, ok := globalCache.GetMethod(qsType, "OrderBy")
	if !ok {
		return qs
	}

	rawFields := strings.Split(ordering, ",")
	args := make([]reflect.Value, 0, len(rawFields)+1)
	args = append(args, qs) // Add receiver
	for _, field := range rawFields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		args = append(args, reflect.ValueOf(field))
	}
	if len(args) <= 1 { // Only receiver, no fields
		return qs
	}
	results := orderByMethod.Func.Call(args)
	if len(results) > 0 {
		if newQS := results[0].Interface(); newQS != nil {
			return reflect.ValueOf(newQS)
		}
	}

	return qs
}

// populateFromMap populates a model instance from a map, optionally ignoring specified keys
func populateFromMap(instance interface{}, data map[string]interface{}, ignoredKeys ...string) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	if instanceValue.Kind() != reflect.Struct {
		return
	}

	ignored := make(map[string]bool, len(ignoredKeys))
	for _, k := range ignoredKeys {
		ignored[strings.ToLower(k)] = true
	}

	var applyToStruct func(target reflect.Value)
	applyToStruct = func(target reflect.Value) {
		targetType := target.Type()
		for i := 0; i < targetType.NumField(); i++ {
			field := targetType.Field(i)
			if !field.IsExported() {
				continue
			}
			fieldValue := target.Field(i)

			if field.Anonymous {
				switch fieldValue.Kind() {
				case reflect.Struct:
					applyToStruct(fieldValue)
				case reflect.Ptr:
					if fieldValue.IsNil() {
						continue
					}
					if fieldValue.Elem().Kind() == reflect.Struct {
						applyToStruct(fieldValue.Elem())
					}
				}
				continue
			}

			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			tagParts := strings.Split(jsonTag, ",")
			key := tagParts[0]
			if key == "" {
				continue
			}
			if ignored[strings.ToLower(key)] || ignored[strings.ToLower(field.Name)] {
				continue
			}
			if value, ok := data[key]; ok {
				if fieldValue.CanSet() {
					setFieldValue(fieldValue, value)
				}
			}
		}
	}

	applyToStruct(instanceValue)
}

// getSchemaPrimaryKeyField returns the primary key field from the model's schema if available.
// Returns the Go struct field name, db column name, and json tag name.
// When no schema is available it returns empty strings and callers fall back
// to the "id"/"ID"/"Id" variants.
func (vs *BaseViewSet) getSchemaPrimaryKeyField() (goName, dbName, jsonName string) {
	if vs.Model == nil {
		return "", "", ""
	}
	// Check if model implements schema.Schema
	modelSchema, ok := vs.Model.(schema.Schema)
	if !ok {
		return "", "", ""
	}
	var pkName, pkDB string
	for _, f := range modelSchema.Fields() {
		if f.PrimaryKey {
			pkName = f.Name
			pkDB = f.DBColumn
			break
		}
	}
	if pkName == "" && pkDB == "" {
		return "", "", ""
	}
	if pkDB == "" {
		pkDB = pkName
	}
	goResolved, jsonResolved := resolveModelFieldNames(vs.Model, pkName, pkDB)
	if goResolved == "" {
		// No matching Go struct field found; still return db/json names so the
		// body value for the primary key is ignored.
		if jsonResolved == "" {
			jsonResolved = pkName
		}
		return "", pkDB, jsonResolved
	}
	if jsonResolved == "" {
		jsonResolved = pkName
	}
	return goResolved, pkDB, jsonResolved
}

// resolveModelFieldNames finds the Go struct field name and json tag name matching
// any of the given schema/db names (case-insensitive on Go name, exact on tags).
func resolveModelFieldNames(model interface{}, names ...string) (goName, jsonName string) {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			t := v.Type().Elem()
			if t.Kind() == reflect.Struct {
				return resolveFieldNamesInType(t, names...)
			}
			return "", ""
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", ""
	}
	return resolveFieldNamesInType(v.Type(), names...)
}

func resolveFieldNamesInType(t reflect.Type, names ...string) (goName, jsonName string) {
	lowered := make(map[string]bool, len(names))
	for _, n := range names {
		if n != "" {
			lowered[strings.ToLower(n)] = true
		}
	}
	var search func(rt reflect.Type) (string, string, bool)
	search = func(rt reflect.Type) (string, string, bool) {
		for i := 0; i < rt.NumField(); i++ {
			field := rt.Field(i)
			if !field.IsExported() {
				continue
			}
			if field.Anonymous {
				ft := field.Type
				if ft.Kind() == reflect.Ptr {
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct {
					if gn, jn, ok := search(ft); ok {
						return gn, jn, true
					}
				}
				continue
			}
			jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
			dbTag := strings.Split(field.Tag.Get("db"), ",")[0]
			matched := lowered[strings.ToLower(field.Name)]
			for _, n := range names {
				if n == "" {
					continue
				}
				if jsonTag != "" && jsonTag != "-" && n == jsonTag {
					matched = true
				}
				if dbTag != "" && n == dbTag {
					matched = true
				}
			}
			if matched {
				return field.Name, jsonTag, true
			}
		}
		return "", "", false
	}
	gn, jn, ok := search(t)
	if !ok {
		return "", ""
	}
	return gn, jn
}

// getPrimaryKeyValue reads the current primary key value from instance, preferring
// the schema primary key field names and falling back to "ID"/"Id"/"id".
func getPrimaryKeyValue(instance interface{}, pkGoName, pkDBName, pkJSONName string) interface{} {
	v := reflect.ValueOf(instance)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil
	}
	if fv, ok := findPrimaryKeyField(v, []string{pkGoName, pkDBName, pkJSONName}); ok {
		return fv.Interface()
	}
	if pkGoName == "" && pkDBName == "" && pkJSONName == "" {
		if fv, ok := findPrimaryKeyField(v, []string{"ID", "Id", "id"}); ok {
			return fv.Interface()
		}
	}
	return nil
}

// restorePrimaryKey restores the primary key value on instance from the original value or urlID.
// Schema primary key field names take precedence; "ID"/"Id"/"id" are only used
// as a fallback when no schema primary key is known. A body value for the
// primary key never changes which row is written.
func restorePrimaryKey(instance interface{}, origPK interface{}, urlID int64, pkGoName, pkDBName, pkJSONName string) {
	if instance == nil {
		return
	}
	hasSchemaPK := pkGoName != "" || pkDBName != "" || pkJSONName != ""
	if !hasSchemaPK {
		if m, ok := instance.(interface{ SetID(int64) }); ok {
			if origInt, ok := origPK.(int64); ok && origInt != 0 {
				m.SetID(origInt)
			} else {
				m.SetID(urlID)
			}
			return
		}
	}
	v := reflect.ValueOf(instance)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	setField := func(f reflect.Value) bool {
		if !f.CanSet() {
			return false
		}
		if origPK != nil {
			origVal := reflect.ValueOf(origPK)
			if origVal.IsValid() && origVal.Type().AssignableTo(f.Type()) {
				f.Set(origVal)
				return true
			}
		}
		if isIntKind(f.Kind()) {
			if f.Kind() >= reflect.Uint && f.Kind() <= reflect.Uint64 {
				f.SetUint(uint64(urlID))
			} else {
				f.SetInt(urlID)
			}
			return true
		}
		return false
	}

	// Prefer the schema primary key when known.
	if pkGoName != "" || pkDBName != "" || pkJSONName != "" {
		if f, ok := findPrimaryKeyField(v, []string{pkGoName, pkDBName, pkJSONName}); ok {
			if setField(f) {
				return
			}
		}
	}

	for _, name := range []string{"ID", "Id", "id"} {
		f := v.FieldByName(name)
		if f.IsValid() {
			if setField(f) {
				return
			}
		}
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		parts := strings.Split(tag, ",")
		if len(parts) > 0 && strings.EqualFold(parts[0], "id") {
			f := v.Field(i)
			if f.IsValid() {
				if setField(f) {
					return
				}
			}
		}
	}
}

// findPrimaryKeyField locates the settable struct field matching any of the given
// names (Go field name case-insensitively, or exact json/db tag match),
// searching into embedded structs.
func findPrimaryKeyField(v reflect.Value, names []string) (reflect.Value, bool) {
	clean := make([]string, 0, len(names))
	lowered := make(map[string]bool)
	for _, n := range names {
		if n == "" {
			continue
		}
		clean = append(clean, n)
		lowered[strings.ToLower(n)] = true
	}
	if len(clean) == 0 {
		return reflect.Value{}, false
	}
	// Fast path: direct Go field name lookup.
	for _, n := range clean {
		if f := v.FieldByName(n); f.IsValid() {
			return f, true
		}
	}
	t := v.Type()
	var search func(rv reflect.Value, rt reflect.Type) (reflect.Value, bool)
	search = func(rv reflect.Value, rt reflect.Type) (reflect.Value, bool) {
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)
			if !sf.IsExported() {
				continue
			}
			fv := rv.Field(i)
			if sf.Anonymous {
				ft := sf.Type
				fvv := fv
				if ft.Kind() == reflect.Ptr {
					if fvv.IsNil() {
						continue
					}
					fvv = fvv.Elem()
					ft = ft.Elem()
				}
				if ft.Kind() == reflect.Struct {
					if res, ok := search(fvv, ft); ok {
						return res, true
					}
				}
				continue
			}
			jsonTag := strings.Split(sf.Tag.Get("json"), ",")[0]
			dbTag := strings.Split(sf.Tag.Get("db"), ",")[0]
			if lowered[strings.ToLower(sf.Name)] {
				return fv, true
			}
			for _, n := range clean {
				if jsonTag != "" && jsonTag != "-" && n == jsonTag {
					return fv, true
				}
				if dbTag != "" && n == dbTag {
					return fv, true
				}
			}
		}
		return reflect.Value{}, false
	}
	return search(v, t)
}

func isIntKind(k reflect.Kind) bool {
	return (k >= reflect.Int && k <= reflect.Int64) || (k >= reflect.Uint && k <= reflect.Uint64)
}

// setFieldValue sets a field value from interface{}
func setFieldValue(field reflect.Value, value interface{}) {
	if !field.CanSet() {
		return
	}

	valueValue := reflect.ValueOf(value)
	if !valueValue.IsValid() {
		switch field.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
			field.Set(reflect.Zero(field.Type()))
		}
		return
	}

	switch field.Kind() {
	case reflect.String:
		if valueValue.Kind() == reflect.String {
			field.SetString(valueValue.String())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valueValue.Kind() == reflect.Float64 {
			field.SetInt(int64(valueValue.Float()))
		} else if valueValue.Kind() == reflect.Int || valueValue.Kind() == reflect.Int64 {
			field.SetInt(valueValue.Int())
		}
	case reflect.Bool:
		if valueValue.Kind() == reflect.Bool {
			field.SetBool(valueValue.Bool())
		}
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			if str, ok := value.(string); ok {
				if str == "" {
					field.Set(reflect.ValueOf(time.Time{}))
					return
				}
				if parsed, err := time.Parse(time.RFC3339, str); err == nil {
					field.Set(reflect.ValueOf(parsed))
					return
				}
				if parsed, err := time.Parse("2006-01-02", str); err == nil {
					field.Set(reflect.ValueOf(parsed))
					return
				}
			}
		}
	default:
		if valueValue.Type().AssignableTo(field.Type()) {
			field.Set(valueValue)
		}
	}
}
