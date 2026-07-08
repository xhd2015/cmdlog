# Scenario

**Feature**: manual nested cmdlog record is not rejected

```
# user manually runs cmdlog record inside another record invocation
cmdlog record --cwd=/tmp -- cmdlog record --cwd=/x -- foo -> full command string logged
```

## Steps

1. Run `cmdlog record --cwd=/tmp -- cmdlog record --cwd=/x -- foo`.

## Context

- Go `record` subcommand must **not** use prefix-based rejection of `cmdlog record`.
- The entire command string after `--` is stored in `cmd`.

```go
func Setup(t *testing.T, req *Request) error {
	req.RecordCWD = "/tmp"
	req.RecordCommand = "cmdlog record --cwd=/x -- foo"
	return nil
}
```