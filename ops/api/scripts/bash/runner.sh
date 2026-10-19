#!/usr/bin/env bash
# runner.sh - lifecycle manager for the Go API.
#
# Usage: runner.sh start | stop | status | watch [logfile] | help
#
# Design:
#   - `start` launches a small detached *supervisor* (this same script with the
#     internal `_supervise` command) in its own session when setsid exists,
#     otherwise under nohup. The supervisor runs `go run .`, appends its
#     stdout/stderr to the session log, records identity in the state file, and
#     writes the termination footer when the API ends - however it ends.
#   - Because the supervisor is detached, closing the menu or the log terminal
#     never touches the API.
#   - `stop` signals only the supervisor recorded in the state file (after
#     verifying its identity); the supervisor terminates the whole `go run`
#     process tree (leaves first), TERM then KILL after a grace period.

SELF="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)/$(basename -- "${BASH_SOURCE[0]}")"
# shellcheck source=common.sh
. "$(dirname -- "$SELF")/common.sh" || exit 1

umask 077
export API_PROJECT_ROOT="$ROOT_DIR"

STOP_GRACE="${API_STOP_GRACE:-8}"
case "$STOP_GRACE" in ''|*[!0-9]*) STOP_GRACE=8 ;; esac

# ------------------------------------------------------------------- lock --
# mkdir is atomic everywhere, so it is used instead of flock (not portable).
# Limitation: a lock left by a SIGKILLed runner is detected by PID liveness
# only; if it ever sticks, remove the lock.d directory by hand.
LOCK_DIR="${API_RUNTIME_DIR}/lock.d"
HAVE_LOCK=0

acquire_lock() {
    local tries=0 holder
    while ! mkdir "$LOCK_DIR" 2>/dev/null; do
        # The runtime dir is pruned when idle; recreate it if it vanished.
        [ -d "$API_RUNTIME_DIR" ] || { ensure_dirs || return 1; continue; }
        holder="$(cat "$LOCK_DIR/pid" 2>/dev/null)"
        if [ -n "$holder" ] && ! proc_alive "$holder"; then
            rm -rf -- "$LOCK_DIR"
            continue
        fi
        tries=$((tries + 1))
        if [ "$tries" -ge 25 ]; then
            err "Another start/stop is in progress (lock: $LOCK_DIR)."
            return 1
        fi
        sleep 0.2
    done
    printf '%s\n' "$$" >"$LOCK_DIR/pid"
    HAVE_LOCK=1
    trap release_lock EXIT
    return 0
}

release_lock() {
    if [ "$HAVE_LOCK" = 1 ]; then
        rm -rf -- "$LOCK_DIR"
        HAVE_LOCK=0
        prune_runtime_dir
    fi
}

# ------------------------------------------------------------- log helpers --
new_log_file() {
    local stamp date_dir base f n=0
    stamp="$(now_log_path)" || return 1
    date_dir="${stamp%/*}"
    base="${API_LOG_DIR}/${stamp}"
    mkdir -p -- "${API_LOG_DIR}/${date_dir}" || return 1
    f="${base}.log"
    # noclobber makes creation atomic: never overwrite an earlier session.
    while ! (set -C; : >"$f") 2>/dev/null; do
        n=$((n + 1))
        [ "$n" -gt 99 ] && return 1
        f="${base}_${n}.log"
    done
    printf '%s' "$f"
}

show_log_tail() {
    local f="$1"
    [ -f "$f" ] || return 0
    info "---- last lines of $f ----"
    tail -n 20 "$f" 2>/dev/null
    info "---------------------------"
}

# --------------------------------------------------------- live log window --
# Why this is careful about signals: a non-interactive shell starts every
# background (`&`) job with SIGINT/SIGQUIT *ignored*, and `nohup` ignores SIGHUP.
# Ignored signals survive exec, so the terminal - and the viewer inside it -
# would never react to Ctrl+C or to the window being closed. So the terminal is
# launched (a) under job control and (b) through a wrapper that resets those
# signals to their defaults (needed when *we* were started by something that
# ignores them, e.g. a menu opened from another script's `&`).
SIG_RESET=()
if env --default-signal=INT,HUP,QUIT true >/dev/null 2>&1; then
    SIG_RESET=(env --default-signal=INT,HUP,QUIT)
elif command -v perl >/dev/null 2>&1; then
    SIG_RESET=(perl -e '$SIG{$_}="DEFAULT" for qw(INT HUP QUIT); exec @ARGV or exit 127;' --)
fi

