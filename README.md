# BlockChain
# Computer Systems Final Project Fall 2017

Proto file generation    
protoc -I ./ --python_out=. --grpc_out=. --plugin=protoc-gen-grpc=`grpc` ./db.proto

### Prerequisite:
    - Go installed and gopath/gobin appropriately set     
    - google protobuf installed      
    - Python 2.7/2.9      
    - FUSE installed      
    - All python dependencies can be installed using requirements.txt       
    

### How to run:
    -  make myserver : Spawns default 3 servers
    -  make runI :  runs the integration test
    -  make runU :  runs the unit test cases
    -  ./run_fs.py  <mount directory>: Runs the blockchain filesystem
    - To use the cli:
        Navigate to MiniBlockChain folder and run ./setup.sh

### Contributors:
    - Akshaya Khare
    - Rishab Khandelwal
    - Sarita Joshi

### Project report and presentation deck available on demand
    


