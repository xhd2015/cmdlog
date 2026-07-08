# Scenario

**Feature**: cmdlog record appends structured events to events.jsonl

```
# user or hook invokes record with cwd and command string
cmdlog record --cwd=<dir> -- <command> -> append Event -> events.jsonl
```

## Preconditions

- Isolated `HOME` with no pre-existing `~/.cmdlog` (unless leaf seeds state).
- `cmdlog` binary built and on PATH for the test process.

## Steps

1. Set `req.Subcommand = "record"`.
2. Configure `req.RecordCWD` and `req.RecordCommand` per leaf.

## Context

- `record` must never reject `cmdlog record` in the command string (manual invocations allowed).
- Storage path: `$HOME/.cmdlog/events.jsonl`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "record"
	return nil
}
```