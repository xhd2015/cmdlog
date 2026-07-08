## Expected

- `err == nil` (Run completes; CLI failure is in `resp.ExitCode`).
- `resp.ExitCode != 0`.
- Stderr or stdout mentions "not yet supported" (case-insensitive).

## Exit Code

- Non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit for unsupported zsh install; stdout=%q stderr=%q",
			resp.Stdout, resp.Stderr)
	}
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, "not yet supported") {
		t.Fatalf("expected 'not yet supported' message; stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
}
```