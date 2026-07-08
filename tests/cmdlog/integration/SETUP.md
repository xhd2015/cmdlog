# Scenario

**Feature**: cmdlog integration installs shell hooks into isolated HOME

```
cmdlog integration <shell> --install/--uninstall -> bash.sh + .bash_profile marker
```

## Preconditions

- Isolated `HOME`; never modifies real `~/.bash_profile`.

## Steps

1. Set `req.Subcommand = "integration"`.
2. Configure `IntegrationShell` and `IntegrationAction` per descendant.

## Context

- Bash install writes `$HOME/.cmdlog/integration/bash.sh` and appends idempotent marker to `~/.bash_profile`.
- Bash uninstall removes marker only; does not delete `events.jsonl`.
- Zsh install is a stub returning "not yet supported".

```go
func Setup(t *testing.T, req *Request) error {
	req.Subcommand = "integration"
	return nil
}
```