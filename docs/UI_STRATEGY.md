# forge UI Strategy

## Philosophy: Django-Style Server-Rendered HTML

forge follows Django's approach: **server-rendered HTML with progressive enhancement**. No build tooling, no complex JavaScript frameworks, just fast, simple, and effective.

## Recommended Stack

### Core Approach: **HTMX + Bootstrap 5 + Server-Rendered HTML**

**Why this combination?**

1. **HTMX** - Already integrated! Provides modern interactivity without JavaScript complexity

   - No build tooling required
   - Works perfectly with Go templates
   - Progressive enhancement (works without JS)
   - Perfect for admin interfaces

2. **Bootstrap 5** - Professional, battle-tested CSS framework

   - No build step (CDN or static files)
   - Comprehensive component library
   - Responsive by default
   - Familiar to most developers
   - Works great with HTMX

3. **Server-Rendered HTML** - Go templates (html/template)
   - Fast initial page loads
   - SEO-friendly
   - Simple debugging
   - No hydration issues
   - Perfect for admin interfaces

## What NOT to Use (For Admin)

❌ **React/Vue/Angular** - Overkill for admin, requires build tooling, adds complexity
❌ **jQuery** - Outdated, HTMX is better
❌ **Complex build pipelines** - Keep it simple
❌ **SPA architecture** - Admin doesn't need it

## Implementation Plan

### Phase 1: Bootstrap Integration ✅

- [x] HTMX already integrated
- [x] Add Bootstrap 5 CSS/JS (CDN)
- [x] Update templates to use Bootstrap classes
- [x] Create reusable admin components

### Phase 2: HTMX Enhancements ✅

- [x] Inline editing with HTMX + Alpine.js
- [x] Dynamic filtering/search with debouncing
- [x] Delete confirmations with HTMX
- [x] Form validation feedback
- [x] Pagination with HTMX swaps

### Phase 3: Advanced Features ✅

- [x] DataTables.js for advanced table features (sorting, filtering, export)
- [x] Select2 for better select dropdowns
- [x] Flatpickr for date/time pickers
- [x] TinyMCE for rich text editing
- [x] Alpine.js for client-side interactivity

## Template Structure

```
templates/
├── base.html          # Base layout with Bootstrap + HTMX
├── list.html          # List view with Bootstrap table
├── form.html          # Form view with Bootstrap forms
├── detail.html        # Detail view
└── components/       # Reusable components
    ├── pagination.html
    ├── filters.html
    └── actions.html
```

## Example: Bootstrap + HTMX Pattern

```html
<!-- Inline edit with HTMX -->
<td
  hx-get="/admin/users/1/edit-field/email/"
  hx-target="this"
  hx-trigger="click"
  hx-swap="outerHTML"
>
  {{.Email}}
</td>

<!-- Delete with confirmation modal -->
<button
  hx-delete="/admin/users/1/delete/"
  hx-confirm="Are you sure?"
  hx-target="#main-content"
  class="btn btn-danger"
>
  Delete
</button>

<!-- Dynamic search -->
<input
  type="text"
  hx-get="/admin/users/"
  hx-trigger="keyup changed delay:500ms"
  hx-target="#user-list"
  class="form-control"
/>
```

## For User-Facing Apps (Not Admin)

If users want to build their own frontend:

### Option 1: REST API + Any Frontend

- forge generates REST APIs
- Users can use React/Vue/Angular/Next.js/etc.
- Framework provides JSON endpoints

### Option 2: Server-Rendered with HTMX

- Same approach as admin
- Fast, simple, SEO-friendly
- Perfect for content sites, dashboards

### Option 3: Hybrid

- Server-rendered for SEO pages
- HTMX for interactive parts
- API for mobile apps

## Benefits of This Approach

✅ **No build tooling** - Just Go templates + CDN
✅ **Fast development** - Bootstrap components ready to use
✅ **Modern UX** - HTMX provides SPA-like feel
✅ **Simple debugging** - View source, inspect HTML
✅ **Progressive enhancement** - Works without JavaScript
✅ **SEO-friendly** - Server-rendered HTML
✅ **Low complexity** - No webpack, no npm, no hydration

## Migration Path

Current state → Add Bootstrap → Enhance with HTMX → Polish

No breaking changes, just progressive enhancement.
