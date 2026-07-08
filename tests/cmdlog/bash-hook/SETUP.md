# Scenario

**Feature**: bash hook captures user commands via history delta detection

```
# install integration, simulate PROMPT_COMMAND capture
history -s "<cmd>" -> _cmdlog_capture -> cmdlog record (automation) -> events.jsonl
```

## Preconditions

- Bash available on PATH.
- Isolated `HOME`; integration installed before hook simulation.
- Pattern from `script/demo/home-isolation-demo.sh`: `bash --noprofile --norc -c` with
  `source "$HOME/.bash_profile"`, `history -s`, `_cmdlog_capture`.

## Steps

1. Set `req.Subcommand = "bash-hook"`.
2. `Run` installs bash integration, then runs hook simulation for each `HookCommands` entry.

## Context

- Hook's internal `cmdlog record` uses `set +o history` — automation plumbing must not
  appear as separate user events in `events.jsonl`.
- Manual `cmdlog record` typed by user **must** be logged.

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "bash-hook"
	return nil
}
```