#!/bin/bash

yaml=(
  "assets/ceph/crds.yaml"
  "assets/ceph/common.yaml"
  "assets/ceph/operator.yaml"
  "assets/ceph/toolbox.yaml"
  "assets/ceph/csi/rbd/storageclass.yaml"
  "assets/ceph/filesystem.yaml"
  "assets/ceph/csi/cephfs/storageclass.yaml"
  "assets/ceph/object.yaml"
  "assets/ceph/storageclass-bucket-delete.yaml"
  "assets/ceph/object-bucket-claim-delete.yaml")

for var in ${yaml[@]}; do
  n=0
  until [ $n -ge 10 ]; do
    kubectl apply -f $var && break
    n=$(($n + 1))
    sleep 10
  done
done