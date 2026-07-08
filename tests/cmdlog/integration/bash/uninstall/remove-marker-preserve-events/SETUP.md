# Scenario

**Feature**: uninstall removes marker but preserves events.jsonl

```
# profile has marker + events.jsonl has data
cmdlog integration bash --uninstall -> marker gone, events intact
```

## Steps

1. Pre-seed `.bash_profile` with cmdlog marker and `events.jsonl` with one event.
2. Run `cmdlog integration bash --uninstall`.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreExistingProfile = `# user config
export EDITOR=vim
# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===
`
	req.PreExistingEvents = `{"ts":"2026-07-08T09:00:00Z","cwd":"/tmp","cmd":"keep-me"}` + "\n"
	return nil
}
```