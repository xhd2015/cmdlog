# Scenario

**Feature**: zsh install stub returns clear error

```
cmdlog integration zsh --install -> stderr "not yet supported", exit != 0
```

## Steps

1. Run `cmdlog integration zsh --install` in isolated HOME.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationShell = "zsh"
	req.IntegrationAction = "install"
	return nil
}
```