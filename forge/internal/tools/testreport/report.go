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
				case "skip":
					c.skipped++
					for _, re := range required {
						if re.MatchString(e.Package) {
							forbidden = append(forbidden, e.Package+" "+e.Test)
							exitCode = 1
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
