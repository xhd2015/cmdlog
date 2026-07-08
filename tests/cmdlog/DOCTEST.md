# cmdlog — Doc-Style Test Tree

Test suite for the **cmdlog** standalone Go CLI: append-only command logging to
`~/.cmdlog/events.jsonl`, query subcommands (`today`, `summary`, `timeline`,
`dirs`), and shell integration (`integration bash --install/--uninstall`,
`integration zsh --install` stub).

All tests use **isolated temporary `HOME`** via `t.Setenv("HOME", t.TempDir())` —
never the real `~/.cmdlog` or `~/.bash_profile`.

## Version

0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **cmdlog CLI** — standalone Go binary with subcommands `record`, `today`,
  `summary`, `timeline`, `dirs`, and `integration`.
- **events.jsonl** — append-only JSONL store at `$HOME/.cmdlog/events.jsonl`;
  each line is `{"ts":"<UTC ISO8601>","cwd":"<path>","cmd":"<string>"}`.
- **bash integration script** — `$HOME/.cmdlog/integration/bash.sh` installed by
  `cmdlog integration bash --install`; hooks into `PROMPT_COMMAND` via
  `_cmdlog_capture` delta detection.
- **~/.bash_profile** — receives an idempotent marker block sourcing bash.sh on
  install; marker removed on uninstall (events preserved).
- **Test harness** — builds cmdlog once per session, sets fake HOME, seeds
  fixtures, runs CLI or bash hook simulation (`history -s` + `_cmdlog_capture`).

**Behaviors**

- `record --cwd=<dir> -- <command>` appends one Event; creates `~/.cmdlog` if
  missing; **never** rejects `cmdlog record` invocations (manual records allowed).
- Query commands read `events.jsonl`, skip malformed lines with stderr warning.
- `today` prints chronological timeline grouped by hour + cwd; empty log → friendly
  message, exit 0.
- `summary [--date YYYY-MM-DD]` prints command/directory counts and top directory.
- `timeline [--from RFC3339] [--to RFC3339]` filters events by time range.
- `dirs [--date YYYY-MM-DD]` groups commands by cwd for a date.
- **Stdout contract:** all user-facing stdout ends with trailing `\n`.
- Bash hook calls `cmdlog record` with `set +o history` so automation plumbing
  does not enter bash history; user commands and manual `cmdlog record` are logged.
- `integration bash --install/--uninstall --dry-run` previews would-write/would-remove
  actions without modifying profile, `bash.sh`, or `events.jsonl`.
- `integration bash --status` reports `installed` / `partial` / `not installed` from
  filesystem inspection only (read-only); exit 0 when installed, 1 otherwise.
- `integration zsh --install` returns non-zero with "not yet supported".

## Decision Tree

```
tests/cmdlog/                                    ROOT: Request{Subcommand, ...}
│                                                 Response{ExitCode, Stdout, Events, ...}
│                                                 Run() builds CLI, fake HOME, dispatches
│
├── record/                                      DECISION: Subcommand = record
│   ├── append-simple-command/                   LEAF: echo hi → one JSONL line
│   ├── append-manual-cmdlog-record/             LEAF: nested cmdlog record not rejected
│   └── create-storage-dir/                      LEAF: missing ~/.cmdlog → auto-created
│
├── today/                                       DECISION: Subcommand = today
│   ├── timeline-with-fixtures/                  LEAF: fixture events → hour+cwd timeline
│   └── empty-log-friendly-message/              LEAF: missing log → friendly msg, exit 0
│
├── summary/                                     DECISION: Subcommand = summary
│   └── stats-for-date/                          LEAF: N commands, M dirs, top path
│
├── timeline/                                    DECISION: Subcommand = timeline
│   └── filter-from-to/                          LEAF: --from/--to returns matching events
│
├── dirs/                                        DECISION: Subcommand = dirs
│   └── group-by-cwd/                            LEAF: groups commands by cwd for date
│
├── integration/                                 DECISION: Subcommand = integration
│   ├── bash/                                    DECISION: shell = bash
│   │   ├── install/                             DECISION: action = install
│   │   │   ├── first-install/                   LEAF: creates bash.sh + profile marker
│   │   │   ├── idempotent-second-install/       LEAF: second install → no duplicate marker
│   │   │   └── dry-run/                         DECISION: --dry-run on install
│   │   │       ├── dry-run-fresh/               LEAF: empty HOME → preview only, no writes
│   │   │       └── dry-run-already-installed/   LEAF: pre-seeded → no changes, exit 0
│   │   ├── uninstall/                           DECISION: action = uninstall
│   │   │   ├── remove-marker-preserve-events/   LEAF: marker gone, events.jsonl intact
│   │   │   └── dry-run/                         DECISION: --dry-run on uninstall
│   │   │       ├── dry-run-with-marker/         LEAF: marker present → removal preview, no writes
│   │   │       └── dry-run-already-uninstalled/ LEAF: no marker → no changes, exit 0
│   │   └── status/                              DECISION: action = status (read-only)
│   │       ├── installed/                       LEAF: script + marker → installed, exit 0
│   │       ├── not-installed/                   LEAF: empty HOME → not installed, exit 1
│   │       └── partial-script-only/             LEAF: bash.sh only → partial, exit 1
│   └── zsh/                                     DECISION: shell = zsh
│       └── not-supported/                       LEAF: error "not yet supported", exit ≠ 0
│
├── bash-hook/                                   DECISION: bash hook capture (HOME isolated)
│   ├── capture-user-command/                    LEAF: history -s git status → event
│   ├── capture-manual-cmdlog-record/            LEAF: manual cmdlog record → event
│   └── suppress-automation-plumbing/            LEAF: hook automation not logged as event
│
└── malformed-jsonl/                             DECISION: corrupt events.jsonl handling
    └── skip-bad-line-with-warning/              LEAF: bad line skipped, valid lines processed
```

