# Scenario

**Feature**: record auto-creates ~/.cmdlog when missing

```
# fresh HOME has no .cmdlog directory
cmdlog record --cwd=/tmp -- pwd -> creates ~/.cmdlog/ and events.jsonl
```

## Steps

1. Start with empty isolated HOME (no `.cmdlog` directory).
2. Run `cmdlog record --cwd=/tmp -- pwd`.

## Context

- Verifies directory creation side effect, not just file append.

```go
func Setup(t *testing.T, req *Request) error {
	req.RecordCWD = "/tmp"
	req.RecordCommand = "pwd"
	return nil
}
```