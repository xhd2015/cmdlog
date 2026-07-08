# cmdlog

track what were we doing in the terminal? what did I do all day?

Standalone Go CLI that appends terminal commands to `~/.cmdlog/events.jsonl` and provides query subcommands (`today`, `summary`, `timeline`, `dirs`). Bash integration hooks into `PROMPT_COMMAND` via `cmdlog integration bash --install`.

## Install

```sh
go install ./cmd/cmdlog
```

## Usage

```sh
cmdlog record --cwd=/path -- <command>
cmdlog today
cmdlog summary [--date YYYY-MM-DD]
cmdlog timeline [--from RFC3339] [--to RFC3339]
cmdlog dirs [--date YYYY-MM-DD]
cmdlog integration bash --install
cmdlog integration bash --install --dry-run
cmdlog integration bash --uninstall
cmdlog integration bash --uninstall --dry-run
cmdlog integration bash --status
```

`--install --dry-run` and `--uninstall --dry-run` preview changes without writing files.
`--status` reports whether bash integration is `installed`, `partial`, or `not installed`
(exit 0 only when fully installed).

## Storage

Events are stored as JSONL at `~/.cmdlog/events.jsonl`:

```json
{"ts":"2026-07-08T09:52:00Z","cwd":"/tmp","cmd":"git status"}
```