## Test Index

| # | Leaf | Description |
|---|------|-------------|
| 1 | `record/append-simple-command/` | `cmdlog record --cwd=/tmp -- echo hi` → one valid JSONL line |
| 2 | `record/append-manual-cmdlog-record/` | Manual nested `cmdlog record` invocation is recorded, not rejected |
| 3 | `record/create-storage-dir/` | Missing `~/.cmdlog` directory created automatically |
| 4 | `today/timeline-with-fixtures/` | Today's fixture events → timeline grouped by hour + cwd |
| 5 | `today/empty-log-friendly-message/` | Empty/missing log → friendly message, exit 0 |
| 6 | `summary/stats-for-date/` | Summary shows command count, directory count, top directory |
| 7 | `timeline/filter-from-to/` | `--from` / `--to` filter returns only matching events |
| 8 | `dirs/group-by-cwd/` | Groups commands by cwd for given date |
| 9 | `integration/bash/install/first-install/` | Install creates `bash.sh` + marker in fake `.bash_profile` |
| 10 | `integration/bash/install/idempotent-second-install/` | Second install is no-op (no duplicate markers) |
| 11 | `integration/bash/uninstall/remove-marker-preserve-events/` | Uninstall removes marker; `events.jsonl` untouched |
| 12 | `integration/zsh/not-supported/` | `integration zsh --install` → "not yet supported", non-zero exit |
| 13 | `bash-hook/capture-user-command/` | `history -s "git status"` + `_cmdlog_capture` → event recorded |
| 14 | `bash-hook/capture-manual-cmdlog-record/` | Manual `cmdlog record` via history → event logged |
| 15 | `bash-hook/suppress-automation-plumbing/` | Hook's internal `cmdlog record` not logged as user event |
| 16 | `malformed-jsonl/skip-bad-line-with-warning/` | Malformed line skipped with stderr warning; valid lines processed |
| 17 | `integration/bash/install/dry-run/dry-run-fresh/` | `--install --dry-run` on empty HOME → preview stdout, no files written |
| 18 | `integration/bash/install/dry-run/dry-run-already-installed/` | `--install --dry-run` when installed → no changes, exit 0 |
| 19 | `integration/bash/uninstall/dry-run/dry-run-with-marker/` | `--uninstall --dry-run` with marker → removal preview, profile unchanged |
| 20 | `integration/bash/uninstall/dry-run/dry-run-already-uninstalled/` | `--uninstall --dry-run` without marker → no changes, exit 0 |
| 21 | `integration/bash/status/installed/` | `--status` after install → `installed`, exit 0 |
| 22 | `integration/bash/status/not-installed/` | `--status` on empty HOME → `not installed`, exit 1 |
| 23 | `integration/bash/status/partial-script-only/` | `--status` with bash.sh only → `partial`, exit 1 |

## How to Run

```sh
# Vet tree structure
doctest vet ./tests/cmdlog

# Run all tests (expect RED until implementation)
doctest test -v ./tests/cmdlog

# Run a single leaf
doctest test -v ./tests/cmdlog/record/append-simple-command
```

