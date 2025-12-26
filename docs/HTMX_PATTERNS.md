# HTMX Patterns for forge Admin

This guide documents HTMX patterns and best practices for the forge admin interface.

## Overview

forge admin uses **HTMX + Bootstrap 5 + Alpine.js** for a modern, interactive admin interface without the complexity of React/Vue.

## Core Libraries

- **HTMX 1.9.10** - Server interactions and dynamic updates
- **Bootstrap 5.3.2** - UI components and styling
- **Alpine.js 3.x** - Lightweight client-side interactivity
- **DataTables.js** - Advanced table features
- **Select2** - Enhanced dropdowns
- **Flatpickr** - Date/time pickers
- **TinyMCE** - Rich text editing

## Common Patterns

### 1. Dynamic Search

Search as you type with automatic updates:

```html
<input
  type="text"
  name="search"
  value="{{.SearchQuery}}"
  placeholder="Search..."
  class="form-control"
  hx-get="{{.BaseURL}}"
  hx-trigger="keyup changed delay:500ms, search"
  hx-target="#table-container"
  hx-swap="outerHTML"
  hx-include="#search-form"
/>
```

**Features:**

- 500ms debounce to avoid excessive requests
- Updates only the table container (not full page)
- Includes form data automatically

### 2. Inline Editing

Edit fields directly in the list view:

```html
<td
  x-data="{ editing: false }"
  class="editable-cell"
  @dblclick="editing = true"
>
  <span x-show="!editing" @click="editing = true">{{.Value}}</span>
  <input
    x-show="editing"
    type="text"
    class="form-control form-control-sm"
    value="{{.Value}}"
    hx-post="/admin/users/1/edit-field/email/"
    hx-trigger="blur, keyup[key=='Enter']"
    hx-target="closest td"
    hx-swap="outerHTML"
    @keyup.escape="editing = false"
    name="value"
  />
</td>
```

**Features:**

- Double-click to edit
- Alpine.js manages edit state
- HTMX handles server update
- ESC key cancels editing

### 3. Delete with Confirmation

Delete rows with HTMX confirmation:

```html
<button
  type="button"
  class="btn btn-outline-danger"
  hx-delete="/admin/users/1/delete/"
  hx-confirm="Are you sure you want to delete this item?"
  hx-target="closest tr"
  hx-swap="outerHTML swap:1s"
>
  Delete
</button>
```

**Features:**

- Native browser confirmation dialog
- Removes row on success
- 1 second swap delay for visual feedback

### 4. Pagination with HTMX

Update table without full page reload:

```html
<a
  class="page-link"
  href="/admin/users/?page=2"
  hx-get="/admin/users/?page=2"
  hx-target="#table-container"
  hx-swap="outerHTML"
  hx-push-url="true"
  >Next</a
>
```

**Features:**

- Updates URL in browser (hx-push-url)
- Only swaps table container
- Maintains browser history

### 5. Form Validation Feedback

Show validation errors with HTMX:

```html
<form
  method="POST"
  action="/admin/users/new/"
  hx-post="/admin/users/new/"
  hx-target="#form-container"
  hx-swap="outerHTML"
>
  <!-- Form fields -->
</form>
```

**Server Response:**

- On error: Return form with error messages
- On success: Redirect or show success message

### 6. Toast Notifications

Show notifications for user actions:

```html
<div class="toast-container position-fixed bottom-0 end-0 p-3">
  <div id="toast-notification" class="toast" role="alert">
    <div class="toast-header">
      <strong class="me-auto">Notification</strong>
      <button type="button" class="btn-close" data-bs-dismiss="toast"></button>
    </div>
    <div class="toast-body"></div>
  </div>
</div>

<script>
  document.body.addEventListener("htmx:responseError", function (evt) {
    const toast = new bootstrap.Toast(
      document.getElementById("toast-notification")
    );
    const toastBody = document.querySelector("#toast-notification .toast-body");
    toastBody.textContent = "An error occurred. Please try again.";
    toast.show();
  });
</script>
```

## HTMX + Alpine.js Patterns

### Toggle Visibility

```html
<div x-data="{ show: false }">
  <button @click="show = !show">Toggle</button>
  <div x-show="show" x-transition>Content</div>
</div>
```

### Conditional Rendering

