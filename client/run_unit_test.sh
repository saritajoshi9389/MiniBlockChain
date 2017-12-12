# Automation script for running the client
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./unit_test.sh <port-number>

#!/usr/bin/env bash
# sleep 5
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=6 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=10 -fee=2
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=6 -fee=1
sleep 17
sleep 7
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=10 -fee=2
sleep 7
./client --p=$1 --ip=localhost -T=GET  --user=00000
./client --p=$1 --ip=localhost -T=GET  --user=12345
./client --p=$1 --ip=localhost -T=GET  --user=00001
./client --p=$1 --ip=localhost -T=GET  --user=12346
