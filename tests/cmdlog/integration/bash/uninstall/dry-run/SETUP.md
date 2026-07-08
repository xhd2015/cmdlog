# Scenario

**Feature**: bash uninstall dry-run previews marker removal without writing files

```
cmdlog integration bash --uninstall --dry-run -> stdout preview only, no filesystem writes
```

## Steps

1. Set `req.IntegrationAction = "uninstall"` and `req.DryRun = true`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationAction = "uninstall"
	req.DryRun = true
	return nil
}
```