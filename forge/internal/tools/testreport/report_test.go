package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	cases := []struct {
		name, input string
		required    patterns
		code        int
		want        string
	}{
		{"pass only", "{\"Action\":\"pass\",\"Package\":\"p\",\"Test\":\"TestOK\"}\n{\"Action\":\"pass\",\"Package\":\"p\"}\n", nil, 0, "TOTAL 1 0 0 0"},
		{"required skip", `{"Action":"skip","Package":"p","Test":"TestDB"}`, patterns{"^p$"}, 1, "REQUIRED TEST SKIPPED: p TestDB"},
		{"other skip", `{"Action":"skip","Package":"q","Test":"TestDB"}`, patterns{"^p$"}, 0, "TOTAL 0 0 1 0"},
		{"failed test", `{"Action":"fail","Package":"p","Test":"TestBad"}`, nil, 1, "TOTAL 0 1 0 0"},
		{"build failure", `{"Action":"fail","Package":"p"}`, nil, 1, "TOTAL 0 0 0 1"},
		{"malformed", "not json", nil, 2, ""},
		{"null", "null", nil, 2, ""},
		{"bad regexp", "", patterns{"["}, 2, ""},
		{"blank lines", "\n \n", nil, 0, "TOTAL 0 0 0 0"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			code, err := run(strings.NewReader(tc.input), &output, options{requireNoSkip: tc.required})
			if code != tc.code || (err != nil) != (tc.code == 2) {
				t.Fatalf("run = (%d, %v), want code %d", code, err, tc.code)
			}
			if !strings.Contains(strings.Join(strings.Fields(output.String()), " "), tc.want) {
				t.Fatalf("report %q missing %q", output.String(), tc.want)
			}
		})
	}
}

func TestReportSortedAndCopied(t *testing.T) {
	input := `{"Action":"pass","Package":"z","Test":"TestZ"}
{"Action":"skip","Package":"a","Test":"TestA"}
{"Action":"fail","Package":"z"}
`
	path := filepath.Join(t.TempDir(), "report.txt")
	var output bytes.Buffer
	code, err := run(strings.NewReader(input), &output, options{out: path})
	if code != 1 || err != nil {
		t.Fatalf("run = (%d, %v)", code, err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != output.String() {
		t.Fatalf("file differs from stdout: %v", err)
	}
	lines := strings.Split(output.String(), "\n")
	want := []string{"PACKAGE PASSED FAILED SKIPPED PACKAGE_FAILURES", "a 0 0 1 0", "z 1 0 0 1", "TOTAL 1 0 1 1"}
	for i, expected := range want {
		if got := strings.Join(strings.Fields(lines[i]), " "); got != expected {
			t.Errorf("line %d = %q, want %q", i, got, expected)
		}
	}
}
