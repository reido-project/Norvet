#!/usr/bin/env bash

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"

# Hand the project root to every sub-menu so none of them depends on $PWD.
export API_PROJECT_ROOT="$ROOT_DIR"

OPS_DIR="${ROOT_DIR}/ops"

API_SCRIPT_DIR="${OPS_DIR}/api"
MIGRATION_SCRIPT_DIR="${OPS_DIR}/prisma"
GIT_SCRIPT_DIR="${OPS_DIR}/git"

# Static shell code. Paths are passed as positional arguments.
# (Single quotes: $1/$2/$status must expand inside the new terminal's bash -c,
# not here where they would be empty.)
TERMINAL_RUNNER='cd -- "$1" || exit 1; bash "$2"; status=$?; printf "\nScript exited with status %s.\n" "$status"; exec bash'

open_in_new_terminal() {
    local script_path="$1"
    local script_dir
    local terminal

    if [[ ! -f "$script_path" || ! -r "$script_path" ]]; then
        printf 'Error: Script not found or unreadable: %s\n' "$script_path" >&2
        return 1
    fi

    script_dir="$(cd -- "$(dirname -- "$script_path")" && pwd)" || return 1
    script_path="${script_dir}/$(basename -- "$script_path")"

    # macOS Terminal
    if [[ "$OSTYPE" == darwin* ]]; then
        if ! command -v osascript >/dev/null 2>&1; then
            printf 'Error: osascript not found.\n' >&2
            return 1
        fi

        osascript - "$script_dir" "$script_path" <<'APPLESCRIPT' >/dev/null 2>&1
on run argv
    set scriptDir to quoted form of (item 1 of argv)
    set scriptPath to quoted form of (item 2 of argv)
    set shellCommand to "cd -- " & scriptDir & " && bash " & scriptPath & "; status=$?; printf '\\nScript exited with status %s.\\n' \"$status\"; exec bash"

    tell application "Terminal"
        do script shellCommand
        activate
    end tell
end run
APPLESCRIPT
        return $?
    fi

    # Linux / BSD: terminals accepting argument-based execution.
    for terminal in gnome-terminal konsole kitty alacritty tilix xterm x-terminal-emulator; do
        command -v "$terminal" >/dev/null 2>&1 || continue

        case "$terminal" in
            gnome-terminal)
                "$terminal" -- bash -c "$TERMINAL_RUNNER" bash "$script_dir" "$script_path" >/dev/null 2>&1 &
                ;;
            konsole|alacritty|tilix|xterm|x-terminal-emulator)
                "$terminal" -e bash -c "$TERMINAL_RUNNER" bash "$script_dir" "$script_path" >/dev/null 2>&1 &
                ;;
            kitty)
                "$terminal" bash -c "$TERMINAL_RUNNER" bash "$script_dir" "$script_path" >/dev/null 2>&1 &
                ;;
        esac

        disown "$!" 2>/dev/null || true
        return 0
    done

    # XFCE / MATE: command-string interfaces.
    for terminal in xfce4-terminal mate-terminal; do
        command -v "$terminal" >/dev/null 2>&1 || continue

        env SCRIPT_DIR="$script_dir" SCRIPT_PATH="$script_path" \
            "$terminal" --command="bash -c 'cd -- \"\$SCRIPT_DIR\" || exit 1; bash \"\$SCRIPT_PATH\"; status=\$?; printf \"\\nScript exited with status %s.\\n\" \"\$status\"; exec bash'" \
            >/dev/null 2>&1 &

        disown "$!" 2>/dev/null || true
        return 0
    done

    printf 'Error: No supported terminal was found!\n' >&2
    return 1
}

run_script() {
    local script_path="${1}/start.sh"
    local ask_new_term

    if [[ ! -f "$script_path" || ! -r "$script_path" ]]; then
        printf 'Error: Script not found or unreadable: %s\n' "$script_path" >&2
        read -r -p "Press [Enter] to continue..."
        return 1
    fi

    read -r -p "Start in new terminal? (y/n): " ask_new_term

    case "$ask_new_term" in
        [yY]|[yY][eE][sS])
            if open_in_new_terminal "$script_path"; then
                read -r -p "Press [Enter] to continue..."
            else
                read -r -p "Failed to launch terminal. Press [Enter] to continue..."
                return 1
            fi
            ;;
        *)
            printf "\nExecuting script in current terminal...\n\n"
            bash "$script_path"
            local status=$?
            printf "\nScript exited with status %s.\n" "$status"
            read -r -p "Press [Enter] to continue..."
            return "$status"
            ;;
    esac
}

show_menu() {
    clear
    printf '%s\n' \
        "——————————————————————————————————" \
        "    Select a program to start" \
        "——————————————————————————————————" \
        "  1. API" \
        "  2. Migration & Seeder" \
        "  3. Git" \
        "  0. Exit" \
        "——————————————————————————————————"
}

while true; do
    show_menu
    read -r -p "> Select: " choice

    case "$choice" in
        1)
            printf "\nStarting API...\n"
            run_script "$API_SCRIPT_DIR"
            ;;
        2)
            printf "\nRunning Migration & Seeder...\n"
            run_script "$MIGRATION_SCRIPT_DIR"
            ;;
        3)
            printf "\nOpening Git commands...\n"
            run_script "$GIT_SCRIPT_DIR"
            ;;
        0)
            printf "\nExiting...\n"
            sleep 1
            clear
            exit 0
            ;;
        *)
            printf "\nInvalid option.\n"
            sleep 1.5
            ;;
    esac
done