```go
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"
)

// Event mirrors the cmdlog JSONL schema under test.
type Event struct {
	TS  string `json:"ts"`
	CWD string `json:"cwd"`
	CMD string `json:"cmd"`
}

// Request drives one cmdlog invocation or bash-hook simulation.
// Defined only at root; descendants must not redefine.
type Request struct {
	Subcommand string // "record" | "today" | "summary" | "timeline" | "dirs" | "integration" | "bash-hook" | "malformed"

	// record
	RecordCWD     string
	RecordCommand string

	// query
	Date string // YYYY-MM-DD for summary/dirs; empty = today
	From string // RFC3339 for timeline
	To   string

	// integration
	IntegrationShell       string // "bash" | "zsh"
	IntegrationAction      string // "install" | "uninstall" | "status"
	DryRun                 bool   // --dry-run on install/uninstall
	RunTwice               bool   // second install for idempotency
	PreInstallIntegration  bool   // run real install before main integration command

	// bash-hook
	HookCommands []string // commands to history -s then _cmdlog_capture

	// fixtures / pre-state
	SeedEvents         []Event // written to events.jsonl before Run
	EventsFixtureFile  string  // path to JSONL fixture file
	EventsRawContent   string  // raw JSONL written verbatim (e.g. includes corrupt lines)
	PreExistingProfile string  // content for ~/.bash_profile before install/uninstall
	PreExistingBashSh  string  // content for ~/.cmdlog/integration/bash.sh before run
	PreExistingEvents  string  // raw JSONL written before uninstall test
}

// Response captures CLI outcome and filesystem state.
type Response struct {
	ExitCode      int
	Stdout        string
	Stderr        string
	StdoutSecond  string // populated when RunTwice
	Home          string
	EventsPath    string
	EventsContent string
	EventLines    []string
	ProfilePath   string
	ProfileContent string
	BashShPath    string
	BashShContent string
	MarkerCount   int
	Error         string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	bin := buildCmdlogOnce(t)

	if req.PreInstallIntegration {
		installReq := Request{
			Subcommand:        "integration",
			IntegrationShell:  req.IntegrationShell,
			IntegrationAction: "install",
		}
		if _, err := runCLI(t, &installReq, home, bin); err != nil {
			return nil, fmt.Errorf("pre-install integration: %w", err)
		}
	}

	if len(req.SeedEvents) > 0 || req.EventsFixtureFile != "" || req.EventsRawContent != "" {
		if err := writeEventsFixture(t, req, home); err != nil {
			return nil, err
		}
	}
	if req.PreExistingProfile != "" {
		if err := os.WriteFile(filepath.Join(home, ".bash_profile"), []byte(req.PreExistingProfile), 0644); err != nil {
			return nil, fmt.Errorf("write pre-existing profile: %w", err)
		}
	}
	if req.PreExistingBashSh != "" {
		integrationDir := filepath.Join(home, ".cmdlog", "integration")
		if err := os.MkdirAll(integrationDir, 0755); err != nil {
			return nil, fmt.Errorf("mkdir integration dir: %w", err)
		}
		if err := os.WriteFile(filepath.Join(integrationDir, "bash.sh"), []byte(req.PreExistingBashSh), 0644); err != nil {
			return nil, fmt.Errorf("write pre-existing bash.sh: %w", err)
		}
	}
	if req.PreExistingEvents != "" {
		cmdlogDir := filepath.Join(home, ".cmdlog")
		if err := os.MkdirAll(cmdlogDir, 0755); err != nil {
			return nil, fmt.Errorf("mkdir .cmdlog: %w", err)
		}
		if err := os.WriteFile(filepath.Join(cmdlogDir, "events.jsonl"), []byte(req.PreExistingEvents), 0644); err != nil {
			return nil, fmt.Errorf("write pre-existing events: %w", err)
		}
	}

	switch req.Subcommand {
	case "bash-hook":
		return runBashHook(t, req, home, bin)
	case "record", "today", "summary", "timeline", "dirs", "integration", "malformed":
		return runCLI(t, req, home, bin)
	default:
		return nil, fmt.Errorf("unknown subcommand %q", req.Subcommand)
	}
}

func runCLI(t *testing.T, req *Request, home, bin string) (*Response, error) {
	args := buildArgs(req)
	runOnce := func() (stdout, stderr string, code int) {
		cmd := exec.Command(bin, args...)
		cmd.Env = append(os.Environ(), "HOME="+home, "PATH="+filepath.Dir(bin)+":"+os.Getenv("PATH"))
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		code = 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				code = exitErr.ExitCode()
			} else {
				return "", "", -1
			}
		}
		return outBuf.String(), errBuf.String(), code
	}

	stdout, stderr, exitCode := runOnce()
	stdoutSecond := ""
	if req.RunTwice {
		stdoutSecond, _, _ = runOnce()
	}

	eventsPath := filepath.Join(home, ".cmdlog", "events.jsonl")
	eventsContent, eventLines := readEventsFile(eventsPath)
	profilePath := filepath.Join(home, ".bash_profile")
	profileContent, _ := readFileIfExists(profilePath)
	bashShPath := filepath.Join(home, ".cmdlog", "integration", "bash.sh")
	bashShContent, _ := readFileIfExists(bashShPath)

	return &Response{
		ExitCode:       exitCode,
		Stdout:         stdout,
		Stderr:         stderr,
		StdoutSecond:   stdoutSecond,
		Home:           home,
		EventsPath:     eventsPath,
		EventsContent:  eventsContent,
		EventLines:     eventLines,
		ProfilePath:    profilePath,
		ProfileContent: profileContent,
		BashShPath:     bashShPath,
		BashShContent:  bashShContent,
		MarkerCount:    countMarkers(profileContent),
	}, nil
}

func runBashHook(t *testing.T, req *Request, home, bin string) (*Response, error) {
	// Install integration into fake HOME so bash.sh and profile marker exist.
	installReq := Request{
		Subcommand:        "integration",
		IntegrationShell:  "bash",
		IntegrationAction: "install",
	}
	if _, err := runCLI(t, &installReq, home, bin); err != nil {
		return nil, fmt.Errorf("install bash integration: %w", err)
	}

	var script strings.Builder
	script.WriteString("set -euo pipefail\n")
	script.WriteString("source \"$HOME/.bash_profile\"\n")
	script.WriteString("cd \"$HOME\"\n")
	for _, cmd := range req.HookCommands {
		escaped := strings.ReplaceAll(cmd, "'", "'\\''")
		script.WriteString(fmt.Sprintf("history -s '%s'\n", escaped))
		script.WriteString("_cmdlog_capture\n")
	}

	cmd := exec.Command("bash", "--noprofile", "--norc", "-c", script.String())
	cmd.Env = append(os.Environ(),
		"HOME="+home,
		"PATH="+filepath.Dir(bin)+":"+os.Getenv("PATH"),
	)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("bash hook simulation failed: %w\nstderr: %s", err, errBuf.String())
	}

	eventsPath := filepath.Join(home, ".cmdlog", "events.jsonl")
	eventsContent, eventLines := readEventsFile(eventsPath)

	return &Response{
		ExitCode:      0,
		Stdout:        outBuf.String(),
		Stderr:        errBuf.String(),
		Home:          home,
		EventsPath:    eventsPath,
		EventsContent: eventsContent,
		EventLines:    eventLines,
	}, nil
}

func buildArgs(req *Request) []string {
	switch req.Subcommand {
	case "record":
		args := []string{"record", "--cwd=" + req.RecordCWD, "--"}
		if req.RecordCommand != "" {
			args = append(args, req.RecordCommand)
		}
		return args
	case "today":
		return []string{"today"}
	case "summary":
		args := []string{"summary"}
		if req.Date != "" {
			args = append(args, "--date", req.Date)
		}
		return args
	case "timeline":
		args := []string{"timeline"}
		if req.From != "" {
			args = append(args, "--from", req.From)
		}
		if req.To != "" {
			args = append(args, "--to", req.To)
		}
		return args
	case "dirs":
		args := []string{"dirs"}
		if req.Date != "" {
			args = append(args, "--date", req.Date)
		}
		return args
	case "integration":
		args := []string{"integration", req.IntegrationShell}
		if req.IntegrationAction == "status" {
			return append(args, "--status")
		}
		args = append(args, "--"+req.IntegrationAction)
		if req.DryRun {
			args = append(args, "--dry-run")
		}
		return args
	case "malformed":
		return []string{"today"}
	default:
		return nil
	}
}

func buildCmdlogOnce(t *testing.T) string {
	t.Helper()
	cacheDir := filepath.Join(os.TempDir(), "cmdlog-doctest-"+DOCTEST_SESSION_ID)
	lockPath := filepath.Join(cacheDir, "build.lock")
	readyPath := filepath.Join(cacheDir, "binaries.ready")
	binPath := filepath.Join(cacheDir, "cmdlog")

	withFileLock(t, lockPath, func() error {
		if fileExists(readyPath) && fileExists(binPath) {
			return nil
		}
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}
		moduleRoot := filepath.Join(DOCTEST_ROOT, "..", "..")
		// Try cmd/cmdlog first, then module root.
		for _, pkg := range []string{"./cmd/cmdlog", "."} {
			build := exec.Command("go", "build", "-o", binPath, pkg)
			build.Dir = moduleRoot
			if out, err := build.CombinedOutput(); err == nil {
				return os.WriteFile(readyPath, []byte("ok"), 0644)
			} else {
				_ = out
			}
		}
		return fmt.Errorf("go build cmdlog failed from %s", moduleRoot)
	})
	return binPath
}

func withFileLock(t *testing.T, lockPath string, fn func() error) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(lockPath), 0755); err != nil {
		t.Fatalf("mkdir lock dir: %v", err)
	}
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open lock: %v", err)
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		t.Fatalf("flock: %v", err)
	}
	defer syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	if err := fn(); err != nil {
		t.Fatalf("locked fn: %v", err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func writeEventsFixture(t *testing.T, req *Request, home string) error {
	t.Helper()
	cmdlogDir := filepath.Join(home, ".cmdlog")
	if err := os.MkdirAll(cmdlogDir, 0755); err != nil {
		return err
	}
	eventsPath := filepath.Join(cmdlogDir, "events.jsonl")

	if req.EventsRawContent != "" {
		content := req.EventsRawContent
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		return os.WriteFile(eventsPath, []byte(content), 0644)
	}

	var lines []string
	if req.EventsFixtureFile != "" {
		data, err := os.ReadFile(req.EventsFixtureFile)
		if err != nil {
			return fmt.Errorf("read fixture %s: %w", req.EventsFixtureFile, err)
		}
		for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
			if line != "" {
				lines = append(lines, line)
			}
		}
	}
	for _, ev := range req.SeedEvents {
		b, err := json.Marshal(ev)
		if err != nil {
			return err
		}
		lines = append(lines, string(b))
	}
	content := strings.Join(lines, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(eventsPath, []byte(content), 0644)
}

func readEventsFile(path string) (string, []string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil
	}
	content := string(data)
	var lines []string
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return content, lines
}

func readFileIfExists(path string) (string, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", false
	}
	return string(data), true
}

func countMarkers(profile string) int {
	return strings.Count(profile, "# === cmdlog integration begin ===")
}

// cmdlogMarkerBlock returns the standard profile marker block appended on install.
func cmdlogMarkerBlock() string {
	return `# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===
