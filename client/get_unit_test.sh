# Automation script for running the unit test cases
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./get_unit_test.sh <port-number>
#!/usr/bin/env bash
./client --p=$1 --ip=localhost -T=GET  --user=00000
./client --p=$1 --ip=localhost -T=GET  --user=12345
./client --p=$1 --ip=localhost -T=GET  --user=00001
./client --p=$1 --ip=localhost -T=GET  --user=12346