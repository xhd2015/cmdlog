# Scenario

**Feature**: second bash install is a no-op

```
cmdlog integration bash --install (twice) -> single marker block unchanged
```

## Steps

1. Run `cmdlog integration bash --install` twice (`req.RunTwice = true`).

```go
func Setup(t *testing.T, req *Request) error {
	req.RunTwice = true
	return nil
}
```