## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout contains `inside-window`.
- Stdout does **not** contain `before-window` or `after-window`.
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
	assertStdoutEndsWithNewline(t, resp.Stdout)

	if !strings.Contains(resp.Stdout, "inside-window") {
		t.Fatalf("expected inside-window in stdout:\n%s", resp.Stdout)
	}
	for _, excluded := range []string{"before-window", "after-window"} {
		if strings.Contains(resp.Stdout, excluded) {
			t.Fatalf("expected %q excluded from filtered timeline; stdout:\n%s", excluded, resp.Stdout)
		}
	}
}
```