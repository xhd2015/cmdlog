# Scenario

**Feature**: hook records normal user commands

```
history -s "git status" -> _cmdlog_capture -> event with cmd "git status"
```

## Steps

1. Simulate one user command: `git status`.

```go
func Setup(t *testing.T, req *Request) error {
	req.HookCommands = []string{"git status"}
	return nil
}
```