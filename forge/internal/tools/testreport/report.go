package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"
)

type options struct {
	out           string
	requireNoSkip patterns
}

type event struct {
	Action  string
	Package string
	Test    string
	Elapsed float64
	Output  string
}

type counts struct{ passed, failed, skipped, packageFailed int }

type testID struct {
	pkg  string
	test string
}

type section struct {
	header string
	lines  []string
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	s = strings.TrimSuffix(s, "\r\n")
	s = strings.TrimSuffix(s, "\n")
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSuffix(line, "\r")
	}
	return lines
}

func run(r io.Reader, w io.Writer, opts options) (exitCode int, err error) {
	var required []*regexp.Regexp
	for _, pattern := range opts.requireNoSkip {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return 2, fmt.Errorf("invalid --require-no-skip pattern: %w", err)
		}
		required = append(required, re)
	}
	packages := make(map[string]*counts)
	var forbidden []string

	testOutputs := make(map[testID]*strings.Builder)
	packageOutputs := make(map[string]*strings.Builder)

	var failedTests []testID
	failedSeen := make(map[testID]bool)

	var requiredSkips []testID
	skipSeen := make(map[testID]bool)

	reader := bufio.NewReader(r)
	for lineNo := 1; ; lineNo++ {
		line, readErr := reader.ReadString('\n')
		if readErr != nil && readErr != io.EOF {
			return 2, fmt.Errorf("read input: %w", readErr)
		}
		if strings.TrimSpace(line) != "" {
			var e event
			if err := json.Unmarshal([]byte(line), &e); err != nil {
				return 2, fmt.Errorf("line %d: %w", lineNo, err)
			}
			if e.Action == "" || e.Package == "" {
				return 2, fmt.Errorf("line %d: missing Action or Package", lineNo)
			}
			c := packages[e.Package]
			if c == nil {
				c = &counts{}
				packages[e.Package] = c
			}

			if e.Action == "output" {
				if e.Test != "" {
					id := testID{pkg: e.Package, test: e.Test}
					b := testOutputs[id]
					if b == nil {
						b = &strings.Builder{}
						testOutputs[id] = b
					}
					b.WriteString(e.Output)
				} else {
					b := packageOutputs[e.Package]
					if b == nil {
						b = &strings.Builder{}
						packageOutputs[e.Package] = b
					}
					b.WriteString(e.Output)
				}
			}

			if e.Test == "" {
				if e.Action == "fail" {
					c.packageFailed++
					exitCode = 1
				}
			} else {
				switch e.Action {
				case "pass":
					c.passed++
				case "fail":
					c.failed++
					exitCode = 1
					id := testID{pkg: e.Package, test: e.Test}
					if !failedSeen[id] {
						failedSeen[id] = true
						failedTests = append(failedTests, id)
					}
				case "skip":
					c.skipped++
					for _, re := range required {
						if re.MatchString(e.Package) {
							forbidden = append(forbidden, e.Package+" "+e.Test)
							exitCode = 1
							id := testID{pkg: e.Package, test: e.Test}
							if !skipSeen[id] {
								skipSeen[id] = true
								requiredSkips = append(requiredSkips, id)
							}
							break
						}
					}
				}
			}
		}
		if readErr == io.EOF {
			break
		}
	}
	var names []string
	for name := range packages {
		names = append(names, name)
	}
	sort.Strings(names)
	var report bytes.Buffer
	table := tabwriter.NewWriter(&report, 0, 4, 2, ' ', 0)
	fmt.Fprintln(table, "PACKAGE\tPASSED\tFAILED\tSKIPPED\tPACKAGE_FAILURES")
	var total counts
	row := func(name string, c counts) {
		fmt.Fprintf(table, "%s\t%d\t%d\t%d\t%d\n", name, c.passed, c.failed, c.skipped, c.packageFailed)
	}
	for _, name := range names {
		c := *packages[name]
		row(name, c)
		total.passed += c.passed
		total.failed += c.failed
		total.skipped += c.skipped
		total.packageFailed += c.packageFailed
	}
	row("TOTAL", total)
	if err := table.Flush(); err != nil {
		return 2, err
	}
	sort.Strings(forbidden)
	for _, test := range forbidden {
		fmt.Fprintf(&report, "REQUIRED TEST SKIPPED: %s\n", test)
	}

	var sections []section

	sort.Slice(failedTests, func(i, j int) bool {
		if failedTests[i].pkg != failedTests[j].pkg {
			return failedTests[i].pkg < failedTests[j].pkg
		}
		return failedTests[i].test < failedTests[j].test
	})
	for _, id := range failedTests {
		var raw string
		if b := testOutputs[id]; b != nil {
			raw = b.String()
		}
		sections = append(sections, section{
			header: fmt.Sprintf("=== FAIL %s %s", id.pkg, id.test),
			lines:  splitLines(raw),
		})
	}

	for _, name := range names {
		c := packages[name]
		if c.packageFailed > 0 && c.failed == 0 {
			var raw string
			if b := packageOutputs[name]; b != nil {
				raw = b.String()
			}
			sections = append(sections, section{
				header: fmt.Sprintf("=== PACKAGE FAIL %s", name),
				lines:  splitLines(raw),
			})
		}
	}

	sort.Slice(requiredSkips, func(i, j int) bool {
		if requiredSkips[i].pkg != requiredSkips[j].pkg {
			return requiredSkips[i].pkg < requiredSkips[j].pkg
		}
		return requiredSkips[i].test < requiredSkips[j].test
	})
	for _, id := range requiredSkips {
		var raw string
		if b := testOutputs[id]; b != nil {
			raw = b.String()
		}
		sections = append(sections, section{
			header: fmt.Sprintf("=== REQUIRED SKIP %s %s", id.pkg, id.test),
			lines:  splitLines(raw),
		})
	}

	const maxSections = 50
	const maxLines = 60

	toPrint := sections
	var omittedSections int
	if len(toPrint) > maxSections {
		omittedSections = len(toPrint) - maxSections
		toPrint = toPrint[:maxSections]
	}

	for _, sec := range toPrint {
		fmt.Fprintln(&report, sec.header)
		lines := sec.lines
		if len(lines) > maxLines {
			omitted := len(lines) - maxLines
			fmt.Fprintf(&report, "... %d earlier lines omitted\n", omitted)
			lines = lines[len(lines)-maxLines:]
		}
		for _, line := range lines {
			fmt.Fprintln(&report, line)
		}
	}

	if omittedSections > 0 {
		fmt.Fprintf(&report, "... %d more failures omitted\n", omittedSections)
	}

	if opts.out != "" {
		if err := os.WriteFile(opts.out, report.Bytes(), 0o644); err != nil {
			return 2, fmt.Errorf("write report: %w", err)
		}
	}
	if _, err := w.Write(report.Bytes()); err != nil {
		return 2, fmt.Errorf("write report: %w", err)
	}
	return exitCode, nil
}
