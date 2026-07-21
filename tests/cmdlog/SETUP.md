# Scenario

**Feature**: cmdlog CLI logs and queries terminal commands in isolated HOME

```
# harness builds cmdlog, sets fake HOME, runs subcommand or bash hook
Test harness -> fake HOME (~/.cmdlog/events.jsonl) -> cmdlog CLI
cmdlog record -> append Event JSONL line
cmdlog today/summary/timeline/dirs -> read events.jsonl -> stdout report
cmdlog integration bash --install -> bash.sh + .bash_profile marker
bash hook -> history -s + _cmdlog_capture -> cmdlog record (automation suppressed)
```

## Preconditions

- Go module at `filepath.Join(d.DOCTEST_ROOT, "..", "..")` with installable `cmdlog` binary.
- `doctest vet` and `doctest test` available on PATH.
- **Isolation (mandatory):** every test sets `HOME` to `t.TempDir()` — never touches real
  `~/.cmdlog` or `~/.bash_profile`.
- Bash available for `bash-hook/` leaves (`bash --noprofile --norc -c`).
- Process-local binary built once under an in-memory mutex (`buildCmdlogOnce`); no session
  flock / free inject vars (one-process suite model).

## Steps

1. `buildCmdlogOnce(t, d)` compiles the cmdlog binary once per process (shared across leaves).
2. `Run` creates `home := t.TempDir()` and `t.Setenv("HOME", home)`.
3. Seed fixtures (`SeedEvents`, `EventsFixtureFile`, `PreExistingProfile`, `PreExistingEvents`)
   when the scenario requires pre-existing state.
4. Dispatch on `req.Subcommand`:
   - CLI subcommands via `runCLI` (exec cmdlog with constructed args).
   - `bash-hook` via `runBashHook` (install integration, then `history -s` + `_cmdlog_capture`).
5. Snapshot `events.jsonl`, `.bash_profile`, and `integration/bash.sh` into `Response`.
6. Return `(*Response, nil)` for leaf `Assert` functions.

## Context

- `d.DOCTEST_ROOT` refers to the `tests/cmdlog/` directory (`d *session.Doctest`).
- Event schema: `{"ts":"...Z","cwd":"...","cmd":"..."}` (UTC ISO8601 timestamp).
- Query tests with date-sensitive fixtures use `localTodayDate()` and `utcTodayAt()` helpers.
- Malformed JSONL leaf runs `cmdlog today` against a file with one bad line + one good line.
- Bash hook automation uses `set +o history` around internal `cmdlog record` calls.

```go
func Setup(t *testing.T, req *Request) error {
	_ = t
	_ = req
	return nil
}
```
