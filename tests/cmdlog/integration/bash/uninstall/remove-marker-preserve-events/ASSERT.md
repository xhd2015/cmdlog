## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Profile no longer contains `# === cmdlog integration begin ===` marker.
- `events.jsonl` still contains the `keep-me` event unchanged.

## Side Effects

- Marker block removed from `.bash_profile`.
- `events.jsonl` **not** deleted or truncated.

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
	if resp.MarkerCount != 0 {
		t.Fatalf("expected marker removed; count=%d profile:\n%s", resp.MarkerCount, resp.ProfileContent)
	}
	if !strings.Contains(resp.ProfileContent, "export EDITOR=vim") {
		t.Fatalf("uninstall must preserve unrelated profile content:\n%s", resp.ProfileContent)
	}
	if len(resp.EventLines) != 1 {
		t.Fatalf("events.jsonl must be preserved; got %d lines: %s", len(resp.EventLines), resp.EventsContent)
	}
	if !strings.Contains(resp.EventLines[0], "keep-me") {
		t.Fatalf("expected keep-me event preserved: %s", resp.EventLines[0])
	}
}
```