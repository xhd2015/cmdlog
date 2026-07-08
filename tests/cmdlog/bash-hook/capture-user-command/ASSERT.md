## Expected

- `err == nil`.
- Exactly one event in `events.jsonl`.
- Event `cmd == "git status"`.

## Exit Code

- `0` (hook simulation succeeds)

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if len(resp.EventLines) != 1 {
		t.Fatalf("expected 1 event, got %d: %s", len(resp.EventLines), resp.EventsContent)
	}
	events := parseEvents(t, resp.EventLines)
	if events[0].CMD != "git status" {
		t.Fatalf("expected cmd 'git status', got %q", events[0].CMD)
	}
}
```