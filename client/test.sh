# Automation script for running the client
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./unit_test.sh <port-number>

#!/usr/bin/env bash
# sleep 5
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=10 -fee=1
sleep 7
./client --p=$1 --ip=localhost -T=GET  --user=00000
