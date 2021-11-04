#!/bin/sh

trap 'echo "Goodbye, Init"; kill -TERM $PID1 $PID2' TERM INT

/app/api.sh &
PID1=$!

/app/parser.sh &
PID2=$!

wait $PID1
wait $PID2

EXIT_STATUS=$?