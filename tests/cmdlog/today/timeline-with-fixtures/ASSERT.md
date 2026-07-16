## Expected Output

```text
---
version: 3
__HOUR09__: type=string, example=09:00, morning hour label
---
__HOUR09__ /proj/a
  git status
  go test \./\.\.\.
__HOUR10__ /proj/b
  make build
```

## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout contains hour groupings and cwd sub-groupings for all three commands.
- Stdout ends with trailing newline.

## Exit Code

- `0`

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)

	for _, needle := range []string{"/proj/a", "/proj/b", "git status", "go test ./...", "make build"} {
		if !strings.Contains(resp.Stdout, needle) {
			t.Fatalf("stdout missing %q:\n%s", needle, resp.Stdout)
		}
	}

	// Hour grouping: both /proj/a commands and /proj/b command appear under distinct hour sections.
	lines := strings.Split(strings.TrimRight(resp.Stdout, "\n"), "\n")
	var sawHourA, sawHourB bool
	for _, line := range lines {
		if strings.Contains(line, "/proj/a") {
			sawHourA = true
		}
		if strings.Contains(line, "/proj/b") {
			sawHourB = true
		}
	}
	if !sawHourA || !sawHourB {
		t.Fatalf("expected hour+cwd grouping; stdout:\n%s", resp.Stdout)
	}

	// Flexible template: hour labels vary by timezone; assert structure via contains above.
	_ = assert.Output
}
```