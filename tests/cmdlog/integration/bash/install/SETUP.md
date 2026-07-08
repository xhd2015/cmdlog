# Scenario

**Feature**: bash integration install writes hook script and profile marker

```
cmdlog integration bash --install -> ~/.cmdlog/integration/bash.sh + marker in .bash_profile
```

## Steps

1. Set `req.IntegrationAction = "install"`.
2. Start with empty isolated HOME (no pre-existing profile).

```go
func Setup(t *testing.T, req *Request) error {
	req.IntegrationAction = "install"
	return nil
}
```