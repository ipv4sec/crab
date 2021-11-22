#!/bin/bash

name=''

while [ $# -gt 0 ]; do
  case "$1" in
  --name)
    name="$2"
    shift
    ;;
  --*)
    echo "Illegal option $1"
    ;;
  esac
  shift $(($# > 0 ? 1 : 0))
done

/usr/local/bin/kubectl label node $name island-storage=local --overwrite