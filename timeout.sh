#!/bin/sh

# wrapper for the 'timeout' cmd which adds output for the failure case
# preserves original SIG handling
# only works with int (seconds) as first argument

# usage example:
# `sh timeout.sh 3 sleep 2`
# `sh timeout.sh 1 sleep 2`

timeout_int=$1; shift
cmd_and_args=$@;

timeout -t $timeout_int $cmd_and_args;

TIMEOUT_EXIT_CODE=$?

if [ ! $TIMEOUT_EXIT_CODE -eq 0 ]
then
    echo "__TIMEOUT__: the command \"$cmd_and_args\" has exceeded the given timeout of $timeout_int seconds and was therefore terminated"
fi

exit $TIMEOUT_EXIT_CODE