_open_monitor() {
    local log="$1" t pid rc title="API Console - Live Log"

    for t in xfce4-terminal gnome-terminal mate-terminal konsole x-terminal-emulator xterm; do
        command -v "$t" >/dev/null 2>&1 || continue
        case "$t" in
            xfce4-terminal|mate-terminal)
                "${SIG_RESET[@]}" "$t" --title="$title" -x bash "$SELF" watch --attached "$log" \
                    >/dev/null 2>&1 </dev/null & ;;
            gnome-terminal)
                "${SIG_RESET[@]}" "$t" --title="$title" -- bash "$SELF" watch --attached "$log" \
                    >/dev/null 2>&1 </dev/null & ;;
            konsole|x-terminal-emulator|xterm)
                "${SIG_RESET[@]}" "$t" -e bash "$SELF" watch --attached "$log" \
                    >/dev/null 2>&1 </dev/null & ;;
        esac
        pid=$!
        sleep 0.7
        if kill -0 "$pid" 2>/dev/null; then
            disown "$pid" 2>/dev/null || true
            return 0
        fi
        # Launcher already exited: fine if it succeeded (it forked/delegated).
        wait "$pid" 2>/dev/null
        rc=$?
        [ "$rc" -eq 0 ] && return 0
    done
    return 1
}

open_monitor() {
    local rc
    [ "$(uname -s 2>/dev/null)" = "Darwin" ] && return 1      # not supported yet
    [ -n "${DISPLAY:-}${WAYLAND_DISPLAY:-}" ] || return 1     # no GUI session
    set -m
    _open_monitor "$1"
    rc=$?
    set +m
    return $rc
}

# ------------------------------------------------------------------ start --
write_state() {
    local child="$1" tmp="${STATE_FILE}.$$"
    [ -d "$API_RUNTIME_DIR" ] || ensure_dirs
    {
        printf 'SUP_PID=%s\n' "$$"
        printf 'SUP_START=%s\n' "$(proc_start "$$")"
        printf 'CHILD_PID=%s\n' "$child"
        printf 'START_EPOCH=%s\n' "$START_EPOCH"
        printf 'START_HUMAN=%s\n' "$START_HUMAN"
        printf 'LOG_FILE=%s\n' "$LOG_FILE"
        printf 'ROOT=%s\n' "$ROOT_DIR"
    } >"$tmp" && mv -f -- "$tmp" "$STATE_FILE"
}

cmd_start() {
    local log epoch human launch_pid i
    command -v go >/dev/null 2>&1 ||
        { err "Go was not found in PATH. Install Go (or fix PATH) and try again."; return 1; }
    [ -d "$API_DIR" ] || { err "API directory not found: $API_DIR"; return 1; }
    ensure_dirs || return 1
    acquire_lock || return 1

    api_refresh_status
    if [ "$API_STATE" = running ]; then
        info "API is already running (PID $API_PID, uptime $(fmt_hms "$API_UPTIME"))."
        info "Log: $API_LOG"
        return 0
    fi

    log="$(new_log_file)" || { err "Could not create a log file under $API_LOG_DIR"; return 1; }
    epoch="$(now_epoch)"
    human="$(now_human)"
    {
        printf '[API Console] Started at: %s\n' "$human"
        printf '[API Console] Project root: %s\n' "$ROOT_DIR"
        printf '%s\n' "$LOG_SEP"
    } >>"$log"

    info "Starting API..."
    if command -v setsid >/dev/null 2>&1; then
        setsid bash "$SELF" _supervise "$log" "$epoch" "$human" </dev/null >/dev/null 2>&1 &
    else
        nohup bash "$SELF" _supervise "$log" "$epoch" "$human" </dev/null >/dev/null 2>&1 &
    fi
    launch_pid=$!
    disown "$launch_pid" 2>/dev/null || true

    # Wait (up to ~5s) for the supervisor to publish its state.
    i=0
    while [ "$i" -lt 25 ]; do
        api_refresh_status
        [ "$API_STATE" = running ] && break
        if ! proc_alive "$launch_pid"; then
            sleep 0.2
            api_refresh_status
            break
        fi
        sleep 0.2
        i=$((i + 1))
    done
    if [ "$API_STATE" != running ]; then
        err "API failed to start."
        show_log_tail "$log"
        return 1
    fi

    # Short grace period: catch an immediate crash before reporting Running.
    i=0
    while [ "$i" -lt 8 ]; do
        sleep 0.3
        api_refresh_status
        if [ "$API_STATE" != running ]; then
            err "API exited right after starting."
            show_log_tail "$log"
            return 1
        fi
        i=$((i + 1))
    done

    ok "API is running (PID $API_PID)."
    info "Log: $log"
    info "(go run may still be compiling - follow the log to confirm it is serving.)"

    if open_monitor "$log"; then
        info "Live log opened in a new terminal."
    else
        info "No GUI terminal available for the live log. Follow it with:"
        info "  bash \"$SELF\" watch"
    fi
    return 0
}

