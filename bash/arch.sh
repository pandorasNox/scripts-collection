#!/usr/bin/env sh
# Script to evaluate to systems process architecture

set -o errexit
set -o nounset
# set -o xtrace

if set +o | grep -q -F 'set +o pipefail'; then
  # shellcheck disable=SC3040
  set -o pipefail
fi

if set +o | grep -q -F 'set +o posix'; then
  # shellcheck disable=SC3040
  set -o posix
fi

function func_command_exist() {
    FUNC_ARG_CMD=${1:?}

    if ! command -v ${FUNC_ARG_CMD} &> /dev/null; then
        return 1
    fi

    return 0
}

function func_arch() {
    ARCH_CMD_PATHS_TO_CHECK=("/usr/bin/arch" "/bin/arch" "arch")

    ARCH_CMD=""
    for ARCH_CMD_PATH_TO_CHECK in "${ARCH_CMD_PATHS_TO_CHECK[@]}"
    do
        if func_command_exist "${ARCH_CMD_PATH_TO_CHECK}"; then
            ARCH_CMD="${ARCH_CMD_PATH_TO_CHECK}"
            break
        fi
    done

    if test -z "${ARCH_CMD}"; then
        ARCH_CMD="uname -m"
    fi

    printf '%s' "$(${ARCH_CMD})"
}

# echo detected arch: $(func_arch)
func_arch
