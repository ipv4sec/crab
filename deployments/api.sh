
#!/bin/sh

trap 'echo "Goodbye, api"; kill $PID' INT TERM

echo "api running"

/app/api &
PID=$!

wait $PID