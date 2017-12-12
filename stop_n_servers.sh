#!/usr/bin/env bash
# Automation script for stopping servers
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./stop_n_servers <total-servers_spawned>
port=5000
value=$1
for number in `seq 1 $value`;
do
    temp=$(expr $port + $number)
    pid=$(lsof -i:$temp -t); kill -TERM $pid || kill -KILL $pid
done
exit 0