```html
<div x-data="{ mode: 'view' }">
  <span x-show="mode === 'view'" @click="mode = 'edit'">{{.Value}}</span>
  <input
    x-show="mode === 'edit'"
    hx-post="/update/"
    hx-trigger="blur"
    @blur="mode = 'view'"
  />
</div>
```

### Loading States

```html
<button hx-post="/save/" hx-indicator="#spinner">Save</button>
<div id="spinner" class="htmx-indicator">
  <span class="spinner-border spinner-border-sm"></span>
</div>
```

## DataTables Integration

Initialize DataTables after HTMX swaps:

```javascript
document.body.addEventListener("htmx:afterSwap", function (evt) {
  if (evt.detail.target.id === "table-container") {
    // Destroy existing DataTable
    if ($.fn.DataTable.isDataTable("#data-table")) {
      $("#data-table").DataTable().destroy();
    }
    // Initialize new DataTable
    $("#data-table").DataTable({
      pageLength: 20,
      order: [[0, "asc"]],
    });
  }
});
```

## Form Libraries Integration

### Select2

```javascript
$(document).ready(function () {
  $(".select2-field").select2({
    theme: "bootstrap-5",
    width: "100%",
  });
});
```

### Flatpickr

```javascript
flatpickr(".flatpickr-field", {
  enableTime: true,
  dateFormat: "Y-m-d H:i",
  time_24hr: true,
});
```

### TinyMCE

```javascript
tinymce.init({
  selector: ".tinymce-field",
  height: 400,
  plugins: ["advlist", "autolink", "lists", "link", "image"],
  toolbar: "undo redo | blocks | bold italic | alignleft aligncenter",
});
```

## Best Practices

### 1. Target Specific Elements

Always target specific containers, not the entire page:

```html
<!-- Good -->
hx-target="#table-container"

<!-- Bad -->
hx-target="body"
```

### 2. Use Debouncing for Search

Prevent excessive requests:

```html
hx-trigger="keyup changed delay:500ms"
```

### 3. Handle HTMX Requests on Server

Check for HTMX requests and return partial HTML:

```go
isHTMXRequest := r.Header.Get("HX-Request") == "true"
if isHTMXRequest {
    // Return only the updated section
    return renderPartial(w, data)
}
// Return full page
return renderFullPage(w, data)
```

### 4. Use hx-push-url for Navigation

Maintain browser history:

```html
hx-push-url="true"
```

### 5. Show Loading Indicators

Provide visual feedback:

```html
hx-indicator="#spinner"
```

### 6. Handle Errors Gracefully

Use HTMX events for error handling:

```javascript
document.body.addEventListener("htmx:responseError", function (evt) {
  // Show error message
});
```

## Common Issues and Solutions

### Issue: DataTables Not Initializing After HTMX Swap

**Solution:** Reinitialize in `htmx:afterSwap` event:

```javascript
document.body.addEventListener("htmx:afterSwap", function (evt) {
  if (evt.detail.target.id === "table-container") {
    $("#data-table").DataTable().destroy();
    $("#data-table").DataTable({
      /* config */
    });
  }
});
```

### Issue: Form Libraries Not Working After HTMX Swap

**Solution:** Reinitialize libraries after swap:

```javascript
document.body.addEventListener("htmx:afterSwap", function (evt) {
  $(".select2-field").select2();
  flatpickr(".flatpickr-field", {
    /* config */
  });
});
```

### Issue: Alpine.js Not Working After HTMX Swap

**Solution:** Alpine.js auto-initializes, but ensure x-data is in swapped content:

```html
<div x-data="{ ... }" hx-swap="outerHTML">
  <!-- Content -->
</div>
```

## Performance Tips

1. **Use hx-swap="outerHTML"** for replacing entire elements
2. **Use hx-swap="innerHTML"** for updating content only
3. **Debounce search inputs** to reduce server load
4. **Target specific containers** to minimize DOM updates
5. **Use hx-indicator** for visual feedback during requests

## Examples

See the admin templates for complete examples:

- `internal/admin/templates/list.html` - List view with search, pagination, inline editing
- `internal/admin/templates/form.html` - Form with validation and enhanced inputs

## Further Reading

- [HTMX Documentation](https://htmx.org/docs/)
- [Alpine.js Documentation](https://alpinejs.dev/)
- [DataTables Documentation](https://datatables.net/)
- [Select2 Documentation](https://select2.org/)
