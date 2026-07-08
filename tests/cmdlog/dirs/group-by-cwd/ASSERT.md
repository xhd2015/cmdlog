## Expected

- `err == nil` and `resp.ExitCode == 0`.
- Stdout groups commands under `/alpha` and `/beta` headers/sections.
- All three commands appear under their respective cwd groups.
- Stdout ends with trailing newline.

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
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d; stderr=%s", resp.ExitCode, resp.Stderr)
	}
	assertStdoutEndsWithNewline(t, resp.Stdout)

	alphaIdx := strings.Index(resp.Stdout, "/alpha")
	betaIdx := strings.Index(resp.Stdout, "/beta")
	if alphaIdx < 0 || betaIdx < 0 {
		t.Fatalf("expected /alpha and /beta in stdout:\n%s", resp.Stdout)
	}

	alphaSection := resp.Stdout[alphaIdx:betaIdx]
	if !strings.Contains(alphaSection, "npm install") || !strings.Contains(alphaSection, "npm test") {
		t.Fatalf("expected /alpha section to list both npm commands:\n%s", alphaSection)
	}
	if !strings.Contains(resp.Stdout[betaIdx:], "cargo build") {
		t.Fatalf("expected /beta section to list cargo build:\n%s", resp.Stdout)
	}
}
```