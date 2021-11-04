
#!/bin/sh

trap 'echo "Goodbye, parser"; kill $PID' INT TERM

echo "parser running"

/app/parser &
PID=$!

wait $PID