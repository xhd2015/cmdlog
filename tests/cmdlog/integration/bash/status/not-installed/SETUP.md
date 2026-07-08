# Scenario

**Feature**: status reports not installed on empty HOME

```
# empty isolated HOME
cmdlog integration bash --status -> not installed, exit 1
```

## Steps

1. Start with empty isolated HOME (no bash.sh, no marker).

```go
func Setup(t *testing.T, req *Request) error {
	// Empty HOME: no script, no marker, no pre-install.
	req.PreExistingProfile = ""
	req.PreExistingBashSh = ""
	req.PreInstallIntegration = false
	return nil
}
```