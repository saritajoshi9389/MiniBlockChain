#!/usr/bin/env bash
# Automation script for stopping servers
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# Usage ./stop_n_servers
rm -rf myserver
go build -o myserver *.go
sleep 2
./myserver --s_id=$1
