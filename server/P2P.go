package main

import (
	"fmt"
	p_b "../protobuf/go"
	"golang.org/x/net/context"
	"time"
	"strings"
)
// User loggerP2P
type P2PServer struct{}

func printAllBlocks() {
	iterator := block_chain.blocks
	for _, tx := range iterator {
		logger.Noticef("Block->%v", tx.json)
	}
}

func CreateCustomData(types p_b.VerifyResponse_Results) *p_b.VerifyResponse_Results {
	return &types
}

func (c *P2PServer) Get(ctx context.Context, in *p_b.GetRequest) (*p_b.GetResponse, error) {
	loggerP2P.Infof("$$$$RPCRequest$$$$$$$$Get$$$$ %v", in)
	block_chain.lockBlocks.RLock()
	defer block_chain.lockBlocks.RUnlock()
	inUserId := *in.UserID
	answer := block_chain.long_chain.GetUserBalance(inUserId)
	printAllBlocks()
	fmt.Println(time.Now().Format(time.RFC850))
	return &p_b.GetResponse{Value: CreateInt(answer)}, nil
}

func (c *P2PServer) Transfer(ctx context.Context, in *p_b.Transaction) (*p_b.BooleanResponse, error) {
	fmt.Println(time.Now().Format(time.RFC850))
	loggerP2P.Infof("$$$$RPC RC TRANSFER$$$$ %v", in)
	printAllBlocks()
	fmt.Println(time.Now().Format(time.RFC850))
	inFromId := *in.FromID
	inToId := *in.ToID
	inMiningFee := *in.MiningFee
	inValue := *in.Value
	// if both same, or mining > actual amt or mining fee < 0
	if inFromId == inToId  || inMiningFee >= inValue || inMiningFee < 0  { return &p_b.BooleanResponse{ Success: CreateBool(false)}, nil }
	block_chain.lockBlocks.RLock()
	// Verification of amount failed
	if !block_chain.long_chain.VerifyTransactionValue(inFromId, inValue) {
		return &p_b.BooleanResponse{Success: CreateBool(false)}, nil
	}
	block_chain.lockBlocks.RUnlock()
	// Adding the transaction fails
	if !transaction.AddNewTransaction(in) { return &p_b.BooleanResponse{Success: CreateBool(false)}, nil }
	// Else success, push the transaction
	success := PushTransactions(in)
	loggerP2P.Info("$$$$RPC RC TRANSFER$$$$", success)
	return &p_b.BooleanResponse{Success: CreateBool(true)}, nil
}

func (c *P2PServer) Verify(ctx context.Context, in *p_b.Transaction) (*p_b.VerifyResponse, error) {
	printAllBlocks()
	inFromId := *in.FromID
	inUUID := *in.UUID
	inValue := *in.Value
	var tmp int32
	tmp = *block_chain.long_chain.block.BlockID
	block_chain.lockBlocks.RLock()
	defer block_chain.lockBlocks.RUnlock()
	loggerP2P.Infof("$$$$ RPC REQ VER $$$$ %v", in)
	if p_id, hash_val := block_chain.VerifyTransaction(in, block_chain.long_chain); p_id == OK {
		if block_chain.long_chain.VerifyTransactionValue(inFromId, inValue) {
			transaction.lockData.RLock()
			defer transaction.lockData.RUnlock()
			if _, success := transaction.list_of_uuid[inUUID]; success {
				return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_PENDING), BlockHash: CreateString(INITIALIZE_STRING)}, nil
			} else {
				return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_FAILED), BlockHash: CreateString(INITIALIZE_STRING)}, nil
			}
		} else { return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_FAILED), BlockHash: CreateString(INITIALIZE_STRING)}, nil }
	} else if p_id == INCORRECT_PID {
		return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_FAILED), BlockHash: CreateString(INITIALIZE_STRING)}, nil
	} else if p_id  <= tmp + BUFFER_BLOCKS {
		return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_SUCCEEDED), BlockHash: CreateString(hash_val)}, nil
	} else {
		return &p_b.VerifyResponse{Result: CreateCustomData(p_b.VerifyResponse_PENDING), BlockHash: CreateString(hash_val)}, nil }

}

func (c *P2PServer) GetHeight(ctx context.Context, in *p_b.Null) (*p_b.GetHeightResponse, error) {
	block_chain.lockBlocks.RLock()
	defer block_chain.lockBlocks.RUnlock()
	return &p_b.GetHeightResponse{Height: block_chain.long_chain.block.BlockID, LeafHash: CreateString(block_chain.long_chain.hash_value)}, nil
}

func (c *P2PServer) GetBlock(ctx context.Context, in *p_b.GetBlockRequest) (*p_b.JsonBlockString, error) {
	loggerP2P.Infof("$$$$RPC REQ GET BLK$$$$ %v", in)
	hash := *in.BlockHash
	if hash == INITIAL_HASH { return &p_b.JsonBlockString{Json: CreateString("")}, nil }
	block_chain.lockBlocks.RLock()
	defer block_chain.lockBlocks.RUnlock()
	if hash == "all"{
		test := "["
		iterator := block_chain.blocks
		for _, tx := range iterator {
			test += tx.json + ","
		}
		if(len(test) > 1) {test = test[:len(test)-1]}
		test += "]"
		test = strings.Replace(test, ",,", ",", -1)
		test = strings.Replace(test, ",]", "]", -1)
		test = strings.Replace(test, "[,", "[", -1)
		return &p_b.JsonBlockString{Json: CreateString(test)}, nil
	}
	val, success := block_chain.blocks[hash]
	if success { return &p_b.JsonBlockString{Json: CreateString(val.json)}, nil } else {
		return &p_b.JsonBlockString{Json: CreateString("")}, nil
	}
}
// Push to n-1 servers
func (c *P2PServer) PushBlock(ctx context.Context, in *p_b.JsonBlockString) (*p_b.Null, error) {
	if b, error := block_chain.ValidateBlock(in); error == nil { block_chain.AddBlock(b, false)} else {
		loggerP2P.Errorf("$$$$RPC PUSH BLOCK$$$$ %v ", error)
	}
	loggerP2P.Infof("$$$$RPC PUSH BLOCK$$$$")
	return &p_b.Null{}, nil
}
// Similar for pushing transactions
func (c *P2PServer) PushTransaction(ctx context.Context, in *p_b.Transaction) (*p_b.Null, error) {
	loggerP2P.Infof("$$$$RPC REQ PUSH TRANS$$$$ %v", in)
	inFromId := *in.FromID
	inToId := *in.ToID
	inMiningFee := *in.MiningFee
	inValue := *in.Value
	if inFromId == inToId || inMiningFee >= inValue || inMiningFee < 0  { return &p_b.Null{}, nil }
	block_chain.lockBlocks.RLock()
	if !block_chain.long_chain.VerifyTransactionValue(inFromId, inValue) { return &p_b.Null{}, nil }
	block_chain.lockBlocks.RUnlock()
	transaction.AddNewTransaction(in)
	loggerP2P.Infof("$$$$RPC REQ PUSH TRANS$$$$")
	return &p_b.Null{}, nil
}