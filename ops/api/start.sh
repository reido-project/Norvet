#!/usr/bin/env bash
# API Console - interactive menu.
#
# Works both when launched from the main start.sh (which exports
# API_PROJECT_ROOT) and when run directly: in that case the project root is
# derived from this script's location, never from $PWD.
#
# The menu only *controls* the API via scripts/bash/runner.sh. The API itself
# runs detached, so exiting this menu never stops it.
#
# Input is line based: type the option, then press [Enter] to confirm. While
# you type, the status / uptime keep refreshing.

SELF_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
SCRIPTS_DIR="${SELF_DIR}/scripts/bash"

if [ ! -r "${SCRIPTS_DIR}/common.sh" ]; then
    printf 'Error: missing %s/common.sh\n' "$SCRIPTS_DIR" >&2
    exit 1
fi
# shellcheck source=scripts/bash/common.sh
. "${SCRIPTS_DIR}/common.sh" || exit 1

RUNNER="${SCRIPTS_DIR}/runner.sh"
COMPILE="${SCRIPTS_DIR}/compile.sh"
TEST="${SCRIPTS_DIR}/test.sh"
PUSH="${SCRIPTS_DIR}/push.sh"

IS_TTY=0
[ -t 1 ] && IS_TTY=1
NEED_CLEAR=1
REFRESH=1
INPUT=''

leave() {
    prune_runtime_dir
    exit "${1:-0}"
}
trap 'printf "\n"; leave 0' INT

# render [refresh]   refresh=1 re-checks the real process state first.
render() {
    local nl=$'\n' el='' home='' below=''
    local buf=''

    if [ "$IS_TTY" = 1 ]; then
        el=$'\033[K'; home=$'\033[H'; below=$'\033[J'
        if [ "$NEED_CLEAR" = 1 ]; then
            clear
            NEED_CLEAR=0
        fi
    fi

    [ "${1:-1}" = 1 ] && api_refresh_status

    buf+="${SEP_LINE}${el}${nl}"
    buf+="            API Console${el}${nl}"
    buf+="${SEP_LINE}${el}${nl}"
    if [ "$API_STATE" = running ]; then
        buf+="  Status: ${C_GREEN}Running${C_RESET}${el}${nl}"
        buf+="  PID: ${API_PID}${el}${nl}"
        buf+="  Uptime: $(fmt_hms "$API_UPTIME")${el}${nl}"
    else
        buf+="  Status: ${C_RED}Offline${C_RESET}${el}${nl}"
    fi
    buf+="${el}${nl}"
    buf+="  1. Start / Stop${el}${nl}"
    buf+="  2. Run Test Series${el}${nl}"
    buf+="  3. Compile${el}${nl}"
    buf+="  4. Push to VPS${el}${nl}"
    buf+="  0. Exit${el}${nl}"
    buf+="${SEP_LINE}${el}${nl}"

    printf '%s%s%s' "$home" "$buf" "$below"
    printf 'Choose an option: %s' "$INPUT"
}

pause() {
    read -r -p "Press [Enter] to continue..." _ || true
}

toggle_api() {
    api_refresh_status
    printf '\n'
    if [ "$API_STATE" = running ]; then
        bash "$RUNNER" stop
    else
        bash "$RUNNER" start
    fi
    pause
}

while :; do
    render "$REFRESH"
    REFRESH=0

    if [ "$IS_TTY" = 1 ]; then
        # Collect the line one key at a time so the screen can keep refreshing
        # (1s timeout) without losing what has been typed so far.
        IFS= read -rsn1 -t 1 ch
        rc=$?
        if [ "$rc" -gt 128 ]; then
            REFRESH=1
            continue
        fi
        [ "$rc" -ne 0 ] && break            # EOF (Ctrl+D)
        case "$ch" in
            '')                              # Enter: confirm
                choice="$INPUT"
                INPUT=''
                printf '\n'
                ;;
            $'\177'|$'\b')                   # Backspace
                INPUT="${INPUT%?}"
                continue
                ;;
            $'\033')                         # swallow arrow-key sequences
                IFS= read -rsn2 -t 0.05 _ || true
                continue
                ;;
            *)
                [ "${#INPUT}" -lt 8 ] && INPUT+="$ch"
                continue
                ;;
        esac
    else
        IFS= read -r choice || break
    fi

    case "$choice" in
        1)
            toggle_api
            NEED_CLEAR=1
            REFRESH=1
            ;;
        2)
            printf '\n'
            bash "$TEST"
            sleep 1.5
            NEED_CLEAR=1
            REFRESH=1
            ;;
        3)
            printf '\n'
            bash "$COMPILE"
            pause
            NEED_CLEAR=1
            REFRESH=1
            ;;
        4)
            printf '\n'
            bash "$PUSH"
            sleep 1.5
            NEED_CLEAR=1
            REFRESH=1
            ;;
        0)
            printf '\nExiting... (the API keeps running if it was started)\n'
            leave 0
            ;;
        '')
            REFRESH=1                        # bare Enter: just redraw
            ;;
        *)
            printf '\nInvalid option.\n'
            sleep 1
            NEED_CLEAR=1
            REFRESH=1
            ;;
    esac
done
leave 0