# ------------------------------------------------------------------- stop --
cmd_stop() {
    local sup pid log epoch i p elapsed
    ensure_dirs || return 1
    acquire_lock || return 1

    api_refresh_status
    if [ "$API_STATE" != running ]; then
        info "API is not running."
        return 0
    fi
    sup="$API_SUP_PID"; pid="$API_PID"; log="$API_LOG"; epoch="$API_START_EPOCH"

    info "Stopping API (PID $pid)..."
    kill -TERM "$sup" 2>/dev/null

    # Supervisor does TERM -> grace -> KILL on the tree itself; give it time.
    i=0
    while [ "$i" -lt $(((STOP_GRACE + 12) * 5)) ]; do
        proc_alive "$sup" || break
        sleep 0.2
        i=$((i + 1))
    done

    if proc_alive "$sup"; then
        warn "Supervisor did not exit in time; forcing termination."
        for p in $(descendants "$sup"); do kill -KILL "$p" 2>/dev/null; done
        kill -KILL "$sup" 2>/dev/null
        sleep 0.3
        elapsed=$(( $(now_epoch) - epoch ))
        if [ -f "$log" ] && ! grep -q '^Terminated at:' "$log" 2>/dev/null; then
            {
                printf '%s\n' "$LOG_SEP"
                printf '[API Console] Process force-killed by API Console.\n'
                printf 'Terminated at: %s\n' "$(now_human)"
                printf 'Elapsed time: %s seconds\n' "$elapsed"
            } >>"$log"
        fi
        rm -f -- "$STATE_FILE"
    fi

    elapsed=$(( $(now_epoch) - epoch ))
    ok "API stopped. It ran for $(fmt_hms "$elapsed")."
    info "Log: $log"
    return 0
}

# ----------------------------------------------------------------- status --
cmd_status() {
    api_refresh_status
    if [ "$API_STATE" = running ]; then
        info "Status: ${C_GREEN}Running${C_RESET}"
        info "PID: $API_PID"
        info "Uptime: $(fmt_hms "$API_UPTIME")"
        info "Log: $API_LOG"
    else
        info "Status: ${C_RED}Offline${C_RESET}"
    fi
    return 0
}

# ------------------------------------------------------------------ watch --
# The log window always shows white text, whatever the user's terminal theme
# (OSC 10 sets the default foreground, so ANSI resets in the API output also
# return to white; SGR 97 is the fallback for terminals without OSC 10).
TERM_TOUCHED=0
term_white() {
    if [ -t 1 ]; then
        printf '\033]10;#ffffff\007\033[97m'
        TERM_TOUCHED=1
    fi
}
term_restore() {
    if [ "$TERM_TOUCHED" = 1 ]; then
        printf '\033[0m\033]110\007'
        TERM_TOUCHED=0
    fi
}

register_monitor() {
    local tmp="${MONITOR_FILE}.$$"
    [ -d "$API_RUNTIME_DIR" ] || ensure_dirs
    printf 'MONITOR_PID=%s\n' "$$" >"$tmp" && mv -f -- "$tmp" "$MONITOR_FILE"
}

unregister_monitor() {
    local pid
    pid="$(monitor_read_pid)" || return 0
    if [ "$pid" = "$$" ]; then
        rm -f -- "$MONITOR_FILE"
        prune_runtime_dir
    fi
}

VIEW_LOG=''
VIEW_TAIL=''

# Ctrl+C in the log window = stop the API (it is the API's own terminal).
viewer_request_stop() {
    api_refresh_status
    if [ "$API_STATE" = running ] && [ "$API_LOG" = "$VIEW_LOG" ]; then
        printf '\n[API Console] Interrupt received - stopping API...\n'
        kill -TERM "$API_SUP_PID" 2>/dev/null
    fi
}

viewer_cleanup() {
    [ -n "$VIEW_TAIL" ] && kill "$VIEW_TAIL" 2>/dev/null
    term_restore
}

# Attached viewer: follows the log, and its lifetime is tied to the API's.
#   - window closed  -> this process dies; the supervisor notices the missing
#                       monitor and stops the API (and writes the log footer)
#   - Ctrl+C         -> API is stopped, then the window waits for Enter
#   - API ends       -> final log lines are shown, then the window waits for Enter
viewer_attached() {
    local log="$1" have_tty=0

    VIEW_LOG="$log"
    : >/dev/tty 2>/dev/null && have_tty=1
    register_monitor
    trap viewer_cleanup EXIT
    trap viewer_request_stop INT
    trap 'exit 0' HUP TERM    # leave the registration: supervisor stops the API

    tail -n +1 -F "$log" &
    VIEW_TAIL=$!

    while :; do
        api_refresh_status
        if [ "$API_STATE" != running ] || [ "$API_LOG" != "$log" ]; then
            break
        fi
        # Controlling terminal vanished (window closed but HUP was missed).
        if [ "$have_tty" = 1 ] && ! : >/dev/tty 2>/dev/null; then
            exit 0
        fi
        sleep 0.5
    done

    sleep 1                     # let tail print the termination footer
    kill "$VIEW_TAIL" 2>/dev/null
    VIEW_TAIL=''
    unregister_monitor
    trap 'exit 0' INT
    printf '\n[API Console] API terminated. Press [Enter] to close this window...'
    read -r _ || true
    return 0
}

