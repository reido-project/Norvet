#!/usr/bin/env bash
# common.sh - shared helpers for the Prisma Console. Meant to be sourced.
#
# Self-contained on purpose: it does not depend on the API console's helpers.

if [ -n "${_PRISMA_CONSOLE_COMMON:-}" ]; then
    return 0
fi
_PRISMA_CONSOLE_COMMON=1

_COMMON_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"

# ---------------------------------------------------------------- colours --
if [ -t 1 ] && [ -z "${NO_COLOR:-}" ] && [ "${TERM:-dumb}" != "dumb" ]; then
    C_RED=$'\033[31m'; C_GREEN=$'\033[32m'; C_YELLOW=$'\033[33m'; C_RESET=$'\033[0m'
else
    C_RED=''; C_GREEN=''; C_YELLOW=''; C_RESET=''
fi

SEP_LINE='——————————————————————————————————'

err()  { printf '%sError:%s %s\n' "$C_RED" "$C_RESET" "$*" >&2; }
warn() { printf '%sWarning:%s %s\n' "$C_YELLOW" "$C_RESET" "$*" >&2; }
ok()   { printf '%s%s%s\n' "$C_GREEN" "$*" "$C_RESET"; }

pause() { read -r -p "Press [Enter] to continue..." _ || true; }

# ------------------------------------------------------------------ paths --
# Root resolution order:
#   1. PROJECT_ROOT, exported by the main start.sh (must be valid, else error)
#   2. derived from this file's location (never from $PWD)
init_paths() {
    local candidate

    if [ -n "${PROJECT_ROOT:-}" ]; then
        [ -d "$PROJECT_ROOT" ] || {
            err "PROJECT_ROOT is set but is not a directory: $PROJECT_ROOT"; return 1; }
        candidate="$(cd -- "$PROJECT_ROOT" && pwd -P)" || return 1
    else
        # this file lives in <root>/ops/prisma/scripts/bash
        candidate="$(cd -- "$_COMMON_DIR/../../../.." && pwd -P)" || return 1
    fi
    [ -d "$candidate/ops/prisma" ] || {
        err "Not a valid project root (missing ops/prisma): $candidate"; return 1; }

    ROOT_DIR="$candidate"
}

# Where `npx prisma` runs: always <root>/external/nodejs (package.json,
# node_modules and the prisma/ folder live there).
resolve_prisma_dir() {
    PRISMA_DIR="${ROOT_DIR}/external/nodejs"
    [ -d "$PRISMA_DIR" ] || {
        err "Prisma directory does not exist: $PRISMA_DIR"; return 1; }
    PRISMA_DIR="$(cd -- "$PRISMA_DIR" && pwd -P)"
}

# ----------------------------------------------------------------- prisma --
# Finds the locally installed Prisma CLI and its version. Plain `npx prisma`
# is never run without this check: when the CLI is not installed, npx would
# offer to download "latest" - possibly a different major version than the
# project uses.
#   sets PRISMA_VER (empty when not usable) and PRISMA_PROBLEM (reason).
prisma_detect() {
    local bin="${PRISMA_DIR}/node_modules/.bin/prisma" out
    PRISMA_VER=''; PRISMA_PROBLEM=''

    command -v npx >/dev/null 2>&1 || {
        PRISMA_PROBLEM="npx was not found in PATH. Install Node.js (which includes npx)."; return 1; }
    [ -x "$bin" ] || {
        PRISMA_PROBLEM="Prisma CLI not found in ${PRISMA_DIR}/node_modules/.bin. Run 'npm install' there."
        return 1; }

    out="$(cd -- "$PRISMA_DIR" && "$bin" --version </dev/null 2>/dev/null)" || true
    PRISMA_VER="$(printf '%s\n' "$out" | awk -F: '$1 ~ /^prisma[[:space:]]*$/ { gsub(/[[:space:]]/, "", $2); print $2; exit }')"
    if [ -z "$PRISMA_VER" ]; then
        PRISMA_PROBLEM="Could not read the Prisma version (is Node.js working?). Try: cd '$PRISMA_DIR' && npx prisma --version"
        return 1
    fi
    return 0
}

# Initialise everything a Prisma menu needs. Exits the script on failure to
# resolve paths; a missing CLI is not fatal (the menu still opens and explains).
prisma_init() {
    init_paths || exit 1
    resolve_prisma_dir || exit 1
    if [ -z "${PRISMA_VER:-}" ]; then
        prisma_detect || true
    else
        PRISMA_PROBLEM=''
    fi
}

# run_prisma <args...>   runs `npx prisma <args>` in PRISMA_DIR, interactive.
# Ctrl+C interrupts the Prisma command only, then control returns to the menu.
#
# Exit status 130 is what Prisma itself returns when it stops without doing
# anything: a prompt was declined or cancelled, or (seen with `migrate dev`
# after drift) it needs a database reset but does not ask for confirmation.
# That is not the same as being interrupted, so the two are reported
# differently.
SEEN_INT=0
run_prisma() {
    local rc prev_trap
    if [ -z "$PRISMA_VER" ]; then
        err "${PRISMA_PROBLEM:-Prisma is not available.}"
        return 1
    fi
    printf '\n> npx prisma %s\n\n' "$*"
    SEEN_INT=0
    prev_trap="$(trap -p INT)"
    trap 'SEEN_INT=1' INT
    (cd -- "$PRISMA_DIR" && npx prisma "$@")
    rc=$?
    if [ -n "$prev_trap" ]; then eval "$prev_trap"; else trap - INT; fi
    printf '\n'
    if [ "$rc" -eq 0 ]; then
        ok "Done."
    elif [ "$SEEN_INT" = 1 ]; then
        warn "Interrupted."
    elif [ "$rc" -eq 130 ]; then
        warn "Prisma stopped without making changes (exit status 130)."
        printf '  Prisma does this when a prompt is declined or cancelled, or when\n'
        printf '  it needs to reset the database but cannot ask for confirmation.\n'
        printf '  If it said "We need to reset the ... database", use Reset in the\n'
        printf '  Migration menu: it shows the target and asks you to confirm.\n'
    else
        err "Prisma exited with status $rc."
    fi
    return "$rc"
}
