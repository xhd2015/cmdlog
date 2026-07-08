package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lessflags "github.com/xhd2015/less-flags"
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
	remain, err := lessflags.Help("-h,--help", integrationHelp).
		StopOnFirstArg().
		Parse(args)
	if err != nil {
		return err
	}

	if len(remain) == 0 {
		printHelp(integrationHelp)
	}

	shell := remain[0]
	switch shell {
	case "bash":
		return runIntegrationBash(remain[1:])
	case "zsh":
		return runIntegrationZsh(remain[1:])
	default:
		return fmt.Errorf("unknown shell %q", shell)
	}
}

func runIntegrationBash(args []string) error {
	action, dryRun, err := parseIntegrationShellArgs(args, integrationBashHelp, map[string]string{
		"--install":   integrationBashInstallHelp,
		"--uninstall": integrationBashUninstallHelp,
		"--status":    integrationBashStatusHelp,
	})
	if err != nil {
		return err
	}
	if action == "" {
		printHelp(integrationBashHelp)
	}
	if dryRun && action == "--status" {
		return fmt.Errorf("unknown integration action %q", "--dry-run")
	}

	switch action {
	case "--install":
		if dryRun {
			return installBashDryRun()
		}
		return installBash()
	case "--uninstall":
		if dryRun {
			return uninstallBashDryRun()
		}
		return uninstallBash()
	case "--status":
		statusBash()
		return nil
	default:
		return fmt.Errorf("unknown integration action %q", action)
	}
}

func runIntegrationZsh(args []string) error {
	action, _, err := parseIntegrationShellArgs(args, integrationZshHelp, map[string]string{
		"--install": integrationZshHelp,
	})
	if err != nil {
		return err
	}
	if action == "" {
		printHelp(integrationZshHelp)
	}
	if action == "--install" {
		fmt.Fprintln(os.Stderr, "zsh integration not yet supported")
		os.Exit(1)
	}
	return fmt.Errorf("unknown integration action %q", action)
}

func parseIntegrationShellArgs(args []string, shellHelp string, actionHelp map[string]string) (action string, dryRun bool, err error) {
	for _, arg := range args {
		switch arg {
		case "--install", "--uninstall", "--status":
			if action != "" {
				return "", false, fmt.Errorf("unknown integration action %q", arg)
			}
			action = arg
		case "--dry-run":
			dryRun = true
		case "-h", "--help":
			if text, ok := actionHelp[action]; ok && action != "" {
				printHelp(text)
			}
			printHelp(shellHelp)
		default:
			return "", false, fmt.Errorf("unknown integration action %q", arg)
		}
	}
	return action, dryRun, nil
}

func bashPaths() (home, scriptPath, profilePath string, err error) {
	home, err = os.UserHomeDir()
	if err != nil {
		return "", "", "", err
	}
	scriptPath = filepath.Join(home, ".cmdlog", "integration", "bash.sh")
	profilePath = filepath.Join(home, ".bash_profile")
	return home, scriptPath, profilePath, nil
}

func scriptPresent(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func markerPresent(profilePath string) bool {
	data, err := os.ReadFile(profilePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), markerBegin)
}

func installBashDryRun() error {
	_, scriptPath, profilePath, err := bashPaths()
	if err != nil {
		return err
	}

	if markerPresent(profilePath) {
		fmt.Println("cmdlog bash integration: already installed")
		fmt.Printf("profile: %s (marker present)\n", profilePath)
		if scriptPresent(scriptPath) {
			fmt.Printf("script: %s (exists)\n", scriptPath)
		} else {
			fmt.Printf("script: %s (absent)\n", scriptPath)
		}
		fmt.Println("no changes needed")
		return nil
	}

	fmt.Println("dry-run: would write ~/.cmdlog/integration/bash.sh")
	fmt.Println("dry-run: would append marker block to ~/.bash_profile")
	fmt.Println()
	fmt.Print(markerBlock)
	return nil
}

func uninstallBashDryRun() error {
	_, _, profilePath, err := bashPaths()
	if err != nil {
		return err
	}

	if !markerPresent(profilePath) {
		fmt.Println("cmdlog bash integration: already uninstalled")
		fmt.Printf("profile: %s (marker absent)\n", profilePath)
		fmt.Println("no changes needed")
		return nil
	}

	fmt.Println("dry-run: would remove marker block from ~/.bash_profile")
	fmt.Println()
	fmt.Print(markerBlock)
	return nil
}

func statusBash() {
	_, scriptPath, profilePath, err := bashPaths()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	scriptExists := scriptPresent(scriptPath)
	markerExists := markerPresent(profilePath)

	var state string
	exitCode := 1
	switch {
	case scriptExists && markerExists:
		state = "installed"
		exitCode = 0
	case !scriptExists && !markerExists:
		state = "not installed"
	default:
		state = "partial"
	}

	fmt.Printf("bash integration: %s\n", state)
	if scriptExists {
		fmt.Printf("script: %s (present)\n", scriptPath)
	} else {
		fmt.Printf("script: %s (absent)\n", scriptPath)
	}
	if markerExists {
		fmt.Printf("profile: %s (marker present)\n", profilePath)
	} else {
		fmt.Printf("profile: %s (marker absent)\n", profilePath)
	}
	fmt.Println()
	os.Exit(exitCode)
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