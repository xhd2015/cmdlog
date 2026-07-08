## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout previews marker block removal.
- Profile still contains exactly one marker block (unchanged).

## Side Effects

- No profile modifications.
- No `bash.sh` delete.
- `events.jsonl` untouched.

## Exit Code

- `0`

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}

	assert.Output(t, resp.Stdout, `---
version: 2
---
dry-run: would remove marker block from ~/.bash_profile

# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===

`)

	if resp.MarkerCount != 1 {
		t.Fatalf("dry-run must not remove marker; count=%d profile:\n%s", resp.MarkerCount, resp.ProfileContent)
	}
	if !strings.Contains(resp.ProfileContent, "export EDITOR=vim") {
		t.Fatalf("dry-run must preserve unrelated profile content:\n%s", resp.ProfileContent)
	}
	assertHomeIsolated(t, resp.ProfilePath, resp.Home)
}
```