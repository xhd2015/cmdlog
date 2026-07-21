## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout reports already uninstalled; no changes needed.
- No profile or bash.sh created.

## Side Effects

- No profile modifications.
- No `bash.sh` write/delete.
- `events.jsonl` untouched.

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
version: 3
__PROFILE__: type=string, example=/tmp/home/.bash_profile, profile path
---
cmdlog bash integration: already uninstalled
profile: __PROFILE__ \(marker absent\)
no changes needed
`)

	if _, statErr := os.Stat(resp.ProfilePath); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create .bash_profile at %s", resp.ProfilePath)
	}
	if _, statErr := os.Stat(resp.BashShPath); !os.IsNotExist(statErr) {
		t.Fatalf("dry-run must not create bash.sh at %s", resp.BashShPath)
	}
}
```
