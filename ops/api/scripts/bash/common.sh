#!/usr/bin/env bash
# common.sh - shared helpers for the API Console. Meant to be sourced, not run.
#
# Provides: root/path resolution, colours, messages, portable process helpers
# and the managed-process state reader (identity-verified, never eval'd).

if [ -n "${_API_CONSOLE_COMMON:-}" ]; then
    return 0
fi
_API_CONSOLE_COMMON=1

_COMMON_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

SEP_LINE='——————————————————————————————————'
LOG_SEP='--------------------------------------------------'

# ---------------------------------------------------------------- colours --
setup_colors() {
    if [ -t 1 ] && [ -z "${NO_COLOR:-}" ] && [ "${TERM:-dumb}" != "dumb" ]; then
        C_RED=$'\033[31m'
        C_GREEN=$'\033[32m'
        C_YELLOW=$'\033[33m'
        C_RESET=$'\033[0m'
    else
        C_RED=''
        C_GREEN=''
        C_YELLOW=''
        C_RESET=''
    fi
}
setup_colors

err()  { printf '%sError:%s %s\n' "$C_RED" "$C_RESET" "$*" >&2; }
warn() { printf '%sWarning:%s %s\n' "$C_YELLOW" "$C_RESET" "$*" >&2; }
info() { printf '%s\n' "$*"; }
ok()   { printf '%s%s%s\n' "$C_GREEN" "$*" "$C_RESET"; }
die()  { err "$*"; exit 1; }

# ------------------------------------------------------------------ paths --
looks_like_root() {
    [ -d "$1/api" ] && [ -d "$1/ops/api" ]
}

# Resolution order:
#   1. API_PROJECT_ROOT (set by the main start.sh) - must be valid, else error.
#   2. Derived from this script's own location (never from $PWD).
init_paths() {
    local candidate

    if [ -n "${API_PROJECT_ROOT:-}" ]; then
        [ -d "$API_PROJECT_ROOT" ] ||
            die "API_PROJECT_ROOT is set but is not a directory: $API_PROJECT_ROOT"
        candidate="$(cd -- "$API_PROJECT_ROOT" && pwd -P)" ||
            die "Cannot enter API_PROJECT_ROOT: $API_PROJECT_ROOT"
        looks_like_root "$candidate" ||
            die "API_PROJECT_ROOT does not look like the project root (needs api/ and ops/api/): $candidate"
    else
        # common.sh lives in <root>/ops/api/scripts/bash
        candidate="$(cd -- "$_COMMON_DIR/../../../.." && pwd -P)" ||
            die "Cannot derive project root from script location."
        looks_like_root "$candidate" ||
            die "Derived project root is invalid (needs api/ and ops/api/): $candidate"
    fi

    ROOT_DIR="$candidate"
    API_DIR="${ROOT_DIR}/api"
    API_LOG_DIR="${ROOT_DIR}/ops/api/logs"
    API_RUNTIME_DIR="${ROOT_DIR}/ops/api/.runtime"
    BIN_DIR="${ROOT_DIR}/bin"
    STATE_FILE="${API_RUNTIME_DIR}/api.state"
    MONITOR_FILE="${API_RUNTIME_DIR}/monitor"
}
init_paths

ensure_dirs() {
    local d
    for d in "$API_LOG_DIR" "$API_RUNTIME_DIR"; do
        if [ ! -d "$d" ]; then
            (umask 077; mkdir -p -- "$d") 2>/dev/null || {
                err "Cannot create directory: $d"
                return 1
            }
        fi
        if [ ! -w "$d" ]; then
            err "Directory is not writable: $d"
            return 1
        fi
    done
    chmod 700 "$API_RUNTIME_DIR" 2>/dev/null || true
    return 0
}

# ------------------------------------------------------------------- time --
system_date() { (unset TZ; LC_ALL=C date "$@"); }
now_human() { system_date '+%d:%m:%y %H:%M:%S'; }
now_epoch() { date '+%s'; }

now_log_path() {
    local parts month day year time suffix
    parts="$(system_date '+%b %d %Y %H:%M:%S')" || return 1
    read -r month day year time <<<"$parts"
    day=$((10#$day))

    case "$day" in
        11|12|13) suffix='th' ;;
        *)
            case $((day % 10)) in
                1) suffix='st' ;;
                2) suffix='nd' ;;
                3) suffix='rd' ;;
                *) suffix='th' ;;
            esac
            ;;
    esac

    printf '%s_%s%s_%s/%s' "$month" "$day" "$suffix" "$year" "$time"
}

