# Scenario

**Feature**: hook records manual cmdlog record invocations

```
history -s "cmdlog record --cwd=/tmp -- hello-manual" -> _cmdlog_capture -> event logged
```

## Steps

1. Simulate user manually running `cmdlog record --cwd=/tmp -- hello-manual`.

```go
func Setup(t *testing.T, req *Request) error {
	req.HookCommands = []string{"cmdlog record --cwd=/tmp -- hello-manual"}
	return nil
}
```