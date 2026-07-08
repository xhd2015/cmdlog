# Scenario

**Feature**: bash integration status reports filesystem-derived state (read-only)

```
cmdlog integration bash --status -> installed | partial | not installed (no writes)
```

## Steps

1. Set `req.IntegrationAction = "status"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationAction = "status"
	return nil
}
```