#!/bin/bash

n=0
until [ $n -ge 10 ]; do
  istioctl operator init --hub harbor1.zlibs.com/istio && break
  n=$(($n + 1))
  sleep 10
done

n=0
until [ $n -ge 10 ]; do
  kubectl apply -f $var && break
  n=$(($n + 1))
  sleep 10
done