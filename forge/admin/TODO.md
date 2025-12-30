# Forge Admin - TODO List

This document tracks missing features, incomplete implementations, and improvements needed for the admin system.

## High Priority

### 1. Inline Editing (Inlines)
- [ ] **Complete inline instance retrieval** - `inlines/inline.go:33` has placeholder comment
  - Need to implement proper parent field filtering using ORM field expressions
  - Currently returns all instances instead of filtering by parent
- [ ] **Complete inline save logic** - `inlines/inline.go:54` has placeholder comment
  - Need to implement proper parent field setting using ORM field accessor
  - Currently doesn't set parent field on inline instances
- [ ] **Integrate inlines with form view** - Inlines need to be rendered in form templates
- [ ] **Tabular inline rendering** - UI component for tabular inline display
- [ ] **Stacked inline rendering** - UI component for stacked inline display

### 2. Form Validation
- [ ] **Complete form validation** - `validation.go:15-20` has simplified validation
  - Need to integrate with form view properly
  - Should validate against admin config fieldsets
  - Should use schema validation rules
- [ ] **Field-level validation** - `validation.go:56-68` has placeholder validation
  - Need to implement proper type checking and schema validation
  - Should respect field validators from schema

### 3. Field Expressions
- [ ] **Auto-create field expressions** - `admin.go:126` has placeholder comment
  - Need to create proper FieldExpr instances from schema fields
  - Currently just has placeholder comment

### 4. Authentication & Authorization
- [ ] **User context integration** - Need to verify user is properly passed through all handlers
- [ ] **Permission checking integration** - Verify permission hooks are called in all views
- [ ] **Default permission checker** - Ensure default implementation works correctly
- [ ] **Session-based authentication** - Verify session manager integration

### 5. Template Rendering
- [ ] **Complete template set** - Verify all templates exist and are functional
  - List view template
  - Form view template (add/edit)
  - Detail view template
  - Delete confirmation template
  - History template
- [ ] **Template error handling** - Proper error pages and messages
- [ ] **Template customization** - Allow custom templates per admin instance

## Medium Priority

### 6. HTTP Handlers
- [x] **HandleIndex implementation** - Admin index/dashboard page ✅
- [ ] **Error handling** - Consistent error responses (JSON/HTML)
- [ ] **Request validation** - Validate all incoming requests
- [ ] **CSRF protection** - Re-integrate CSRF protection (was removed)
- [ ] **Rate limiting** - Add rate limiting for admin endpoints

### 7. List View Features
- [ ] **Date hierarchy** - Implement date hierarchy navigation
- [ ] **Column sorting** - Verify sorting works for all field types
- [ ] **Bulk actions UI** - Frontend for bulk action selection
- [ ] **Export functionality** - Complete export implementation (CSV/JSON)
- [ ] **Advanced filtering UI** - Filter sidebar with all filter types

### 8. Form View Features
- [ ] **Prepopulated fields** - JavaScript for auto-population
- [ ] **Raw ID fields** - Lookup popup for foreign keys
- [ ] **Autocomplete fields** - Searchable select with AJAX
- [ ] **Radio fields** - Radio button widget with layout support
- [ ] **File upload handling** - File/image upload widgets
- [ ] **Rich text editor** - WYSIWYG editor integration

### 9. Advanced Features
- [ ] **Change history** - Complete history logging and viewing
- [ ] **Bulk actions** - Complete bulk action execution
- [ ] **Custom actions** - Verify custom action hooks work
- [ ] **View hooks** - Test all view hooks (ChangelistViewHook, ChangeViewHook, etc.)
- [ ] **Response hooks** - Test response hooks (ResponseAddHook, etc.)

### 10. Widgets
- [ ] **Widget registry** - Complete widget selection logic
- [ ] **Custom widgets** - Allow registering custom widgets
- [ ] **Widget validation** - Widget-level validation
- [ ] **Widget help text** - Display help text for fields

## Low Priority / Nice to Have

### 11. Testing
- [ ] **Unit tests** - Test all core functionality
- [ ] **Integration tests** - Test HTTP handlers
- [ ] **E2E tests** - Test full admin workflows
- [ ] **Test fixtures** - Complete test data fixtures

### 12. Documentation
- [ ] **API documentation** - Complete API reference
- [ ] **Usage examples** - More comprehensive examples
- [ ] **Migration guide** - Guide for migrating from old admin
- [ ] **Best practices** - Admin configuration best practices

### 13. Performance
- [ ] **Query optimization** - Optimize queryset operations
- [ ] **Caching** - Add caching for frequently accessed data
- [ ] **Lazy loading** - Lazy load related objects
- [ ] **Pagination optimization** - Efficient pagination for large datasets

### 14. Developer Experience
- [ ] **Code generation** - Complete codegen for admin setup
- [ ] **CLI commands** - Admin-specific CLI commands
- [ ] **Debug mode** - Admin debug mode with detailed errors
- [ ] **Development server** - Hot reload for templates

### 15. UI/UX Improvements
- [ ] **Responsive design** - Mobile-friendly admin interface
- [ ] **Dark mode** - Dark theme support
- [ ] **Accessibility** - ARIA labels and keyboard navigation
- [ ] **Internationalization** - i18n support for admin interface

## Known Issues

### Compilation Errors
- [ ] **ORM package errors** - Fix ORM queryset type issues
  - `orm/queryset.go:92` - OrderField type mismatch
  - `orm/queryset.go:129` - SelectRelated type mismatch
  - `orm/queryset.go:137` - PrefetchRelated type mismatch

### Architecture
- [ ] **Type registry** - Verify type registry works for all admin operations
- [ ] **Import cycles** - Ensure no import cycles exist
- [ ] **Error handling** - Consistent error handling across all packages

## Completed ✅

- ✅ Core admin registration (`admin.Register`)
- ✅ Schema-first design and auto-discovery
- ✅ ORM integration (AdminManager, AdminQuerySet)
- ✅ Filter integration (AdminFilterSet)
- ✅ List view with pagination, search, filtering
- ✅ Form view (add/edit) with fieldsets
- ✅ Detail view
- ✅ HTTP routing (CoreRouter)
- ✅ HTTP handlers (CoreHandler)
- ✅ Type registry for type-safe operations
- ✅ Permission system structure
- ✅ View hooks structure
- ✅ Widget system
- ✅ Component system (Table, Form, Pagination, FilterSidebar)
- ✅ Message system (flash messages)
- ✅ List editable handler implementation

## Notes

- Most core functionality is implemented
- Main gaps are in inline editing, validation, and some advanced features
- Template rendering needs verification
- Authentication/authorization needs integration testing
- ORM package errors need to be resolved (separate from admin package)
