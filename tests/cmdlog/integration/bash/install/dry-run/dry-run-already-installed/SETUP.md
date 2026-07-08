# Scenario

**Feature**: install dry-run when already installed reports no changes

```
# HOME has marker + bash.sh
cmdlog integration bash --install --dry-run -> already installed message, files unchanged
```

## Steps

1. Pre-seed `.bash_profile` with cmdlog marker and `bash.sh` with known content.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingProfile = `# user config
export EDITOR=vim
` + cmdlogMarkerBlock()
	req.PreExistingBashSh = minimalBashSh()
	return nil
}
```