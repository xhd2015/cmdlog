#!/usr/bin/env bash
# Demo: isolate cmdlog bash integration tests with a temporary HOME.
#
# Shows two reliable patterns (no real ~/.cmdlog or ~/.bash_profile touched):
#   A) file/install isolation — writes only under $FAKE_HOME
#   B) hook logic isolation — bash -c + history -s + _cmdlog_capture (no PTY needed)
#
# Usage: bash script/demo/home-isolation-demo.sh

set -euo pipefail

REAL_HOME="${HOME:?}"
WORKDIR="$(mktemp -d /tmp/cmdlog-home-isolation.XXXXXX)"
FAKE_HOME="$WORKDIR/home"
FAKE_BIN="$WORKDIR/bin"

cleanup() {
  rm -rf "$WORKDIR"
}
trap cleanup EXIT

log() {
  printf '==> %s\n' "$*"
}

pass() {
  printf 'PASS: %s\n' "$*"
}

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

assert_eq() {
  local desc=$1 expected=$2 actual=$3
  if [[ "$actual" == "$expected" ]]; then
    pass "$desc"
  else
    fail "$desc (expected '$expected', got '$actual')"
  fi
}

assert_file_exists() {
  local path=$1
  [[ -f "$path" ]] || fail "missing file: $path"
  pass "file exists: $path"
}

assert_not_exists() {
  local path=$1
  [[ ! -e "$path" ]] || fail "should not exist: $path"
  pass "not created in real home: $path"
}

snapshot_real_home() {
  REAL_CMDLOG_EVENTS="$REAL_HOME/.cmdlog/events.jsonl"
  if [[ -f "$REAL_CMDLOG_EVENTS" ]]; then
    REAL_EVENTS_BEFORE="$(wc -l < "$REAL_CMDLOG_EVENTS" | tr -d ' ')"
  else
    REAL_EVENTS_BEFORE="0"
  fi
  REAL_PROFILE_BEFORE="$(md5 -q "$REAL_HOME/.bash_profile" 2>/dev/null || echo missing)"
}

verify_real_home_untouched() {
  if [[ -f "$REAL_CMDLOG_EVENTS" ]]; then
    local after
    after="$(wc -l < "$REAL_CMDLOG_EVENTS" | tr -d ' ')"
    assert_eq "real ~/.cmdlog/events.jsonl line count unchanged" "$REAL_EVENTS_BEFORE" "$after"
  else
    assert_not_exists "$REAL_CMDLOG_EVENTS"
  fi

  local profile_after
  profile_after="$(md5 -q "$REAL_HOME/.bash_profile" 2>/dev/null || echo missing)"
  assert_eq "real ~/.bash_profile unchanged" "$REAL_PROFILE_BEFORE" "$profile_after"
}

setup_fake_home() {
  log "create isolated HOME at $FAKE_HOME"
  mkdir -p "$FAKE_HOME/.cmdlog/integration" "$FAKE_BIN"

  cat >"$FAKE_BIN/cmdlog" <<'STUB'
#!/usr/bin/env bash
set -euo pipefail

cmd=${1:-}
shift || true

case "$cmd" in
  record)
    cwd=""
    while [[ $# -gt 0 ]]; do
      case "$1" in
        --cwd=*)
          cwd="${1#--cwd=}"
          shift
          ;;
        --)
          shift
          break
          ;;
        *)
          echo "stub cmdlog: unknown arg: $1" >&2
          exit 2
          ;;
      esac
    done
    command_text="$*"
    mkdir -p "$HOME/.cmdlog"
    ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
    cwd_json="$(printf '%s' "$cwd" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"
    cmd_json="$(printf '%s' "$command_text" | python3 -c 'import json,sys; print(json.dumps(sys.stdin.read()))')"
    printf '{"ts":"%s","cwd":%s,"cmd":%s}\n' "$ts" "$cwd_json" "$cmd_json" >>"$HOME/.cmdlog/events.jsonl"
    ;;
  *)
    echo "stub cmdlog: unsupported command: $cmd" >&2
    exit 2
    ;;
esac
STUB
  chmod +x "$FAKE_BIN/cmdlog"

  cat >"$FAKE_HOME/.cmdlog/integration/bash.sh" <<'INTEGRATION'
#!/usr/bin/env bash

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
INTEGRATION

  cat >"$FAKE_HOME/.bash_profile" <<'PROFILE'
export HISTFILE="$HOME/.bash_history"
export HISTSIZE=2000
export HISTFILESIZE=4000
shopt -s histappend
# === cmdlog integration begin ===
[[ -f "$HOME/.cmdlog/integration/bash.sh" ]] && source "$HOME/.cmdlog/integration/bash.sh"
# === cmdlog integration end ===
PROFILE
}

demo_install_isolation() {
  log "A) simulate integration install into fake HOME only"
  local marker='# === cmdlog integration begin ==='
  if ! grep -qF "$marker" "$FAKE_HOME/.bash_profile"; then
    fail "fake profile missing cmdlog marker"
  fi
  pass "fake .bash_profile contains cmdlog marker"

  local install_touch="$FAKE_HOME/.cmdlog/integration/.installed"
  date -u +%Y-%m-%dT%H:%M:%SZ >"$install_touch"
  assert_file_exists "$install_touch"
}

demo_hook_via_bash_c() {
  log "B) exercise hook with HOME=$FAKE_HOME bash -c (history -s + _cmdlog_capture)"
  HOME="$FAKE_HOME" PATH="$FAKE_BIN:$PATH" bash --noprofile --norc -c '
    set -euo pipefail
    source "$HOME/.bash_profile"
    cd "$HOME"

    # Simulate user commands entering bash history, then PROMPT_COMMAND capture.
    history -s "git status"
    _cmdlog_capture

    history -s "cmdlog record --cwd=/tmp -- hello-manual"
    _cmdlog_capture

    history -s "cmdlog today"
    _cmdlog_capture
  '
}

verify_fake_home_results() {
  local events="$FAKE_HOME/.cmdlog/events.jsonl"
  assert_file_exists "$events"

  local line_count
  line_count="$(wc -l < "$events" | tr -d ' ')"
  assert_eq "fake events.jsonl line count" "3" "$line_count"

  grep -q '"cmd":"git status"' "$events" || fail 'missing git status event'
  pass 'captured user command: git status'

  grep -q 'hello-manual' "$events" || fail 'missing manual cmdlog record event'
  pass 'captured manual cmdlog record'

  grep -q '"cmd":"cmdlog today"' "$events" || fail 'missing cmdlog today event'
  pass 'captured cmdlog today'

  # Automation plumbing should not appear as its own event.
  if grep -q '"cmd":"cmdlog record --cwd=' "$events"; then
    local plumbing
    plumbing="$(grep -c '"cmd":"cmdlog record --cwd=' "$events" || true)"
    if [[ "$plumbing" -eq 1 ]]; then
      pass 'only manual cmdlog record appears once (automation calls not logged)'
    else
      fail "unexpected plumbing cmdlog record events: $plumbing"
    fi
  else
    pass 'automation cmdlog record calls not present in events'
  fi

  log "events.jsonl contents:"
  sed 's/^/    /' "$events"
}

main() {
  log "real HOME=$REAL_HOME"
  log "temp workdir=$WORKDIR"
  snapshot_real_home
  setup_fake_home
  demo_install_isolation
  demo_hook_via_bash_c
  verify_fake_home_results
  verify_real_home_untouched
  log "demo complete — HOME isolation approach looks good"
  log "recommended Go test pattern: t.Setenv(\"HOME\", t.TempDir()) + exec.Command(\"bash\", ...)"
}

main "$@"