cmd_watch() {
    local attached=0 log=""
    if [ "${1:-}" = "--attached" ]; then
        attached=1
        shift
    fi
    log="${1:-}"
    if [ -z "$log" ]; then
        api_refresh_status
        log="$API_LOG"
        [ -n "$log" ] || log="$(find "$API_LOG_DIR" -type f -name '*.log' -exec ls -t {} + 2>/dev/null | head -n 1)"
    fi
    if [ -z "$log" ] || [ ! -f "$log" ]; then
        err "No log file found in $API_LOG_DIR"
        return 1
    fi

    trap term_restore EXIT
    term_white
    printf 'API Console - live log: %s\n' "$log"
    if [ "$attached" = 1 ]; then
        printf '(Ctrl+C or closing this window stops the API.)\n'
    else
        printf '(Read-only view: Ctrl+C or closing this window does not stop the API.)\n'
    fi
    printf '%s\n' "$LOG_SEP"

    if [ "$attached" = 1 ]; then
        viewer_attached "$log"
        return 0
    fi
    tail -n +1 -F "$log"
}

# ------------------------------------------------------------- supervisor --
SUP_STOPPING=0

terminate_tree() {
    local targets p i alive
    targets="$(descendants $$)"
    for p in $targets; do kill -TERM "$p" 2>/dev/null; done
    i=0
    while [ "$i" -lt $((STOP_GRACE * 5)) ]; do
        alive=0
        for p in $targets; do
            if proc_alive "$p"; then alive=1; break; fi
        done
        [ "$alive" = 0 ] && return 0
        sleep 0.2
        i=$((i + 1))
    done
    for p in $targets; do
        proc_alive "$p" && kill -KILL "$p" 2>/dev/null
    done
    return 0
}

SUP_REASON=''

on_term() {
    [ "$SUP_STOPPING" = 1 ] && return
    SUP_STOPPING=1
    SUP_REASON='stop requested via API Console'
    terminate_tree
}

supervise() {
    local child rc elapsed how
    LOG_FILE="$1"; START_EPOCH="$2"; START_HUMAN="$3"

    # Survive terminal hangups / Ctrl-C aimed at the menu's process group.
    trap '' HUP INT QUIT
    trap on_term TERM

    cd -- "$API_DIR" || {
        printf '[API Console] Cannot enter %s\n' "$API_DIR" >>"$LOG_FILE"
        rc=1
    }

    if [ -z "${rc:-}" ]; then
        go run . >>"$LOG_FILE" 2>&1 </dev/null &
        child=$!
        write_state "$child"

        # Poll (instead of a blocking wait) so we can also notice that the
        # live-log terminal was closed; that stops the API like a normal Stop.
        while proc_alive "$child"; do
            if [ "$SUP_STOPPING" = 0 ] && monitor_gone; then
                SUP_STOPPING=1
                SUP_REASON='log terminal closed'
                terminate_tree
                break
            fi
            sleep 0.5
        done
        # Reap the child (it may still be running a moment after a signal).
        wait "$child"
        rc=$?
    fi

    elapsed=$(( $(now_epoch) - START_EPOCH ))
    how="${SUP_REASON:-exited on its own}"
    rm -f -- "$MONITOR_FILE"
    {
        printf '%s\n' "$LOG_SEP"
        printf '[API Console] Process exit code: %s (%s)\n' "$rc" "$how"
        printf 'Terminated at: %s\n' "$(now_human)"
        printf 'Elapsed time: %s seconds\n' "$elapsed"
    } >>"$LOG_FILE"

    # Only remove state that still describes us.
    if read_state && [ "$S_SUP_PID" = "$$" ]; then
        rm -f -- "$STATE_FILE"
    fi
    prune_runtime_dir
    exit 0
}

usage() {
    cat <<'EOF'
Usage: runner.sh <command>

  start            Start the API in the background (detached from this terminal)
  stop             Stop the API started by this tool
  status           Show whether the API is running
  watch [logfile]  Follow a log live (defaults to the current/latest session)
EOF
}

case "${1:-help}" in
    start)      cmd_start;  exit $? ;;
    stop)       cmd_stop;   exit $? ;;
    status)     cmd_status; exit $? ;;
    watch)      shift; cmd_watch "$@"; exit $? ;;
    _supervise) shift; supervise "$@" ;;
    help|-h|--help) usage ;;
    *)          usage >&2; exit 2 ;;
esac
