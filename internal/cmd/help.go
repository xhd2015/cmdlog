package cmd

import (
	"fmt"
	"os"
)

const HelpText = `
Usage: cmdlog <command> [options]

Track terminal activity in ~/.cmdlog/events.jsonl.

Commands:
  record       append one command event
  today        timeline for today
  summary      stats for a date
  timeline     filter events by time range
  dirs         group commands by directory
  integration  shell integration (bash, zsh)

Examples:
  cmdlog record --cwd=/path -- git status
  cmdlog today
  cmdlog integration bash --install

Options:
  -h, --help   show this help
`

const recordHelp = `
Usage: cmdlog record --cwd=<dir> -- <command>

Append one command event to ~/.cmdlog/events.jsonl.

Options:
  --cwd <dir>   working directory for the command (required)
  -h, --help    show this help

Examples:
  cmdlog record --cwd=/tmp -- echo hello
`

const todayHelp = `
Usage: cmdlog today

Print a chronological timeline of today's commands, grouped by hour and directory.

Options:
  -h, --help   show this help
`

const summaryHelp = `
Usage: cmdlog summary [--date YYYY-MM-DD]

Print command count, directory count, and top directory for a date.

Options:
  --date <date>   filter by date (default: today)
  -h, --help      show this help

Examples:
  cmdlog summary
  cmdlog summary --date 2026-07-08
`

const timelineHelp = `
Usage: cmdlog timeline [--from RFC3339] [--to RFC3339]

Print a chronological timeline filtered by time range.

Options:
  --from <time>   include events at or after this time (RFC3339)
  --to <time>     include events at or before this time (RFC3339)
  -h, --help      show this help

Examples:
  cmdlog timeline --from 2026-07-08T09:00:00Z --to 2026-07-08T18:00:00Z
`

const dirsHelp = `
Usage: cmdlog dirs [--date YYYY-MM-DD]

Group commands by working directory for a date.

Options:
  --date <date>   filter by date (default: today)
  -h, --help      show this help

Examples:
  cmdlog dirs
  cmdlog dirs --date 2026-07-08
`

const integrationHelp = `
Usage: cmdlog integration <shell> <action> [options]

Install shell hooks that record commands to ~/.cmdlog/events.jsonl.

Shells:
  bash   bash integration via PROMPT_COMMAND
  zsh    not yet supported

Examples:
  cmdlog integration bash --install
  cmdlog integration bash --status

Options:
  -h, --help   show this help
`

const integrationBashHelp = `
Usage: cmdlog integration bash <action> [options]

Manage bash integration via ~/.cmdlog/integration/bash.sh.

Actions:
  --install      write bash.sh and add marker to ~/.bash_profile
  --uninstall    remove marker from ~/.bash_profile
  --status       report installed / partial / not installed

Options:
  --dry-run      preview changes without writing files (install, uninstall)
  -h, --help     show this help

Examples:
  cmdlog integration bash --install
  cmdlog integration bash --install --dry-run
  cmdlog integration bash --status
`

const integrationBashInstallHelp = `
Usage: cmdlog integration bash --install [--dry-run]

Install bash integration:
  - writes ~/.cmdlog/integration/bash.sh
  - appends idempotent marker block to ~/.bash_profile

Options:
  --dry-run    preview changes without writing files
  -h, --help   show this help
`

const integrationBashUninstallHelp = `
Usage: cmdlog integration bash --uninstall [--dry-run]

Remove cmdlog marker block from ~/.bash_profile.
Does not delete ~/.cmdlog/events.jsonl or integration/bash.sh.

Options:
  --dry-run    preview changes without writing files
  -h, --help   show this help
`

const integrationBashStatusHelp = `
Usage: cmdlog integration bash --status

Report bash integration state (read-only).

States:
  installed      script and profile marker both present (exit 0)
  partial        only one present (exit 1)
  not installed  neither present (exit 1)

Options:
  -h, --help   show this help
`

const integrationZshHelp = `
Usage: cmdlog integration zsh <action>

Zsh integration is not yet supported.

Actions:
  --install    not yet supported

Options:
  -h, --help   show this help
`

func printHelp(text string) {
	fmt.Print(text)
	os.Exit(0)
}