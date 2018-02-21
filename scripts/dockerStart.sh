#!/bin/sh -e

# Remove old pid file if it exists
rm -f "/var/run/vflow.pid"
# Touch the log file so tail has something to follow if it doesn't exist
touch "/var/log/vflow.log"

# Continuously provide logs so that 'docker logs' can produce them
tail -F "/var/log/vflow.log" &
"/usr/bin/vflow" &
vflow_pid="$!"

trap "echo Received trapped signal, beginning shutdown...;" KILL TERM HUP INT EXIT;

echo vFlow running with PID ${vflow_pid}.
wait ${vflow_pid}