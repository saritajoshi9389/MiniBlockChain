# Automation script for spawing n servers.
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./spawn_n_servers.sh <No-of-server>
#!/usr/bin/env bash
port=5000
value=$1
#for number in {1..99}
for number in `seq 1 $value`;
do
    temp=$(expr $port + $number)
    pid=$(lsof -i:$temp -t); kill -TERM $pid || kill -KILL $pid
    ./start_server.sh $number &
#echo $number
done
exit 0
