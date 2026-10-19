#!/usr/bin/env bash
# Prisma Console - interactive menu.
#
# Works both when launched from the main start.sh (which exports PROJECT_ROOT)
# and when run directly: the project root is then derived from this script's
# location, never from $PWD. Type an option and press [Enter] to confirm.

SELF_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
SCRIPTS_DIR="${SELF_DIR}/scripts/bash"

if [ ! -r "${SCRIPTS_DIR}/common.sh" ]; then
    printf 'Error: missing %s/common.sh\n' "$SCRIPTS_DIR" >&2
    exit 1
fi
# shellcheck source=scripts/bash/common.sh
. "${SCRIPTS_DIR}/common.sh" || exit 1

prisma_init
export PROJECT_ROOT="$ROOT_DIR" PRISMA_VER    # migration.sh reuses these

trap 'printf "\n"; exit 0' INT

show_menu() {
    clear
    printf '%s\n' \
        "$SEP_LINE" \
        "    Prisma Console" \
        "    Prisma Version ${PRISMA_VER:-unknown}" \
        "$SEP_LINE" \
        "  1. Migration" \
        "  2. Seed Database" \
        "  3. Format, Validate & Generate Prisma Client" \
        "  0. Exit" \
        "$SEP_LINE"
    if [ -z "$PRISMA_VER" ]; then
        printf '%s%s%s\n' "$C_RED" "$PRISMA_PROBLEM" "$C_RESET"
    fi
}

seed() {
    run_prisma db seed
}

format_validate_generate() {
    # Stop at the first failing step, like `a && b && c`.
    run_prisma format && run_prisma validate && run_prisma generate
}

while true; do
    show_menu
    read -r -p "> Select: " choice || break      # EOF (Ctrl+D) exits

    case "$choice" in
        1)
            printf "\nOpening Migration Console...\n"
            bash "${SCRIPTS_DIR}/migration.sh"
            ;;
        2)
            printf "\nRunning Seeder...\n"
            seed
            pause
            ;;
        3)
            printf "\nFormatting, Validating & Generating Prisma Client...\n"
            format_validate_generate
            pause
            ;;
        0)
            printf "\nExiting...\n"
            sleep 1
            clear
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
