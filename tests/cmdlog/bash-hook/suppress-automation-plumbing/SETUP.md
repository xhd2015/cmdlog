# Scenario

**Feature**: hook automation cmdlog record calls do not create plumbing events

```
# user runs git status and cmdlog today; hook's internal record calls are suppressed
history -s "git status" + history -s "cmdlog today" -> only 2 user events, no plumbing
```

## Steps

1. Simulate `git status` and `cmdlog today` via history + capture.
2. Verify events.jsonl has exactly 2 lines — no extra `cmdlog record --cwd=...` plumbing events.

## Context

- Mirrors `script/demo/home-isolation-demo.sh` verification: automation `cmdlog record`
  from hook must not appear as its own event (only user-visible commands).

```go
func Setup(t *testing.T, req *Request) error {
	req.HookCommands = []string{
		"git status",
		"cmdlog today",
	}
	return nil
}
```