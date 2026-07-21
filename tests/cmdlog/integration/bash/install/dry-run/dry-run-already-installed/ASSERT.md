## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout reports already installed with profile and script paths; no changes needed.
- Pre-seeded `bash.sh` and profile content unchanged.

## Side Effects

- No profile modifications.
- No `bash.sh` overwrite.
- `events.jsonl` untouched.

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
	for _, want := range []string{
		`cmdlog bash integration: already installed`,
		`(marker present)`,
		`(exists)`,
		`no changes needed`,
	} {
		if !strings.Contains(resp.Stdout, want) {
			t.Fatalf("stdout missing %q:\n%s", want, resp.Stdout)
		}
	}
	if !strings.Contains(resp.Stdout, "profile: ") || !strings.Contains(resp.Stdout, "script: ") {
		t.Fatalf("stdout missing profile/script lines:\n%s", resp.Stdout)
	}
	if resp.MarkerCount != 1 {
		t.Fatalf("dry-run must not change marker count; got %d profile:\n%s", resp.MarkerCount, resp.ProfileContent)
	}
	if !strings.Contains(resp.ProfileContent, "export EDITOR=vim") {
		t.Fatalf("dry-run must preserve unrelated profile content:\n%s", resp.ProfileContent)
	}
	if !strings.Contains(resp.BashShContent, "_cmdlog_capture") {
		t.Fatalf("dry-run must not overwrite bash.sh:\n%s", resp.BashShContent)
	}
	if _, statErr := os.Stat(resp.BashShPath); statErr != nil {
		t.Fatalf("bash.sh must still exist: %v", statErr)
	}
	assertHomeIsolated(t, resp.BashShPath, resp.Home)
	assertHomeIsolated(t, resp.ProfilePath, resp.Home)
}
```
