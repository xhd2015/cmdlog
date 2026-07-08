## Expected

- `err == nil` and `resp.ExitCode == 0`.
- `events.jsonl` exists with exactly one line.
- Parsed event has `cwd == "/tmp"` and `cmd == "echo hi"`.
- Timestamp `ts` is valid UTC RFC3339.

## Side Effects

- `$HOME/.cmdlog/events.jsonl` created (append-only write).
- `$HOME/.cmdlog/` directory created if missing.

## Exit Code

- `0`

```go
import (
	"encoding/json"
	"testing"
	"time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}
	if len(resp.EventLines) != 1 {
		t.Fatalf("expected 1 event line, got %d: %s", len(resp.EventLines), resp.EventsContent)
	}

	var ev Event
	if err := json.Unmarshal([]byte(resp.EventLines[0]), &ev); err != nil {
		t.Fatalf("invalid JSONL: %v\n%s", err, resp.EventLines[0])
	}
	if ev.CWD != "/tmp" {
		t.Fatalf("expected cwd /tmp, got %q", ev.CWD)
	}
	if ev.CMD != "echo hi" {
		t.Fatalf("expected cmd 'echo hi', got %q", ev.CMD)
	}
	if _, parseErr := time.Parse(time.RFC3339, ev.TS); parseErr != nil {
		t.Fatalf("expected valid RFC3339 ts, got %q: %v", ev.TS, parseErr)
	}
	assertHomeIsolated(t, resp.EventsPath, resp.Home)
}
```