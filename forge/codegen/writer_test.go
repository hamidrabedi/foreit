package generator

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteTemplate_FailurePreservesDestinationAndCleansTemp(t *testing.T) {
	tmpDir := t.TempDir()
	dest := filepath.Join(tmpDir, "output.go")

	initialContent := []byte("// initial valid file content\npackage test\n")
	require.NoError(t, os.WriteFile(dest, initialContent, 0640))

	failTmpl := template.Must(template.New("fail").Funcs(template.FuncMap{
		"failFunc": func() (string, error) {
			return "", errors.New("template execution failure")
		},
	}).Parse("partially written text before error: {{failFunc}} after text"))

	w := NewWriter()
	err := w.writeTemplate(failTmpl, nil, dest)
	require.Error(t, err)

	// Verify destination file is byte-for-byte unchanged
	content, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, initialContent, content, "destination file must be byte-for-byte unchanged")

	// Verify permissions are preserved
	fi, err := os.Stat(dest)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0640), fi.Mode().Perm())

	// Verify no temp files left behind
	matches, err := filepath.Glob(filepath.Join(tmpDir, ".*tmp*"))
	require.NoError(t, err)
	assert.Empty(t, matches, "no temp files should remain")

	entries, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.Len(t, entries, 1)
	assert.Equal(t, "output.go", entries[0].Name())
}

func TestWriteTemplate_SuccessReplacesContentAndPreservesPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	dest := filepath.Join(tmpDir, "success.go")

	initialContent := []byte("// old content\n")
	require.NoError(t, os.WriteFile(dest, initialContent, 0750))

	succTmpl := template.Must(template.New("success").Parse("package {{.Package}}\n// new content\n"))
	data := map[string]string{"Package": "mypkg"}

	w := NewWriter()
	err := w.writeTemplate(succTmpl, data, dest)
	require.NoError(t, err)

	// Verify destination content is replaced fully
	content, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "package mypkg\n// new content\n", string(content))

	// Verify previous file's permissions are preserved
	fi, err := os.Stat(dest)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0750), fi.Mode().Perm())

	// Verify no temp files left behind
	matches, err := filepath.Glob(filepath.Join(tmpDir, ".*tmp*"))
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestWriteTemplate_DefaultPermissionsWhenFileDoesNotExist(t *testing.T) {
	tmpDir := t.TempDir()
	dest := filepath.Join(tmpDir, "new_file.go")

	succTmpl := template.Must(template.New("new").Parse("package test\n"))

	w := NewWriter()
	err := w.writeTemplate(succTmpl, nil, dest)
	require.NoError(t, err)

	// Verify file was created
	content, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "package test\n", string(content))

	// Verify new-file permissions respect the process umask.
	fi, err := os.Stat(dest)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o666&^processUmask()), fi.Mode().Perm())

	// Verify no temp files left behind
	matches, err := filepath.Glob(filepath.Join(tmpDir, ".*tmp*"))
	require.NoError(t, err)
	assert.Empty(t, matches)
}

func TestWriteTemplate_UpdatesSymlinkTargetWithoutReplacingLink(t *testing.T) {
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "generated")
	require.NoError(t, os.Mkdir(targetDir, 0755))
	target := filepath.Join(targetDir, "models.go")
	require.NoError(t, os.WriteFile(target, []byte("old contents\n"), 0644))
	link := filepath.Join(tmpDir, "gen.go")
	require.NoError(t, os.Symlink(filepath.Join("generated", "models.go"), link))

	tmpl := template.Must(template.New("symlink").Parse("new contents\n"))
	require.NoError(t, NewWriter().writeTemplate(tmpl, nil, link))

	contents, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "new contents\n", string(contents))
	info, err := os.Lstat(link)
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink)
}
