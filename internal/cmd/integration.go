package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	markerBegin = "# === cmdlog integration begin ==="
	markerEnd   = "# === cmdlog integration end ==="
)

const bashIntegrationScript = `#!/usr/bin/env bash

CMDLOG_HIST_FIFO="$(mktemp -u /tmp/cmdlog-history-fifo.XXXXXX)"
trap 'rm -f "$CMDLOG_HIST_FIFO"' EXIT

CMDLOG_LAST_SUM=''
CMDLOG_LAST_CMD=''

_cmdlog_call_record() {
  local cmd=$1
  {
    set +o history
    CMDLOG_FROM_HOOK=1 cmdlog record --cwd="$PWD" -- "$cmd"
  }
}

_cmdlog_capture() {
  local saved_underscore="$_"
  : >"$CMDLOG_HIST_FIFO"
  local sum_cmd
  sum_cmd="$(HISTFILE=$CMDLOG_HIST_FIFO history -a && cat "$CMDLOG_HIST_FIFO")"
  local new_cmd=${sum_cmd#"$CMDLOG_LAST_SUM"}
  new_cmd=${new_cmd//$'\r'/}
  new_cmd=${new_cmd#$'\n'}
  CMDLOG_LAST_SUM=$sum_cmd
  if [[ -n $new_cmd && "$new_cmd" != "$CMDLOG_LAST_CMD" ]]; then
    CMDLOG_LAST_CMD=$new_cmd
    _cmdlog_call_record "$new_cmd"
  fi
  export _="$saved_underscore"
}

shopt -s histappend
PROMPT_COMMAND="_cmdlog_capture; ${PROMPT_COMMAND:-}"
`

const markerBlock = `# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===
`

func RunIntegration(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cmdlog integration <bash|zsh> --install|--uninstall")
	}

	shell := args[0]
	action := args[1]

	switch shell {
	case "bash":
		switch action {
		case "--install":
			return installBash()
		case "--uninstall":
			return uninstallBash()
		default:
			return fmt.Errorf("unknown integration action %q", action)
		}
	case "zsh":
		if action == "--install" {
			fmt.Fprintln(os.Stderr, "zsh integration not yet supported")
			os.Exit(1)
		}
		return fmt.Errorf("unknown integration action %q", action)
	default:
		return fmt.Errorf("unknown shell %q", shell)
	}
}

func installBash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	integrationDir := filepath.Join(home, ".cmdlog", "integration")
	if err := os.MkdirAll(integrationDir, 0o755); err != nil {
		return err
	}
	bashPath := filepath.Join(integrationDir, "bash.sh")
	if err := os.WriteFile(bashPath, []byte(bashIntegrationScript), 0o644); err != nil {
		return err
	}

	profilePath := filepath.Join(home, ".bash_profile")
	profile, err := os.ReadFile(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	content := string(profile)
	if strings.Contains(content, markerBegin) {
		return nil
	}

	var builder strings.Builder
	if len(content) > 0 && !strings.HasSuffix(content, "\n") {
		builder.WriteString(content)
		builder.WriteString("\n")
	} else {
		builder.WriteString(content)
	}
	builder.WriteString(markerBlock)
	return os.WriteFile(profilePath, []byte(builder.String()), 0o644)
}

func uninstallBash() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	profilePath := filepath.Join(home, ".bash_profile")
	data, err := os.ReadFile(profilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	lines := strings.Split(string(data), "\n")
	var out []string
	inBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == markerBegin {
			inBlock = true
			continue
		}
		if trimmed == markerEnd {
			inBlock = false
			continue
		}
		if inBlock {
			continue
		}
		out = append(out, line)
	}

	newContent := strings.Join(out, "\n")
	if len(data) > 0 && strings.HasSuffix(string(data), "\n") && newContent != "" {
		newContent += "\n"
	}
	return os.WriteFile(profilePath, []byte(newContent), 0o644)
}