#!/bin/bash

yaml=(
  "ceph/crds.yaml"
  "ceph/common.yaml"
  "ceph/operator.yaml"
  "ceph/toolbox.yaml"
  "ceph/csi/rbd/storageclass.yaml"
  "ceph/filesystem.yaml"
  "ceph/csi/cephfs/storageclass.yaml"
  "ceph/object.yaml"
  "ceph/storageclass-bucket-delete.yaml"
  "ceph/object-bucket-claim-delete.yaml")

for var in ${yaml[@]}; do
  n=0
  until [ $n -ge 10 ]; do
    kubectl apply -f $var && break
    n=$(($n + 1))
    sleep 10
  done
done