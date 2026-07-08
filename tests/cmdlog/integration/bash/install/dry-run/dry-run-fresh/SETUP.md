# Scenario

**Feature**: install dry-run on fresh HOME shows would-write preview

```
# empty isolated HOME
cmdlog integration bash --install --dry-run -> preview stdout, no bash.sh or profile created
```

## Steps

1. Start with empty isolated HOME (no pre-existing profile or bash.sh).

```go
func Setup(t *testing.T, req *Request) error {
	// Fresh HOME: no pre-seeded profile or bash.sh (parent sets install + dry-run).
	req.PreExistingProfile = ""
	req.PreExistingBashSh = ""
	req.PreInstallIntegration = false
	return nil
}
```