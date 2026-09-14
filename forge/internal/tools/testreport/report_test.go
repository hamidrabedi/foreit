package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func makeLongOutput(n int) string {
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "{\"Action\":\"output\",\"Package\":\"p\",\"Test\":\"TestLong\",\"Output\":\"line %d\\n\"}\n", i)
	}
	sb.WriteString("{\"Action\":\"fail\",\"Package\":\"p\",\"Test\":\"TestLong\"}\n")
	return sb.String()
}

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
		{
			"failing test output",
			"{\"Action\":\"output\",\"Package\":\"p\",\"Test\":\"TestBad\",\"Output\":\"bad_test.go:10: assertion failed\\n\"}\n{\"Action\":\"fail\",\"Package\":\"p\",\"Test\":\"TestBad\"}\n",
			nil,
			1,
			"=== FAIL p TestBad bad_test.go:10: assertion failed",
		},
		{
			"package failure output",
			"{\"Action\":\"output\",\"Package\":\"p\",\"Output\":\"syntax error: unexpected token\\n\"}\n{\"Action\":\"fail\",\"Package\":\"p\"}\n",
			nil,
			1,
			"=== PACKAGE FAIL p syntax error: unexpected token",
		},
		{
			"required skip output",
			"{\"Action\":\"output\",\"Package\":\"p\",\"Test\":\"TestDB\",\"Output\":\"skipping: db unavailable\\n\"}\n{\"Action\":\"skip\",\"Package\":\"p\",\"Test\":\"TestDB\"}\n",
			patterns{"^p$"},
			1,
			"=== REQUIRED SKIP p TestDB skipping: db unavailable",
		},
		{
			"60-line cap message",
			makeLongOutput(65),
			nil,
			1,
			"... 5 earlier lines omitted",
		},
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

func TestDetailedFailureReporting(t *testing.T) {
	t.Run("failing test output appears after fail header", func(t *testing.T) {
		input := `{"Action":"output","Package":"pkg/a","Test":"TestCrash","Output":"line from test crash\n"}
{"Action":"fail","Package":"pkg/a","Test":"TestCrash"}
`
		var out bytes.Buffer
		code, err := run(strings.NewReader(input), &out, options{})
		if code != 1 || err != nil {
			t.Fatalf("run = (%d, %v), want code 1", code, err)
		}
		s := out.String()
		headerIdx := strings.Index(s, "=== FAIL pkg/a TestCrash")
		outputIdx := strings.Index(s, "line from test crash")
		if headerIdx == -1 {
			t.Fatalf("missing === FAIL header in %q", s)
		}
		if outputIdx == -1 {
			t.Fatalf("missing output line in %q", s)
		}
		if outputIdx <= headerIdx {
			t.Fatalf("output line appears before or at header: headerIdx=%d, outputIdx=%d", headerIdx, outputIdx)
		}
	})

	t.Run("package-level failure prints output under package fail", func(t *testing.T) {
		input := `{"Action":"output","Package":"pkg/b","Output":"build error: cannot find package\n"}
{"Action":"fail","Package":"pkg/b"}
`
		var out bytes.Buffer
		code, err := run(strings.NewReader(input), &out, options{})
		if code != 1 || err != nil {
			t.Fatalf("run = (%d, %v), want code 1", code, err)
		}
		s := out.String()
		headerIdx := strings.Index(s, "=== PACKAGE FAIL pkg/b")
		outputIdx := strings.Index(s, "build error: cannot find package")
		if headerIdx == -1 {
			t.Fatalf("missing === PACKAGE FAIL header in %q", s)
		}
		if outputIdx == -1 {
			t.Fatalf("missing output line in %q", s)
		}
		if outputIdx <= headerIdx {
			t.Fatalf("output line appears before or at header: headerIdx=%d, outputIdx=%d", headerIdx, outputIdx)
		}
	})

	t.Run("required skip prints skip message", func(t *testing.T) {
		input := `{"Action":"output","Package":"pkg/c","Test":"TestNetwork","Output":"skipping: no network\n"}
{"Action":"skip","Package":"pkg/c","Test":"TestNetwork"}
`
		var out bytes.Buffer
		code, err := run(strings.NewReader(input), &out, options{requireNoSkip: patterns{"^pkg/c$"}})
		if code != 1 || err != nil {
			t.Fatalf("run = (%d, %v), want code 1", code, err)
		}
		s := out.String()
		headerIdx := strings.Index(s, "=== REQUIRED SKIP pkg/c TestNetwork")
		outputIdx := strings.Index(s, "skipping: no network")
		if headerIdx == -1 {
			t.Fatalf("missing === REQUIRED SKIP header in %q", s)
		}
		if outputIdx == -1 {
			t.Fatalf("missing output line in %q", s)
		}
		if outputIdx <= headerIdx {
			t.Fatalf("output line appears before or at header: headerIdx=%d, outputIdx=%d", headerIdx, outputIdx)
		}
	})

	t.Run("60-line cap omits earlier lines and keeps last 60", func(t *testing.T) {
		input := makeLongOutput(70)
		var out bytes.Buffer
		code, err := run(strings.NewReader(input), &out, options{})
		if code != 1 || err != nil {
			t.Fatalf("run = (%d, %v), want code 1", code, err)
		}
		s := out.String()
		wantCapMsg := "... 10 earlier lines omitted"
		if !strings.Contains(s, wantCapMsg) {
			t.Fatalf("expected %q in %q", wantCapMsg, s)
		}
		if strings.Contains(s, "\nline 1\n") || strings.Contains(s, "\nline 10\n") {
			t.Fatalf("earlier omitted lines should not appear in %q", s)
		}
		if !strings.Contains(s, "\nline 11\n") || !strings.Contains(s, "\nline 70\n") {
			t.Fatalf("retained lines (11..70) should appear in %q", s)
		}
	})

	t.Run("cap at 50 failure sections", func(t *testing.T) {
		var sb strings.Builder
		for i := 1; i <= 55; i++ {
			fmt.Fprintf(&sb, "{\"Action\":\"fail\",\"Package\":\"p\",\"Test\":\"TestFail%02d\"}\n", i)
		}
		var out bytes.Buffer
		code, err := run(strings.NewReader(sb.String()), &out, options{})
		if code != 1 || err != nil {
			t.Fatalf("run = (%d, %v), want code 1", code, err)
		}
		s := out.String()
		wantMoreMsg := "... 5 more failures omitted"
		if !strings.Contains(s, wantMoreMsg) {
			t.Fatalf("expected %q in %q", wantMoreMsg, s)
		}
		// Count occurrences of "=== FAIL"
		count := strings.Count(s, "=== FAIL")
		if count != 50 {
			t.Fatalf("expected 50 === FAIL sections, got %d", count)
		}
	})
}
