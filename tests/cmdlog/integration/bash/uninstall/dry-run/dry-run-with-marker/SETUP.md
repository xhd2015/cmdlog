# Scenario

**Feature**: uninstall dry-run with marker shows removal preview

```
# profile has cmdlog marker
cmdlog integration bash --uninstall --dry-run -> removal preview, profile unchanged
```

## Steps

1. Pre-seed `.bash_profile` with cmdlog marker block.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingProfile = `# user config
export EDITOR=vim
` + cmdlogMarkerBlock()
	return nil
}
```