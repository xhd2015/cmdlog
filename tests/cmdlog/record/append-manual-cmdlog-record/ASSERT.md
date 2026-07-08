## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Exactly one event appended.
- Event `cmd` equals the full string `cmdlog record --cwd=/x -- foo` (not rejected).

## Exit Code

- `0`

```go
import (
	"encoding/json"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}
	if len(resp.EventLines) != 1 {
		t.Fatalf("expected 1 event, got %d", len(resp.EventLines))
	}

	var ev Event
	if err := json.Unmarshal([]byte(resp.EventLines[0]), &ev); err != nil {
		t.Fatalf("invalid JSONL: %v", err)
	}
	want := "cmdlog record --cwd=/x -- foo"
	if ev.CMD != want {
		t.Fatalf("manual cmdlog record must be recorded; got cmd=%q want=%q", ev.CMD, want)
	}
}
```