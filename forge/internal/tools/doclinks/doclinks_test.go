package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Backtick and fence literals are spelled with hex escapes to keep the
// test fixtures explicit about which characters they contain.
const (
	bt    = "\x60"
	fence = "\x60\x60\x60"
)

func writeFixture(t *testing.T, root, rel, content string) {
	t.Helper()
	abs := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func checkFixture(t *testing.T, root string) []problem {
	t.Helper()
	_, problems, err := check(root)
	if err != nil {
		t.Fatal(err)
	}
	return problems
}

func TestValidRelativeLinkPasses(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\nSee [guide](docs/guide.md).\n")
	writeFixture(t, root, "docs/guide.md", "# Guide\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestMissingRelativeLinkReportedWithLineNumber(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\nIntro line.\nSee [gone](docs/nope.md).\n")
	problems := checkFixture(t, root)
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem, got %v", problems)
	}
	got := problems[0]
	if got.file != "README.md" || got.line != 4 || got.target != "docs/nope.md" {
		t.Fatalf("unexpected problem: %v", got)
	}
}

func TestExternalAndAnchorLinksSkipped(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n"+
		"[web](https://example.com/docs/guide.md)\n"+
		"[plain](http://example.com/x.md)\n"+
		"[mail](mailto:someone@example.com)\n"+
		"[anchor](#section)\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestDocsSiteRouteSkipped(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n[guide](/docs/guides/models)\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestBacktickedPaths(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n"+
		"See "+bt+"docs/MISSING.md"+bt+" for details.\n"+
		"Run "+bt+"go test ./..."+bt+" to test.\n")
	problems := checkFixture(t, root)
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem, got %v", problems)
	}
	got := problems[0]
	if got.file != "README.md" || got.line != 3 || got.target != "docs/MISSING.md" {
		t.Fatalf("unexpected problem: %v", got)
	}
}

func TestBacktickedPathsMustStartAtRepositoryRoot(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n"+
		"Skip "+bt+"config/config.yaml"+bt+" and "+bt+"models/post.go"+bt+".\n"+
		"Report "+bt+"docs/MISSING.md"+bt+".\n")
	problems := checkFixture(t, root)
	if len(problems) != 1 {
		t.Fatalf("expected 1 problem, got %v", problems)
	}
	if got := problems[0]; got.file != "README.md" || got.line != 4 || got.target != "docs/MISSING.md" {
		t.Fatalf("unexpected problem: %v", got)
	}
}

func TestIgnoredLineSkipped(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n"+
		"[gone](docs/nope.md) <!-- doclinks:ignore -->\n"+
		"See "+bt+"docs/MISSING.md"+bt+" <!-- doclinks:ignore -->\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestFencedCodeBlocksIgnored(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\n"+
		fence+"\n"+
		"[gone](docs/nope.md)\n"+
		"see "+bt+"docs/nope.md"+bt+"\n"+
		fence+"\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestFragmentStripped(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "docs/REVIEWING.md", "# Reviewing\n")
	writeFixture(t, root, "docs/page.md", "# Page\n\nSee [review](REVIEWING.md#acceptance-review).\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestQueryStripped(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "docs/guide.md", "# Guide\n")
	writeFixture(t, root, "docs/page.md", "# Page\n\nSee [guide](guide.md?x=1).\n")
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestDirectoryLinkPasses(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\nSee [sample](examples/ecommerce/).\n")
	if err := os.MkdirAll(filepath.Join(root, "examples", "ecommerce"), 0o755); err != nil {
		t.Fatal(err)
	}
	if problems := checkFixture(t, root); len(problems) != 0 {
		t.Fatalf("expected no problems, got %v", problems)
	}
}

func TestRunSuccessOutput(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\nSee [guide](docs/guide.md).\n")
	writeFixture(t, root, "docs/guide.md", "# Guide\n")
	var out bytes.Buffer
	code, err := run(root, &out)
	if err != nil {
		t.Fatal(err)
	}
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if got := out.String(); got != "doclinks: 2 files, 0 missing\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestRunReportsProblems(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n\nSee [gone](docs/nope.md).\n")
	var out bytes.Buffer
	code, err := run(root, &out)
	if err != nil {
		t.Fatal(err)
	}
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
	if got := out.String(); got != "README.md:3: missing docs/nope.md\n" {
		t.Fatalf("unexpected output: %q", got)
	}
}

func TestCollectFilesScope(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, "README.md", "# Title\n")
	writeFixture(t, root, "docs/top.md", "# Top\n")
	writeFixture(t, root, "docs/nested/deep.md", "# Deep\n")
	writeFixture(t, root, "skills/demo/SKILL.md", "# Skill\n")
	writeFixture(t, root, "skills/demo/references/ref.md", "# Ref\n")
	writeFixture(t, root, "skills/demo/references/notes.txt", "notes\n")
	writeFixture(t, root, "tests/README.md", "# Tests\n")
	files, err := collectFiles(root)
	if err != nil {
		t.Fatal(err)
	}
	var rels []string
	for _, f := range files {
		rel, err := filepath.Rel(root, f)
		if err != nil {
			t.Fatal(err)
		}
		rels = append(rels, filepath.ToSlash(rel))
	}
	got := strings.Join(rels, ",")
	want := "README.md,docs/top.md,skills/demo/SKILL.md,skills/demo/references/ref.md,tests/README.md"
	if got != want {
		t.Fatalf("got files %q, want %q", got, want)
	}
}
