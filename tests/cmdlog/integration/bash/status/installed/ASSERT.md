## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout reports `bash integration: installed` with script present and marker present.
- No filesystem changes from status itself.

## Side Effects

- Read-only inspection; profile and bash.sh unchanged.

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
__SCRIPT__: type=string, example=/tmp/home/.cmdlog/integration/bash.sh, bash.sh path
__PROFILE__: type=string, example=/tmp/home/.bash_profile, profile path
---
bash integration: installed
script: __SCRIPT__ \(present\)
profile: __PROFILE__ \(marker present\)

`)

	if resp.MarkerCount != 1 {
		t.Fatalf("expected installed marker; count=%d profile:\n%s", resp.MarkerCount, resp.ProfileContent)
	}
	if _, statErr := os.Stat(resp.BashShPath); statErr != nil {
		t.Fatalf("expected bash.sh present: %v", statErr)
	}
	assertHomeIsolated(t, resp.BashShPath, resp.Home)
	assertHomeIsolated(t, resp.ProfilePath, resp.Home)
}
```