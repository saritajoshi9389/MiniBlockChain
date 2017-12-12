# Automation script for running the client
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./unit2_test.sh <port-number>

#!/usr/bin/env bash
./client --p=$1 --ip=localhost -T=TRANSFER --from=00002 --to=12347 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00002 --to=12347 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=12347 --to=00002 --value=6 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=12347 --to=00002 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=12347 --to=00002 --value=10 -fee=2
sleep 7
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12347 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12347 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=12347 --to=00000 --value=6 -fee=1
./client --p=$1 --ip=localhost  -T=TRANSFER --from=12347 --to=00000 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=12347 --to=00000 --value=10 -fee=2
#./client --p=$1 --ip=localhost -T=GET  --user=00001
#./client --p=$1 --ip=localhost  -T=GET  --user=12346
#./client --p=$1 --ip=localhost  -T=GET  --user=00000
#./client --p=$1 --ip=localhost -T=GET  --user=12345