## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout contains a friendly message (e.g. mentions no commands or empty log).
- Stdout ends with trailing newline.
- No events file created by query.

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
	assertStdoutEndsWithNewline(t, resp.Stdout)

	lower := strings.ToLower(resp.Stdout)
	if !strings.Contains(lower, "no") && !strings.Contains(lower, "empty") && !strings.Contains(lower, "0 command") {
		t.Fatalf("expected friendly empty-log message; stdout=%q", resp.Stdout)
	}
}
```