# Scenario

**Feature**: zsh integration is not yet implemented

```
cmdlog integration zsh --install -> error "not yet supported"
```

## Steps

1. Set `req.IntegrationShell = "zsh"` and `req.IntegrationAction = "install"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationShell = "zsh"
	req.IntegrationAction = "install"
	return nil
}
```