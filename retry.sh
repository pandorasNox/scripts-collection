#!/usr/bin/env bash

# inspiration: https://unix.stackexchange.com/questions/82598/how-do-i-write-a-retry-logic-in-script-to-keep-retrying-to-run-it-upto-5-times

if [ $# -ne 3 ]; then
    echo 'usage: retry <num retries> <wait retry secs> "<command>"'
    exit 1
fi

retries=$1
wait_retry=$2
command=$3

for i in `seq 0 $retries`; do
    if [ ! $i -eq 0 ]; then echo ""; echo ""; echo ""; fi #just formatting
    echo "retry no: ${i}"
    echo "$command"
    $command
    ret_value=$?
    [ $ret_value -eq 0 ] && break
    echo "> failed with return value '$ret_value'"

    if [ ! $i -eq $retries ]; then
        echo ""; echo ""; echo "";
        echo "... waiting ${wait_retry} seconds to retry ..."
        sleep $wait_retry
    fi
done

exit $ret_value
