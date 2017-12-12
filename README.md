# BlockChain
Computer Systems Final Project Fall 2017

Proto file generation    
protoc -I ./ --python_out=. --grpc_out=. --plugin=protoc-gen-grpc=`grpc` ./db.proto