`
}

// minimalBashSh returns a tiny bash.sh stub for pre-seeding partial/installed states.
func minimalBashSh() string {
	return `#!/usr/bin/env bash
# cmdlog integration stub for doctest pre-seed
_cmdlog_capture() { :; }
`
}

// localTodayDate returns YYYY-MM-DD in local timezone for fixture seeding.
func localTodayDate() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

// utcTodayMorning returns RFC3339 UTC timestamp for today at given hour:minute.
func utcTodayAt(hour, minute int) string {
	loc := time.Local
	now := time.Now().In(loc)
	ts := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, loc).UTC()
	return ts.Format(time.RFC3339)
}

// parseEvents unmarshals all JSONL lines into Event structs.
func parseEvents(t *testing.T, lines []string) []Event {
	t.Helper()
	var events []Event
	for i, line := range lines {
		var ev Event
		if err := json.Unmarshal([]byte(line), &ev); err != nil {
			t.Fatalf("parse event line %d: %v\n%s", i+1, err, line)
		}
		events = append(events, ev)
	}
	return events
}

// eventCommands returns the cmd field from each parsed event.
func eventCommands(events []Event) []string {
	cmds := make([]string, len(events))
	for i, ev := range events {
		cmds[i] = ev.CMD
	}
	return cmds
}

// assertStdoutEndsWithNewline fails when user-facing stdout lacks trailing newline.
func assertStdoutEndsWithNewline(t *testing.T, stdout string) {
	t.Helper()
	if stdout == "" {
		return
	}
	if !strings.HasSuffix(stdout, "\n") {
		t.Fatalf("stdout must end with trailing newline; got %q", stdout)
	}
}

// assertHomeIsolated fails if path is outside the test HOME directory.
func assertHomeIsolated(t *testing.T, path, home string) {
	t.Helper()
	absPath, err := filepath.Abs(path)
	if err != nil {
		t.Fatalf("abs %q: %v", path, err)
	}
	homeAbs, err := filepath.Abs(home)
	if err != nil {
		t.Fatalf("abs home %q: %v", home, err)
	}
	if !strings.HasPrefix(absPath, homeAbs+string(filepath.Separator)) && absPath != homeAbs {
		t.Fatalf("path %q is outside isolated HOME %q", absPath, homeAbs)
	}
}
```