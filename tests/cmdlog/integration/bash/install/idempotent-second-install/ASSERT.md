## Expected

- `err == nil` and both install runs exit 0.
- `resp.MarkerCount == 1` after second install (no duplicate markers).
- `bash.sh` still exists and is unchanged functionally.

## Exit Code

- `0`

```go
import (
	"os"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}
	if resp.MarkerCount != 1 {
		t.Fatalf("idempotent install must not duplicate marker; count=%d profile:\n%s",
			resp.MarkerCount, resp.ProfileContent)
	}
	if _, statErr := os.Stat(resp.BashShPath); statErr != nil {
		t.Fatalf("bash.sh missing after idempotent install: %v", statErr)
	}
}
```