#!/bin/bash

/app/parser &
status=$?
if [ $status -ne 0 ]; then
  exit $status
fi

/app/api &
status=$?
if [ $status -ne 0 ]; then
  exit $status
fi

while /bin/true; do
  ps aux |grep parser |grep -q -v grep
  PROCESS_1_STATUS=$?
  ps aux |grep api |grep -q -v grep
  PROCESS_2_STATUS=$?

  if [ $PROCESS_1_STATUS -ne 0 -o $PROCESS_2_STATUS -ne 0 ]; then
    exit -1
  fi
  sleep 6
done