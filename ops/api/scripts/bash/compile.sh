#!/usr/bin/env bash
# compile.sh - build the Go API into ${ROOT_DIR}/bin/
#
# Build target (OS / architecture), highest priority first:
#   1. GOOS / GOARCH in the shell environment (one-off override, Go's own names)
#   2. OS / ARCH in <root>/.env               (what the VPS runs - the normal way)
#   3. the local machine, as reported by `go env`
# Each value is resolved independently.
#
# The only config file used is <root>/.env. It is only *read* for the two keys
# OS and ARCH - it is never sourced, so nothing else in it is executed or
# exposed.
#
# The binary is built to a temp file and moved into place only on success, so a
# failed build never replaces a previously good binary.

# shellcheck source=common.sh
. "$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)/common.sh" || exit 1

command -v go >/dev/null 2>&1 ||
    die "Go was not found in PATH. Install Go (or fix PATH) and try again."
[ -d "$API_DIR" ] || die "API directory not found: $API_DIR"
[ -f "$API_DIR/go.mod" ] || die "go.mod not found in $API_DIR"

ENV_FILE="${ROOT_DIR}/.env"

# resolve_target <label> <shell override value> <.env key> <go env name>
# Sets RESOLVED_VALUE and RESOLVED_SRC.
resolve_target() {
    local label="$1" envval="$2" key="$3" goname="$4" v
    if [ -n "$envval" ]; then
        RESOLVED_VALUE="$envval"
        RESOLVED_SRC="$label env var"
    elif v="$(env_file_get "$key" "$ENV_FILE")" && [ -n "$v" ]; then
        RESOLVED_VALUE="$v"
        RESOLVED_SRC=".env ($key)"
    else
        RESOLVED_VALUE="$(go env "$goname" 2>/dev/null)"
        RESOLVED_SRC="local Go"
    fi
    case "$RESOLVED_VALUE" in
        ''|*[!a-z0-9]*)
            die "Invalid $label value '$RESOLVED_VALUE' (from $RESOLVED_SRC). Expected lowercase letters/digits, e.g. linux or amd64."
            ;;
    esac
}

resolve_target GOOS   "${GOOS:-}"   OS   GOOS;   TARGET_OS="$RESOLVED_VALUE";   SRC_OS="$RESOLVED_SRC"
resolve_target GOARCH "${GOARCH:-}" ARCH GOARCH; TARGET_ARCH="$RESOLVED_VALUE"; SRC_ARCH="$RESOLVED_SRC"

# Reject combinations Go does not support (skipped if `go tool dist list`
# is unavailable).
SUPPORTED="$(go tool dist list 2>/dev/null)"
if [ -n "$SUPPORTED" ] && ! printf '%s\n' "$SUPPORTED" | grep -qx "${TARGET_OS}/${TARGET_ARCH}"; then
    die "Go does not support the target ${TARGET_OS}/${TARGET_ARCH} (OS from ${SRC_OS}, arch from ${SRC_ARCH}). See: go tool dist list"
fi

NAME="api-${TARGET_OS}-${TARGET_ARCH}"
[ "$TARGET_OS" = "windows" ] && NAME="${NAME}.exe"

mkdir -p -- "$BIN_DIR" || die "Cannot create output directory: $BIN_DIR"

OUT="${BIN_DIR}/${NAME}"
TMP="${BIN_DIR}/.${NAME}.tmp.$$"
trap 'rm -f -- "$TMP"' EXIT

HOST_OS="$(go env GOHOSTOS 2>/dev/null)"
HOST_ARCH="$(go env GOHOSTARCH 2>/dev/null)"

printf 'Source : %s\n' "$API_DIR"
printf 'Target : %s/%s   (OS: %s, arch: %s)\n' "$TARGET_OS" "$TARGET_ARCH" "$SRC_OS" "$SRC_ARCH"
if [ -n "$HOST_OS" ] && [ -n "$HOST_ARCH" ] &&
    { [ "$HOST_OS" != "$TARGET_OS" ] || [ "$HOST_ARCH" != "$TARGET_ARCH" ]; }; then
    printf 'Note   : cross-compiling from %s/%s\n' "$HOST_OS" "$HOST_ARCH"
fi
printf 'Output : %s\n\n' "$OUT"

if (cd -- "$API_DIR" && GOOS="$TARGET_OS" GOARCH="$TARGET_ARCH" go build -o "$TMP" .); then
    mv -f -- "$TMP" "$OUT" || die "Could not move the new binary into place: $OUT"
    ok "Build succeeded: $OUT"
    exit 0
else
    rc=$?
    [ "$rc" -ne 0 ] || rc=1
    err "Build failed (exit code $rc). Any previous binary was left untouched."
    exit "$rc"
fi
