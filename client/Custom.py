#  Dummy api call to implement fuse to the current implementation
import json
import grpc
import db_pb2_grpc as gp
import db_pb2 as p_b
import simplejson as json

# APIURL = 'https://insight.bitpay.com/api'

class Custom(object):
    def __init__(self):
        self.data = self.getAllBlocks()
        # print self.data

    def getAllBlocks(self):
        channel = grpc.insecure_channel('localhost:5001')
        stub = gp.BlockChainMinerStub(channel)
        val = stub.GetBlock(p_b.GetBlockRequest(BlockHash="all"))
        resp = str(val).lstrip("Json:")
        resp = resp.replace('\\','', -1)
        resp = resp[2:-2]
        data = json.dumps(resp.decode('utf-8'))
        new = json.loads(data.encode())
        l = json.loads(new)
        return l
        # print type(new), type(data), l[0]['BlockID']

    def blockhash_by_index(self, index):
        # 0-3
        if index+1 >= len(self.data):
            return self.data[len(self.data) -1]['PrevHash']
        for i in range(len(self.data)):
            if self.data[i]['BlockID'] == index+1:
                    return self.data[i]['PrevHash']

    def blockinfo(self, blockhash):
        for i in range(len(self.data)):
            if blockhash == self.data[i]['PrevHash']:
                return self.data[i-1]

    def txinfo(self, txhash):
        for i in range(len(self.data)):
            if txhash == self.data[i]['Transactions'][0]['UUID']:
                return self.data[i]['Transactions'][0]
def main():
    x = Custom()
    x.getAllBlocks()

if __name__ == '__main__':
    main()
