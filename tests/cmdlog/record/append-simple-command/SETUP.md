# Scenario

**Feature**: basic record appends one JSONL event

```
cmdlog record --cwd=/tmp -- echo hi -> {"ts":"...","cwd":"/tmp","cmd":"echo hi"}
```

## Steps

1. Run `cmdlog record --cwd=/tmp -- echo hi` in isolated HOME.

## Context

- No pre-existing `~/.cmdlog` directory; record must create it.

```go
func Setup(t *testing.T, req *Request) error {
	req.RecordCWD = "/tmp"
	req.RecordCommand = "echo hi"
	return nil
}
```