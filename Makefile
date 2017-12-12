# Automation script for spawing n servers.
# Final Project CS5600
# Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
# This is the sample Makefile provided by professor.
# All rules here can be used as , make <rule-name>
# This Makefile will be to simulate a working environment for 3 servers and running and
# validating some results.
CC=gcc -g -O3
GO=go build
all: 	clean
default: check
clean:
	rm -rf server/myserver client/client client/IntegrationTest
myclient:
	./run_client.sh
unit:
	cd client && \
	./run_unit_test.sh 5001 && printf "\n Unit Test Completed!!!\n"
dependencies: requirement.txt
		pip3 install -r requirement.txt
myserver:
	./spawn_n_servers.sh 3
stopserver:
	./stop_n_servers.sh 3
integration:
	./build_integration.sh
all:
	cd client && \
	./myintegration && printf "\n Integration Test Completed!!!\n"
runU: myclient unit
	printf "\n Great Job! Simulation Completed (Unit Test) \n"
runI:
	./build_integration.sh && sleep 5 && cd client && ./myintegration && printf "\n Great Job! Simulation Completed (Integration Test)\n"
stop: stopserver
fs:
	./run_fs.sh client/FS
