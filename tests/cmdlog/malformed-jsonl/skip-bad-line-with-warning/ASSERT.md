## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stderr contains a warning about malformed/skipped line.
- Stdout contains `valid-command` from the good JSONL line.
- Stdout ends with trailing newline.

## Exit Code

- `0`

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}

	lowerErr := strings.ToLower(resp.Stderr)
	if !strings.Contains(lowerErr, "malformed") && !strings.Contains(lowerErr, "skip") && !strings.Contains(lowerErr, "invalid") && !strings.Contains(lowerErr, "warn") {
		t.Fatalf("expected stderr warning for bad JSONL line; stderr=%q", resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "valid-command") {
		t.Fatalf("expected valid line processed; stdout=%q", resp.Stdout)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)
}
```