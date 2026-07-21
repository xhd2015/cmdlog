## Expected

- `err == nil` and `resp.ExitCode == 1`.
- Stdout reports `bash integration: not installed` with script absent and marker absent.

## Side Effects

- Read-only inspection; no files created.

## Exit Code

- `1`

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
	if resp.ExitCode != 1 {
		t.Fatalf("expected exit 1, got %d; stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}

	assert.Output(t, resp.Stdout, `---
version: 3
__SCRIPT__: type=string, example=/tmp/home/.cmdlog/integration/bash.sh, bash.sh path
__PROFILE__: type=string, example=/tmp/home/.bash_profile, profile path
---
bash integration: not installed
script: __SCRIPT__ \(absent\)
profile: __PROFILE__ \(marker absent\)

`)

	if _, statErr := os.Stat(resp.BashShPath); !os.IsNotExist(statErr) {
		t.Fatalf("status must not create bash.sh at %s", resp.BashShPath)
	}
	if _, statErr := os.Stat(resp.ProfilePath); !os.IsNotExist(statErr) {
		t.Fatalf("status must not create .bash_profile at %s", resp.ProfilePath)
	}
}
```