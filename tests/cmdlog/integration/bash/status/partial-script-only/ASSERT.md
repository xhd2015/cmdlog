## Expected

- `err == nil` and `resp.ExitCode == 1`.
- Stdout reports `bash integration: partial` with script present and marker absent.
- Pre-seeded `bash.sh` unchanged; no profile created.

## Side Effects

- Read-only inspection; no filesystem writes.

## Exit Code

- `1`

```go
import (
	"os"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp.ExitCode != 1 {
		t.Fatalf("expected exit 1, got %d; stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}

	assert.Output(t, resp.Stdout, `---
version: 2
__SCRIPT__: type=string, example=/tmp/home/.cmdlog/integration/bash.sh, bash.sh path
__PROFILE__: type=string, example=/tmp/home/.bash_profile, profile path
---
bash integration: partial
script: __SCRIPT__ (present)
profile: __PROFILE__ (marker absent)

`)

	if _, statErr := os.Stat(resp.BashShPath); statErr != nil {
		t.Fatalf("expected pre-seeded bash.sh present: %v", statErr)
	}
	if !strings.Contains(resp.BashShContent, "_cmdlog_capture") {
		t.Fatalf("status must not modify bash.sh:\n%s", resp.BashShContent)
	}
	if _, statErr := os.Stat(resp.ProfilePath); !os.IsNotExist(statErr) {
		t.Fatalf("status must not create .bash_profile at %s", resp.ProfilePath)
	}
	assertHomeIsolated(t, resp.BashShPath, resp.Home)
}
```