fmt_hms() {
    local s="${1:-0}"
    [ "$s" -lt 0 ] 2>/dev/null && s=0
    printf '%02d:%02d:%02d' $((s / 3600)) $((s % 3600 / 60)) $((s % 60))
}

# -------------------------------------------------------- process helpers --
# Alive and not a zombie. Uses only kill -0 and ps -p (POSIX-ish).
proc_alive() {
    local pid="$1" st
    case "$pid" in ''|*[!0-9]*) return 1 ;; esac
    kill -0 "$pid" 2>/dev/null || return 1
    st="$(ps -p "$pid" -o stat= 2>/dev/null)"
    case "$st" in Z*) return 1 ;; esac
    return 0
}

# Process start time as text (identity token that survives PID reuse).
proc_start() {
    ps -p "$1" -o lstart= 2>/dev/null | awk '{$1=$1; print}'
}

proc_args() {
    ps -p "$1" -o args= 2>/dev/null
}

# All descendants of PID $1, deepest first (so leaves can be signalled before
# their parents). Works with GNU/procps and BSD/macOS ps.
descendants() {
    ps -A -o pid= -o ppid= -o comm= 2>/dev/null | awk -v root="$1" '
        $3 ~ /(^|\/)(ps|awk)$/ { next }
        { n++; id[n] = $1; par[$1] = $2 }
        END {
            want[root] = 1
            changed = 1
            while (changed) {
                changed = 0
                for (i = 1; i <= n; i++) {
                    c = id[i]
                    if (!(c in want) && (par[c] in want)) {
                        want[c] = 1
                        order[++m] = c
                        changed = 1
                    }
                }
            }
            for (i = m; i >= 1; i--) print order[i]
        }'
}

# ------------------------------------------------------------- .env reader --
# Reads ONE key from dotenv-style files without ever sourcing/eval'ing them.
# Handles `export KEY=...`, "double" / 'single' quotes, inline `# comments`
# on unquoted values, CRLF line endings. Last definition in a file wins; the
# first file (in argument order) that defines the key wins overall.
#   env_file_get KEY FILE...   -> prints the value, returns 1 if not found
env_file_get() {
    local key="$1" f line k v found have
    shift
    for f in "$@"; do
        [ -n "$f" ] && [ -r "$f" ] || continue
        found=''
        have=0
        while IFS= read -r line || [ -n "$line" ]; do
            line="${line%$'\r'}"
            line="${line#"${line%%[![:space:]]*}"}"
            case "$line" in ''|'#'*) continue ;; esac
            case "$line" in export[[:space:]]*) line="${line#export}"; line="${line#"${line%%[![:space:]]*}"}" ;; esac
            k="${line%%=*}"
            [ "$k" = "$line" ] && continue
            k="${k%"${k##*[![:space:]]}"}"
            [ "$k" = "$key" ] || continue
            v="${line#*=}"
            v="${v#"${v%%[![:space:]]*}"}"
            case "$v" in
                \"*) v="${v#\"}"; v="${v%%\"*}" ;;
                \'*) v="${v#\'}"; v="${v%%\'*}" ;;
                *)   v="${v%%[[:space:]]#*}"; v="${v%"${v##*[![:space:]]}"}" ;;
            esac
            found="$v"
            have=1
        done <"$f"
        if [ "$have" = 1 ]; then
            printf '%s' "$found"
            return 0
        fi
    done
    return 1
}

# ---------------------------------------------------------------- monitor --
# The live-log terminal registers itself here (MONITOR_PID=<pid>). While a
# registered monitor exists, the API's lifetime is tied to it: when that
# process disappears (terminal closed) the supervisor stops the API.
monitor_read_pid() {
    local line pid=''
    [ -f "$MONITOR_FILE" ] || return 1
    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in MONITOR_PID=*) pid="${line#MONITOR_PID=}" ;; esac
    done <"$MONITOR_FILE"
    case "$pid" in ''|*[!0-9]*) return 1 ;; esac
    printf '%s' "$pid"
}

