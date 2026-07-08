# Scenario

**Feature**: first bash install creates integration files

```
# fresh HOME
cmdlog integration bash --install -> bash.sh exists + single profile marker
```

## Steps

1. Run first install on empty HOME.

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationShell = "bash"
	req.IntegrationAction = "install"
	req.RunTwice = false
	return nil
}
```