package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

// GetJSON parses JSON from request body
func GetJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.Unmarshal(body, v)
}

// GetQueryInt gets an integer query parameter
func GetQueryInt(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// GetQueryString gets a string query parameter
func GetQueryString(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryBool gets a boolean query parameter
func GetQueryBool(r *http.Request, key string, defaultValue bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

// GetParam gets a URL parameter (from chi router)
func GetParam(r *http.Request, key string) string {
	return URLParam(r, key)
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// SendError sends a JSON error response
func SendError(w http.ResponseWriter, statusCode int, message string) error {
	return SendJSON(w, statusCode, map[string]interface{}{
		"error":   true,
		"message": message,
		"code":    statusCode,
	})
}

// SendSuccess sends a JSON success response
func SendSuccess(w http.ResponseWriter, statusCode int, data interface{}) error {
	return SendJSON(w, statusCode, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}
