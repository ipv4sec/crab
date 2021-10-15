#!/bin/bash

n=0
until [ $n -ge 10 ]; do
  istioctl operator init --hub harbor1.zlibs.com/istio && break
  n=$(($n + 1))
  sleep 10
done