## Expected

- `err == nil` and `resp.ExitCode == 0`.
- `$HOME/.cmdlog/integration/bash.sh` exists and contains `_cmdlog_capture`.
- `$HOME/.bash_profile` contains exactly one `# === cmdlog integration begin ===` marker.
- Profile sources bash.sh via the marker block.

## Side Effects

- `integration/bash.sh` written under fake HOME.
- Marker block appended to fake `.bash_profile`.

## Exit Code

- `0`

```go
import (
	"os"
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

	if _, statErr := os.Stat(resp.BashShPath); statErr != nil {
		t.Fatalf("expected bash.sh at %s: %v", resp.BashShPath, statErr)
	}
	if !strings.Contains(resp.BashShContent, "_cmdlog_capture") {
		t.Fatalf("bash.sh missing _cmdlog_capture:\n%s", resp.BashShContent)
	}
	if !strings.Contains(resp.BashShContent, "set +o history") {
		t.Fatalf("bash.sh missing set +o history around automation record:\n%s", resp.BashShContent)
	}
	if resp.MarkerCount != 1 {
		t.Fatalf("expected 1 marker block, got %d; profile:\n%s", resp.MarkerCount, resp.ProfileContent)
	}
	if !strings.Contains(resp.ProfileContent, ".cmdlog/integration/bash.sh") {
		t.Fatalf("profile must source bash.sh:\n%s", resp.ProfileContent)
	}
	assertHomeIsolated(t, resp.BashShPath, resp.Home)
	assertHomeIsolated(t, resp.ProfilePath, resp.Home)
}
```