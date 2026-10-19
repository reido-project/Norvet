#!/usr/bin/env bash
# Prisma Migration Console. Opened from ops/prisma/start.sh in the same
# terminal; can also be run directly.

SCRIPTS_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
# shellcheck source=common.sh
. "${SCRIPTS_DIR}/common.sh" || exit 1

prisma_init

trap 'printf "\n"; exit 0' INT

show_menu() {
    clear
    printf '%s\n' \
        "$SEP_LINE" \
        "    Prisma Migration Console" \
        "$SEP_LINE" \
        "  1. Migrate" \
        "  2. Reset" \
        "  3. Status" \
        "  4. Resolve" \
        "  0. Back" \
        "$SEP_LINE"
    if [ -z "$PRISMA_VER" ]; then
        printf '%s%s%s\n' "$C_RED" "$PRISMA_PROBLEM" "$C_RESET"
    fi
}

# `prisma migrate resolve` needs --applied or --rolled-back plus a migration
# name, so ask for both instead of running it bare (which only errors).
# Returns 2 when the user cancelled (no pause needed).
resolve_migration() {
    local mode choice name dir="${PRISMA_DIR}/prisma/migrations"

    clear
    printf '%s\n' \
        "$SEP_LINE" \
        "    Resolve a Migration" \
        "$SEP_LINE" \
        "  1. Mark as applied      (--applied)" \
        "  2. Mark as rolled back  (--rolled-back)" \
        "  0. Cancel" \
        "$SEP_LINE"
    read -r -p "> Select: " choice || return 2
    case "$choice" in
        1) mode="--applied" ;;
        2) mode="--rolled-back" ;;
        *) return 2 ;;
    esac

    if [ -d "$dir" ]; then
        printf '\nMigrations found in %s:\n' "$dir"
        find "$dir" -mindepth 1 -maxdepth 1 -type d -exec basename {} \; 2>/dev/null | sort | sed 's/^/  /'
    fi
    printf '\n'
    read -r -p "> Migration name (empty to cancel): " name || return 2
    [ -n "$name" ] || return 2

    case "$name" in
        *[!A-Za-z0-9_-]*)
            err "Invalid migration name. Use only letters, digits, '_' and '-' (the folder name)."
            return 1
            ;;
    esac

    run_prisma migrate resolve "$mode" "$name"
}

# Reset does its own confirmation instead of relying on Prisma's prompt (which
# `migrate dev` skips after drift, leaving no way to proceed). It shows which
# database will be wiped, requires the word RESET, then runs
# `prisma migrate reset --force` (--force = skip Prisma's own prompt).
# Returns 2 when cancelled (no pause needed).
reset_database() {
    local target answer

    printf 'Looking up the target database...\n'
    target="$(cd -- "$PRISMA_DIR" && npx prisma migrate status </dev/null 2>&1 \
        | sed 's/\x1b\[[0-9;?]*[A-Za-z]//g' | grep -m1 '^Datasource ')"
    [ -n "$target" ] || target="(could not determine it - check DATABASE_URL in .env)"

    printf '\n'
    printf '%sThis drops the database, re-applies every migration and runs the seeder.%s\n' "$C_RED" "$C_RESET"
    printf '%sAll data will be lost.%s\n\n' "$C_RED" "$C_RESET"
    printf '  Target: %s\n\n' "$target"
    read -r -p "> Type RESET to continue (anything else cancels): " answer || return 2
    [ "$answer" = "RESET" ] || { printf '\nCancelled.\n'; sleep 1; return 2; }

    run_prisma migrate reset --force
}

while true; do
    show_menu
    read -r -p "> Select: " choice || break

    case "$choice" in
        1)
            printf "\nRunning Migrate...\n"
            run_prisma migrate dev
            pause
            ;;
        2)
            printf "\nRunning Reset...\n"
            reset_database
            [ $? -eq 2 ] || pause
            ;;
        3)
            printf "\nChecking Migration Status...\n"
            run_prisma migrate status
            pause
            ;;
        4)
            resolve_migration
            [ $? -eq 2 ] || pause
            ;;
        0)
            exit 0
            ;;
        '')
            ;;
        *)
            printf "\nInvalid option.\n"
            sleep 1.5
            ;;
    esac
done
