# Scenario

**Feature**: uninstall dry-run without marker reports no changes

```
# profile has no cmdlog marker
cmdlog integration bash --uninstall --dry-run -> already uninstalled message, no writes
```

## Steps

1. Start with empty isolated HOME (no marker in profile).

```go
func Setup(t *testing.T, req *Request) error {
	// No marker: ensure profile and bash.sh are not pre-seeded.
	req.PreExistingProfile = ""
	req.PreExistingBashSh = ""
	return nil
}
```