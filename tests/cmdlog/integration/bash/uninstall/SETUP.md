# Scenario

**Feature**: bash integration uninstall removes profile marker

```
cmdlog integration bash --uninstall -> marker removed from .bash_profile
```

## Steps

1. Set `req.IntegrationAction = "uninstall"`.
2. Pre-seed profile with marker and existing `events.jsonl`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationAction = "uninstall"
	return nil
}
```