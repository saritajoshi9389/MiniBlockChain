#!/usr/bin/env bash

sleep 7
for run in {1..100}
do
  ./client --p=$1 --ip=localhost -T=TRANSFER --from=00002 --to=12347 --value=3 -fee=1
done
