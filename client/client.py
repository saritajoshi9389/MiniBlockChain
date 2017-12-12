import grpc
import sys
import imp
sys.path.append('../..')
gp = imp.load_source('db_pb2_grpc', './protobuf/db_pb2_grpc.py')
p_b = imp.load_source('db_pb2', './protobuf/db_pb2.py')
# import protobuf.db_pb2_grpc as gp
# import protobuf.db_pb2 as p_b
def run():
    channel = grpc.insecure_channel('localhost:5001')
    stub = gp.BlockChainMinerStub(channel)
    print("-------------- GetFeature --------------")
    print(stub.Get(p_b.GetRequest(UserID="12345")))
    print("-------------- CheckBlock --------------")
    print((stub.GetBlock(p_b.GetBlockRequest(BlockHash="all"))))

if __name__ == '__main__':
  run()
