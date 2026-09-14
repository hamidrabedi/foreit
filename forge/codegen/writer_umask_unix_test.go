//go:build unix

package generator

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

func TestWriteTemplate_RespectsRestrictiveUmask(t *testing.T) {
	// Do not run in parallel: umask is process-wide.
	previous := syscall.Umask(0o077)
	t.Cleanup(func() { syscall.Umask(previous) })
	dest := filepath.Join(t.TempDir(), "new.go")
	tmpl := template.Must(template.New("new").Parse("package test\n"))
	require.NoError(t, NewWriter().writeTemplate(tmpl, nil, dest))
	info, err := os.Stat(dest)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0o600), info.Mode().Perm())
}
