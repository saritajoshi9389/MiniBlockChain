package main

import (
	"container/list"
	"sync"
	p_b "../protobuf/go"
	"os"
)

type MyChain struct {
	block      p_b.Block
	hash_value string
	json       string
	prev_block []*MyChain
	data       map[string]int32
	lockData sync.RWMutex
}
/*
	Check if id exists, if not set to initial balance
*/
func (self *MyChain) GetUserBalance(id string) int32 {
	if value, ok := self.data[id]; ok{
		return value
	}
	return INITIAL_BALANCE
}

func (self *MyChain) MinerTransactionFilter() {
	loggerMiner.Infof(" !!!!!!!!!!!!!!!!! Check for Pending transactions!!!!!!!!!!!")
	self.lockData.RLock()
	tmp := MyChain{
		block: p_b.Block{
			BlockID:      CreateInt(INITIAL_BLOCK_ID),
			PrevHash:     CreateString(PREV_HASH),
			Transactions: []*p_b.Transaction{},
			MinerID:      CreateString(INITIAL_MINER_ID),
			Nonce:        CreateString(INITIAL_NONCE),
		},
		hash_value: INITIAL_HASH,
		json:       INITIALIZE_STRING,
		prev_block: []*MyChain{},
		data:       map[string]int32{},
	}
	for key, value := range self.data { tmp.data[key] = value}
	self.lockData.RUnlock()
	var next_ele *list.Element
	for ele := transaction.data.Front(); ele != nil; ele = next_ele {
		next_ele = ele.Next()
		t, ok := ele.Value.(*p_b.Transaction)
		if !ok {
			loggerMiner.Errorf("Panic!!!!")
			os.Exit(2)
		}
		// If valid transaction, update the amount for both to and from
		if tmp.VerifyTransactionValue(*t.FromID, *t.Value) {
			tmp.ModifyBalance(*t.FromID, -*t.Value)
			tmp.ModifyBalance(*t.ToID, *t.Value - *t.MiningFee)
		} else {
			// If not, just delete for now. Later we can see how to figure this out
			// The delete built-in function deletes the element with the specified key
			// (m[key]) from the map. If m is nil or there is no such element, delete
			// is a no-op.
			transaction.data.Remove(ele)
			delete(transaction.list_of_uuid, *t.UUID)
		}
	}
}