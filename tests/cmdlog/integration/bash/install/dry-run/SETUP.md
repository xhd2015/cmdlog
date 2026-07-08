# Scenario

**Feature**: bash install dry-run previews changes without writing files

```
cmdlog integration bash --install --dry-run -> stdout preview only, no filesystem writes
```

## Steps

1. Set `req.IntegrationAction = "install"` and `req.DryRun = true`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationAction = "install"
	req.DryRun = true
	return nil
}
```