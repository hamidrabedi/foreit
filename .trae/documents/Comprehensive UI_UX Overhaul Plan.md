I have performed a comprehensive audit and identified several critical issues (broken Search, broken Relations, Sidebar bugs) alongside the requested enhancements. Here is the master plan to fix everything and deliver the requested functionality.

## Phase 1: Fix Broken Core Features

1. **Unify API Architecture**:

   * **Issue**: `GlobalSearch` and `SearchableSelect` use a raw, misconfigured `axiosInstance` that points to the wrong URLs (`/api/...` instead of `/admin/api/...`), causing 404 errors.

   * **Fix**: Refactor these components to use the centralized `adminAPI` client which is correctly configured.
2. **Fix Sidebar Navigation**:

   * **Issue**: Nested plugin menus (e.g., Reports -> Sales) are not rendered.

   * **Fix**: Refactor `AdminLayout` sidebar to support recursive/nested menu rendering.

## Phase 2: Interactive Features & UI Components

1. **Install Missing UI Libraries**:

   * Run `npm install @radix-ui/react-alert-dialog` to enable professional modal dialogs.
2. **Create Confirmation System**:

   * Implement `src/components/ui/confirmation-dialog.tsx`.

   * Replace all `window.confirm()` calls in `ModelListPage` (Delete) and `ModelUpsertPage` (Unsaved changes) with this new dialog.

## Phase 3: Advanced Functionality

1. **Advanced Filtering**:

   * Update `ModelListPage` to support **Date Range Filters** (Start/End date) instead of just simple text matching.

   * Add "Select" filters for boolean fields.
2. **Enhanced Validation**:

   * Improve `ModelUpsertPage` to display field-level validation errors returned from the API more effectively.

## Phase 4: Final Verification

1. **Responsiveness Check**: Verify mobile sidebar toggle and table scrolling on small screens.
2. **Full Regression Test**: Manually verify Search, Relations, CRUD, and Navigation.

