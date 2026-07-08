## Expected Output

```text
---
version: 2
---
4 commands, 2 directories, top: /work/main (3)
```

## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout matches summary format with correct counts.
- Stdout ends with trailing newline.

## Exit Code

- `0`

```go
import (
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
4 commands, 2 directories, top: /work/main (3)
`)
}
```