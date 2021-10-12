#!/bin/bash

n=0
until [ $n -ge 10 ]; do
  kubectl apply -f assets/island/island-console.yaml && break
  n=$(($n + 1))
  sleep 10
done