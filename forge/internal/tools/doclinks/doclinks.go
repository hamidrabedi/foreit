package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// problem is a single missing link target.
type problem struct {
	file   string // path relative to root, slash-separated
	line   int    // 1-based line number
	target string // cleaned link target (fragment/query stripped)
}

func (p problem) String() string {
	return fmt.Sprintf("%s:%d: missing %s", p.file, p.line, p.target)
}

var (
	linkRe         = regexp.MustCompile(`\[[^\]]*\]\(([^)]*)\)`)
	referenceDefRe = regexp.MustCompile(`^\s*\[[^\]]+\]:\s*(<[^>]+>|\S+)`)
	tickRe         = regexp.MustCompile(`\x60([^\x60]+)\x60`)
)

const fenceMarker = "\x60\x60\x60"

var backtickSuffixes = []string{".md", ".go", ".yml", ".yaml", ".json"}

var backtickPathRoots = map[string]bool{
	"forge":     true,
	"docs":      true,
	"docs-site": true,
	"skills":    true,
	"tests":     true,
	"examples":  true,
	"scripts":   true,
	".github":   true,
}

// collectFiles returns the sorted list of absolute paths to scan under root.
func collectFiles(root string) ([]string, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absRoot)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("doclinks: root %s is not a directory", root)
	}
	var files []string
	seen := map[string]bool{}
	add := func(rel string) {
		abs := filepath.Join(absRoot, filepath.FromSlash(rel))
		info, err := os.Stat(abs)
		if err != nil || info.IsDir() {
			return
		}
		if !seen[abs] {
			seen[abs] = true
			files = append(files, abs)
		}
	}
	for _, rel := range []string{
		"README.md",
		"AGENTS.md",
		"CONTRIBUTING.md",
		"SECURITY.md",
	} {
		add(rel)
	}
	for _, tree := range []string{"docs", "docs-site/docs", "skills", "tests"} {
		treePath := filepath.Join(absRoot, filepath.FromSlash(tree))
		if _, err := os.Stat(treePath); os.IsNotExist(err) {
			continue
		}
		if err := filepath.WalkDir(treePath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				if path != treePath && (strings.HasPrefix(d.Name(), ".") || d.Name() == "node_modules" || d.Name() == "dist" || d.Name() == "build") {
					return filepath.SkipDir
				}
				return nil
			}
			if strings.HasSuffix(d.Name(), ".md") {
				rel, err := filepath.Rel(absRoot, path)
				if err != nil {
					return err
				}
				add(filepath.ToSlash(rel))
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	sort.Strings(files)
	return files, nil
}

// targetExists reports whether path exists as a file or a directory.
// Links to existing directories (for example "examples/ecommerce/")
// are valid targets, so directories count as existing.
func targetExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// cleanLinkTarget normalizes a markdown link target. It reports false for
// targets the checker must skip (external URLs, anchors, empty targets).
func cleanLinkTarget(raw string) (string, bool) {
	t := strings.TrimSpace(raw)
	if t == "" {
		return "", false
	}
	if strings.HasPrefix(t, "<") {
		if i := strings.Index(t, ">"); i >= 0 {
			t = t[1:i]
		}
	} else if fields := strings.Fields(t); len(fields) > 0 {
		t = fields[0]
	}
	for _, prefix := range []string{"http://", "https://", "mailto:", "#"} {
		if strings.HasPrefix(t, prefix) {
			return "", false
		}
	}
	if strings.HasPrefix(t, "/") {
		return "", false
	}
	if i := strings.Index(t, "#"); i >= 0 {
		t = t[:i]
	}
	if i := strings.Index(t, "?"); i >= 0 {
		t = t[:i]
	}
	t = strings.TrimSpace(t)
	if t == "" {
		return "", false
	}
	return t, true
}

// backtickCandidate reports whether a backticked token looks like a
// repository file path worth checking.
func backtickCandidate(tok string) bool {
	if !strings.Contains(tok, "/") {
		return false
	}
	matched := strings.HasSuffix(tok, "/")
	for _, suffix := range backtickSuffixes {
		if strings.HasSuffix(tok, suffix) {
			matched = true
			break
		}
	}
	if !matched {
		return false
	}
	if strings.ContainsAny(tok, " *<>{\x60$") {
		return false
	}
	firstSegment, _, _ := strings.Cut(tok, "/")
	return backtickPathRoots[firstSegment]
}

// checkFile scans a single file for missing link targets.
func checkFile(absRoot, absPath string) ([]problem, error) {
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return nil, err
	}
	rel = filepath.ToSlash(rel)
	dir := filepath.Dir(absPath)
	var problems []problem
	inFence := false
	for i, line := range strings.Split(string(data), "\n") {
		if strings.Contains(line, "<!-- doclinks:ignore -->") {
			continue
		}
		if strings.Contains(line, fenceMarker) {
			if strings.Count(line, fenceMarker)%2 == 1 {
				inFence = !inFence
			}
			continue
		}
		if inFence {
			continue
		}
		lineno := i + 1
		for _, m := range linkRe.FindAllStringSubmatch(line, -1) {
			target, ok := cleanLinkTarget(m[1])
			if !ok {
				continue
			}
			if !targetExists(filepath.Join(dir, filepath.FromSlash(target))) {
				problems = append(problems, problem{file: rel, line: lineno, target: target})
			}
		}
		if m := referenceDefRe.FindStringSubmatch(line); m != nil {
			target, ok := cleanLinkTarget(m[1])
			if ok && !targetExists(filepath.Join(dir, filepath.FromSlash(target))) {
				problems = append(problems, problem{file: rel, line: lineno, target: target})
			}
		}
		for _, m := range tickRe.FindAllStringSubmatch(line, -1) {
			tok := strings.TrimSpace(m[1])
			if !backtickCandidate(tok) {
				continue
			}
			// Markdown links are relative to their containing document, while
			// backticked repository paths conventionally start at repository root.
			if !targetExists(filepath.Join(absRoot, filepath.FromSlash(tok))) {
				problems = append(problems, problem{file: rel, line: lineno, target: tok})
			}
		}
	}
	return problems, nil
}

// check scans every documentation file under root and returns the missing
// link targets in deterministic (file, line) order.
func check(root string) (int, []problem, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return 0, nil, err
	}
	files, err := collectFiles(absRoot)
	if err != nil {
		return 0, nil, err
	}
	if len(files) == 0 {
		return 0, nil, fmt.Errorf("doclinks: no documentation files found under %s", root)
	}
	var problems []problem
	for _, f := range files {
		found, err := checkFile(absRoot, f)
		if err != nil {
			return 0, nil, err
		}
		problems = append(problems, found...)
	}
	sort.Slice(problems, func(i, j int) bool {
		if problems[i].file != problems[j].file {
			return problems[i].file < problems[j].file
		}
		if problems[i].line != problems[j].line {
			return problems[i].line < problems[j].line
		}
		return problems[i].target < problems[j].target
	})
	return len(files), problems, nil
}

func run(root string, w io.Writer) (int, error) {
	n, problems, err := check(root)
	if err != nil {
		return 2, err
	}
	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(w, p.String())
		}
		return 1, nil
	}
	fmt.Fprintf(w, "doclinks: %d files, 0 missing\n", n)
	return 0, nil
}
