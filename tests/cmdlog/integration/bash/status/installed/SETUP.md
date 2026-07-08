# Scenario

**Feature**: status reports installed when script and marker are present

```
# real install into isolated HOME
cmdlog integration bash --status -> installed, exit 0
```

## Steps

1. Run real `cmdlog integration bash --install` before status (`req.PreInstallIntegration = true`).

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstallIntegration = true
	return nil
}
```