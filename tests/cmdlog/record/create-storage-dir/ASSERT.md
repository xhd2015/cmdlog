## Expected

- `err == nil` and `resp.ExitCode == 0`.
- `$HOME/.cmdlog/` directory exists after record.
- `$HOME/.cmdlog/events.jsonl` exists with one event (`cmd == "pwd"`).

## Side Effects

- `~/.cmdlog/` created automatically (was absent before Run).

## Exit Code

- `0`

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}

	cmdlogDir := filepath.Join(resp.Home, ".cmdlog")
	if info, statErr := os.Stat(cmdlogDir); statErr != nil || !info.IsDir() {
		t.Fatalf("expected .cmdlog directory created: %v", statErr)
	}
	if len(resp.EventLines) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.EventLines))
	}

	events := parseEvents(t, resp.EventLines)
	if events[0].CMD != "pwd" {
		t.Fatalf("expected cmd pwd, got %q", events[0].CMD)
	}
	assertHomeIsolated(t, resp.EventsPath, resp.Home)
}
```