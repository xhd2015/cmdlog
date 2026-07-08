## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout previews would-write `bash.sh` and would-append marker block (no actual writes).
- `$HOME/.cmdlog/integration/bash.sh` does not exist.
- `$HOME/.bash_profile` does not exist.

## Side Effects

- No profile modifications.
- No `bash.sh` write.
- `events.jsonl` untouched (not created).

## Exit Code

- `0`

```go
import (
	"os"
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
dry-run: would write ~/.cmdlog/integration/bash.sh
dry-run: would append marker block to ~/.bash_profile

# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===

`)

	if _, statErr := os.Stat(resp.BashShPath); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create bash.sh at %s", resp.BashShPath)
	}
	if _, statErr := os.Stat(resp.ProfilePath); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create .bash_profile at %s", resp.ProfilePath)
	}
}
```