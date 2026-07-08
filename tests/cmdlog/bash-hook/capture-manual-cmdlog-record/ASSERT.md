## Expected

- `err == nil`.
- Event recorded with full command string containing `hello-manual`.
- Manual `cmdlog record` is **not** filtered out by hook or Go record subcommand.

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
	if len(resp.EventLines) != 1 {
		t.Fatalf("expected 1 event, got %d: %s", len(resp.EventLines), resp.EventsContent)
	}
	events := parseEvents(t, resp.EventLines)
	if !strings.Contains(events[0].CMD, "hello-manual") {
		t.Fatalf("expected manual cmdlog record logged; cmd=%q", events[0].CMD)
	}
}
```