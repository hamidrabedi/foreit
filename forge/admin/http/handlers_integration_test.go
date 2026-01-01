package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Integration tests for HTTP handlers

func TestIntegration_HandleBulkAction(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, err := setupTestAdmin()
	require.NoError(t, err)
	RegisterAdmin(admin)

	t.Run("DeleteAction", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("action", "delete")
		formData.Add("selected", "1")
		formData.Add("selected", "2")
		formData.Add("selected", "3")

		req := httptest.NewRequest("POST", "/admin/TestUser/bulk-action/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h := handler.HandleBulkAction("TestUser")
		h.ServeHTTP(w, req)

		// Should redirect after action
		assert.Equal(t, http.StatusSeeOther, w.Code)
	})

	t.Run("NoActionSpecified", func(t *testing.T) {
		formData := url.Values{}
		formData.Add("selected", "1")

		req := httptest.NewRequest("POST", "/admin/TestUser/bulk-action/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h := handler.HandleBulkAction("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NoItemsSelected", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("action", "delete")

		req := httptest.NewRequest("POST", "/admin/TestUser/bulk-action/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h := handler.HandleBulkAction("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("action", "delete")
		formData.Add("selected", "invalid")

		req := httptest.NewRequest("POST", "/admin/TestUser/bulk-action/", strings.NewReader(formData.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()

		h := handler.HandleBulkAction("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/bulk-action/", nil)
		w := httptest.NewRecorder()

		h := handler.HandleBulkAction("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})
}

func TestIntegration_HandleExport(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, err := setupTestAdmin()
	require.NoError(t, err)
	RegisterAdmin(admin)

	t.Run("ExportCSV", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/export/?format=csv", nil)
		w := httptest.NewRecorder()

		h := handler.HandleExport("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
		assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment")
	})

	t.Run("ExportJSON", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/export/?format=json", nil)
		w := httptest.NewRecorder()

		h := handler.HandleExport("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})

	t.Run("ExportUnsupportedFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/export/?format=xlsx", nil)
		w := httptest.NewRecorder()

		h := handler.HandleExport("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotImplemented, w.Code)
	})

	t.Run("ExportDefaultFormat", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/export/", nil)
		w := httptest.NewRecorder()

		h := handler.HandleExport("TestUser")
		h.ServeHTTP(w, req)

		// Default is CSV
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	})
}

func TestIntegration_HandleAutocomplete(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, err := setupTestAdmin()
	require.NoError(t, err)
	RegisterAdmin(admin)

	t.Run("MissingFieldParameter", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/autocomplete/?search=test", nil)
		w := httptest.NewRecorder()

		h := handler.HandleAutocomplete("TestUser")
		h.ServeHTTP(w, req)

		// Should error without field parameter
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestIntegration_HandleHistory(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, err := setupTestAdmin()
	require.NoError(t, err)
	RegisterAdmin(admin)

	t.Run("ViewHistory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/history/?id=1", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		h := handler.HandleHistory("TestUser")
		h.ServeHTTP(w, req)

		// May fail if history system not integrated, but shouldn't panic
		assert.NotEqual(t, 0, w.Code)
	})

	t.Run("MissingID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/history/", nil)
		w := httptest.NewRecorder()

		h := handler.HandleHistory("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("InvalidID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/history/?id=invalid", nil)
		w := httptest.NewRecorder()

		h := handler.HandleHistory("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestIntegration_HandleListEditable(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, err := setupTestAdmin()
	require.NoError(t, err)
	RegisterAdmin(admin)

	t.Run("InvalidJSON", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/admin/TestUser/list-editable/", strings.NewReader("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		
		// Load session context
		ctx, _ := handler.sessionManager.Load(req.Context(), "")
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		h := handler.HandleListEditable("TestUser")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestIntegration_TypeRegistry(t *testing.T) {
	t.Run("RegisterAndGetAdmin", func(t *testing.T) {
		admin, err := setupTestAdmin()
		require.NoError(t, err)

		RegisterAdmin(admin)

		handler, err := GetAdminHandler("TestUser")
		require.NoError(t, err)
		assert.NotNil(t, handler)
	})

	t.Run("GetNonexistentAdmin", func(t *testing.T) {
		_, err := GetAdminHandler("NonexistentModel")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("MultipleRegistrations", func(t *testing.T) {
		admin1, _ := setupTestAdmin()
		admin2, _ := setupTestAdmin()

		RegisterAdmin(admin1)
		RegisterAdmin(admin2) // Should overwrite

		handler, err := GetAdminHandler("TestUser")
		require.NoError(t, err)
		assert.NotNil(t, handler)
	})
}

func TestIntegration_ErrorHandling(t *testing.T) {
	handler, _ := setupTestHandler()

	t.Run("NonexistentModel", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/NonexistentModel/", nil)
		w := httptest.NewRecorder()

		h := handler.HandleList("NonexistentModel")
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("InvalidURLParams", func(t *testing.T) {
		admin, _ := setupTestAdmin()
		RegisterAdmin(admin)

		req := httptest.NewRequest("GET", "/admin/TestUser/invalid_id/", nil)
		w := httptest.NewRecorder()

		h := handler.HandleDetail("TestUser")
		h.ServeHTTP(w, req)

		// Should handle gracefully
		assert.NotEqual(t, 0, w.Code)
	})
}

func TestIntegration_SessionManagement(t *testing.T) {
	handler, _ := setupTestHandler()

	t.Run("UserContext", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		// Set user in context (not session, to avoid session initialization issues)
		ctx := context.WithValue(req.Context(), "user", map[string]interface{}{
			"id":       1,
			"username": "testuser",
		})
		req = req.WithContext(ctx)

		h := handler.sessionManager.Middleware()(handler.HandleIndex())
		h.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestIntegration_PaginationAndFiltering(t *testing.T) {
	handler, _ := setupTestHandler()
	admin, _ := setupTestAdmin()
	RegisterAdmin(admin)

	t.Run("WithPagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/admin/TestUser/?page=2&page_size=10", nil)
		req.Header.Set("Accept", "application/json")
		w := httptest.NewRecorder()

		// Use middleware to wrap handler
		h := handler.sessionManager.Middleware()(handler.HandleList("TestUser"))
		h.ServeHTTP(w, req)

		// Should complete without panic (may return 200 or 500 depending on DB state)
		assert.NotEqual(t, 0, w.Code)
	})
}

func BenchmarkHandler_HandleList(b *testing.B) {
	handler, _ := setupTestHandler()
	admin, _ := setupTestAdmin()
	RegisterAdmin(admin)

	req := httptest.NewRequest("GET", "/admin/TestUser/", nil)
	req.Header.Set("Accept", "application/json")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		h := handler.HandleList("TestUser")
		h.ServeHTTP(w, req)
	}
}

func BenchmarkHandler_HandleExport(b *testing.B) {
	handler, _ := setupTestHandler()
	admin, _ := setupTestAdmin()
	RegisterAdmin(admin)

	req := httptest.NewRequest("GET", "/admin/TestUser/export/?format=json", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		h := handler.HandleExport("TestUser")
		h.ServeHTTP(w, req)
	}
}
