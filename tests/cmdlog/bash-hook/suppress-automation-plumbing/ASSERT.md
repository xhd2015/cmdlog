## Expected

- `err == nil`.
- Exactly **2** events: `git status` and `cmdlog today`.
- No event whose `cmd` starts with `cmdlog record --cwd=` (hook automation plumbing).

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
	if len(resp.EventLines) != 2 {
		t.Fatalf("expected 2 user events, got %d:\n%s", len(resp.EventLines), resp.EventsContent)
	}

	events := parseEvents(t, resp.EventLines)
	cmds := eventCommands(events)

	want := []string{"git status", "cmdlog today"}
	for _, w := range want {
		found := false
		for _, c := range cmds {
			if c == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing user command %q in events %v", w, cmds)
		}
	}

	for _, ev := range events {
		if strings.HasPrefix(ev.CMD, "cmdlog record --cwd=") {
			t.Fatalf("hook automation plumbing must not appear as event: %q", ev.CMD)
		}
	}
}
```