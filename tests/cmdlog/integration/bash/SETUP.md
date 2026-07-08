# Scenario

**Feature**: bash shell integration install and uninstall

```
cmdlog integration bash --install -> bash.sh + profile marker
cmdlog integration bash --uninstall -> remove marker, preserve events
```

## Steps

1. Set `req.IntegrationShell = "bash"`.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationShell = "bash"
	return nil
}
```