# Is PID a live log-viewer of ours?
monitor_valid() {
    local args
    proc_alive "$1" || return 1
    args="$(proc_args "$1")"
    case "$args" in *runner.sh*watch*) return 0 ;; esac
    return 1
}

# True only if a monitor was registered and is now gone.
monitor_gone() {
    local pid
    pid="$(monitor_read_pid)" || return 1
    monitor_valid "$pid" && return 1
    return 0
}

# ------------------------------------------------------------------ state --
# State is plain KEY=VALUE text, parsed line by line - never sourced/eval'd.
S_SUP_PID=''; S_SUP_START=''; S_CHILD_PID=''; S_START_EPOCH=''
S_START_HUMAN=''; S_LOG_FILE=''

read_state() {
    local line key val
    S_SUP_PID=''; S_SUP_START=''; S_CHILD_PID=''; S_START_EPOCH=''
    S_START_HUMAN=''; S_LOG_FILE=''
    [ -f "$STATE_FILE" ] || return 1
    while IFS= read -r line || [ -n "$line" ]; do
        key="${line%%=*}"
        val="${line#*=}"
        case "$key" in
            SUP_PID)     S_SUP_PID="$val" ;;
            SUP_START)   S_SUP_START="$val" ;;
            CHILD_PID)   S_CHILD_PID="$val" ;;
            START_EPOCH) S_START_EPOCH="$val" ;;
            START_HUMAN) S_START_HUMAN="$val" ;;
            LOG_FILE)    S_LOG_FILE="$val" ;;
        esac
    done <"$STATE_FILE"
    case "$S_SUP_PID" in ''|*[!0-9]*) return 1 ;; esac
    case "$S_START_EPOCH" in ''|*[!0-9]*) return 1 ;; esac
    return 0
}

# Is the supervisor named in the state file really *our* supervisor?
# PID alive + command line is our runner + start-time matches the recorded one.
sup_valid() {
    local args cur
    proc_alive "$S_SUP_PID" || return 1
    args="$(proc_args "$S_SUP_PID")"
    case "$args" in *runner.sh*_supervise*) ;; *) return 1 ;; esac
    if [ -n "$S_SUP_START" ]; then
        cur="$(proc_start "$S_SUP_PID")"
        if [ -n "$cur" ] && [ "$cur" != "$S_SUP_START" ]; then
            return 1
        fi
    fi
    return 0
}

# Remove the runtime dir once nothing lives in it. rmdir only succeeds on an
# empty directory, so this can never delete state, a monitor file or a lock
# that is still in use.
prune_runtime_dir() {
    rmdir -- "$API_RUNTIME_DIR" 2>/dev/null
    return 0
}

clean_stale_state() {
    if [ -n "$S_LOG_FILE" ] && [ -f "$S_LOG_FILE" ] &&
        ! grep -q '^Terminated at:' "$S_LOG_FILE" 2>/dev/null; then
        {
            printf '%s\n' "$LOG_SEP"
            printf '[API Console] Supervisor ended unexpectedly; termination time unknown.\n'
        } >>"$S_LOG_FILE" 2>/dev/null || true
    fi
    rm -f -- "$STATE_FILE"
    prune_runtime_dir
}

# Re-checks reality (not just the stored state). Sets:
#   API_STATE (running|offline), API_PID, API_SUP_PID, API_LOG,
#   API_START_EPOCH, API_UPTIME
api_refresh_status() {
    API_STATE='offline'; API_PID=''; API_SUP_PID=''; API_LOG=''
    API_START_EPOCH=''; API_UPTIME=0

    read_state || {
        # Unparseable leftovers are treated as stale.
        if [ -f "$STATE_FILE" ]; then
            rm -f -- "$STATE_FILE"
            prune_runtime_dir
        fi
        return 0
    }

    if sup_valid; then
        API_STATE='running'
        API_PID="$S_CHILD_PID"
        API_SUP_PID="$S_SUP_PID"
        API_LOG="$S_LOG_FILE"
        API_START_EPOCH="$S_START_EPOCH"
        API_UPTIME=$(( $(now_epoch) - S_START_EPOCH ))
        [ "$API_UPTIME" -lt 0 ] && API_UPTIME=0
    else
        clean_stale_state
    fi
    return 0
}
