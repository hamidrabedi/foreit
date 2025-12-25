package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
	"unicode"
)

func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateToken(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func TruncateWords(s string, wordCount int) string {
	words := strings.Fields(s)
	if len(words) <= wordCount {
		return s
	}
	return strings.Join(words[:wordCount], " ") + "..."
}

func Pluralize(word string, count int) string {
	if count == 1 {
		return word
	}
	
	if strings.HasSuffix(word, "y") {
		return strings.TrimSuffix(word, "y") + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || 
		strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") || 
		strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0f seconds", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0f minutes", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.0f hours", d.Hours())
	}
	return fmt.Sprintf("%.0f days", d.Hours()/24)
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func ToSnakeCase(s string) string {
	var result strings.Builder
	result.Grow(len(s) + len(s)/3)
	
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

func ToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	
	return result
}

func ToPascalCase(s string) string {
	camel := ToCamelCase(s)
	if len(camel) == 0 {
		return camel
	}
	return strings.ToUpper(camel[:1]) + camel[1:]
}

func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0)
	
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	
	return result
}

func Map[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

func Filter[T any](slice []T, fn func(T) bool) []T {
	result := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Now() time.Time {
	return time.Now()
}

func UTCNow() time.Time {
	return time.Now().UTC()
}

func ParseTime(layout, value string) (time.Time, error) {
	return time.Parse(layout, value)
}

func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

func BeginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

func BeginningOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year++
	}
	return time.Date(year, nextMonth, 0, 23, 59, 59, 999999999, t.Location())
}

func BeginningOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, t.Location())
}

func IsToday(t time.Time) bool {
	now := Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}

func IsPast(t time.Time) bool {
	return t.Before(Now())
}

func IsFuture(t time.Time) bool {
	return t.After(Now())
}

