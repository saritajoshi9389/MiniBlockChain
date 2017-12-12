package main

import (
	hasher "./hash"
	"errors"
	"fmt"
	//"regexp"
	"github.com/golang/protobuf/jsonpb"
	p_b "../protobuf/go"
	"regexp"
)
/*
       Add a new block.
*/
func (self *BlockChain) AddBlock(block *MyChain, flag bool) {
	self.lockPendingBlocks.Lock()
	defer self.lockPendingBlocks.Unlock()
	if _, ok := self.list_processed[block.hash_value]; !ok {
		self.list_processed[block.hash_value] = true
		if flag {
			// server queries and adds to the block
			self.pending_blocks.PushFront(block)
		} else {
			// P2P push block, not new, flag = false
			self.pending_blocks.PushBack(block)
		}
		self.msg.wake()
	}
}
// All below condition should match
func (self *BlockChain) ValidateBlock(json *p_b.JsonBlockString) (*MyChain, error) {
	loggerB.Warning("Enters block validation")
	js := json.Json
	var tmp_string string
	tmp_string = *js
	hash := hasher.GetHashString(tmp_string)
	if !hasher.CheckHash(hash) {
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Error!!!!!!!!")
	}
	b := p_b.Block{}
	error := jsonpb.UnmarshalString(tmp_string, &b)
	if error != nil {
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Error!!!!!!!!")
	}
	prevHash := b.PrevHash
	var tmp_hash string
	tmp_hash = *prevHash
	if !hasher.CheckHash(tmp_hash) {
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Error!!!!!!!!")
	}
	var miner_id int
	var tmp_miner_id string
	tmp_miner_id = *b.MinerID
	_, miner_error := fmt.Sscanf(tmp_miner_id, "Server%02d", &miner_id)
	if miner_error != nil || !(miner_id <= len(server_list) && miner_id >= 1) {
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Error!!!!!!!")
	}
	var tmp_nouce string
	tmp_nouce = *b.Nonce
	if match, regerror := regexp.MatchString("\\d{8}", tmp_nouce); len(tmp_nouce) != 8 || !match || regerror != nil {
		loggerB.Warningf(" Make sure u check this %v %v %d", match, regerror, len(tmp_nouce))
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Errror!!!!!!!!!")
	}
	if len(b.Transactions) > MAX_TRANSACTIONS_IN_BLOCK {
		loggerB.Warning("Enters block validation")
		return nil, errors.New("Transactions threshold reached!!!")
	}
	for _, t := range b.Transactions {
		var tmp1 int32
		var tmp2 int32
		tmp1 = *t.Value
		tmp2 = *t.MiningFee
		if tmp1 < tmp2 {
			loggerB.Warning("Enters block validation")
			return nil, errors.New("invalid mining fee!!!!!!!!!!")
		}
	}
	loggerB.Warning("Enters block validation")
	return &MyChain{block: b, hash_value: hash, json: tmp_string}, nil

}
/*
	This is when the nounce is generated for the block and we need to add it to the chain
*/
func (self *BlockChain) InsertBlock(block *MyChain) (bool, error) {
	loggerB.Info("!!!!!!!!!!Insert BLock!!!!!!!!!!!!!!")
	self.lockBlocks.Lock()
	defer self.lockBlocks.Unlock()
	var tmp_hash string
	tmp_hash = *block.block.PrevHash
	parent := self.blocks[tmp_hash]
	var tmp1, tmp2 int32
	tmp1 = *block.block.BlockID
	tmp2 = *parent.block.BlockID
	shuffle := false
	if tmp1 != tmp2 + 1 {
		// Invalid Height
		return false, errors.New(fmt.Sprintf("Error::  block >>>>> ::%dparent >>>>::%d",
			block.block.BlockID, parent.block.BlockID))
	}
	block.data = make(map[string]int32)
	for k, v := range parent.data {
		block.data[k] = v
	}
	// Check for duplicacy at block level, like transaction
	ifExists := make(map[string]bool)
	for _, t := range block.block.Transactions {
		var tmp_uuid string
		tmp_uuid = *t.UUID
		id, hash := self.FetchUUID(tmp_uuid, parent)
		if id > OK {
			return false, errors.New(fmt.Sprintf("Errror UUID>>>> :: %s, hash >>>:: %s", tmp_uuid, hash))
		}

		var tmp_fromid string
		tmp_fromid = *t.FromID
		var tmp_val int32
		tmp_val = *t.Value // When block has to be added, bal +/- activity
		if _, ok := block.ModifyBalance(tmp_fromid, - tmp_val); ok != nil {
			return false, errors.New(fmt.Sprintf("Oops! balance lesser for :: >>>  %s", tmp_fromid))
		}
		if _, ok := ifExists[tmp_uuid]; ok {
			return false, errors.New(fmt.Sprintf("If already found, no duplicacy please!! %v", t))
		}
		var tmp_toid string
		tmp_toid = *t.ToID
		var tmp_fee int32
		tmp_fee = *t.MiningFee
		var tmp_minerid string
		tmp_minerid = *block.block.MinerID
		block.ModifyBalance(tmp_toid, tmp_val - tmp_fee) // all good case
		block.ModifyBalance(tmp_minerid, tmp_fee) // all good case
		ifExists[tmp_uuid] = true // processed, self flag to be true now
	}
	block.prev_block = make([]*MyChain, 0, BATCH_TRANSACTIONS) // keep track until 5 blocks
	temphash := INITIAL_HASH
	counter := 0
	for tmp := int32(1); tmp < tmp1; tmp *= 2 {
		if tmp2 % tmp != OK {
			block.prev_block = append(block.prev_block, parent.prev_block[counter])
			temphash += fmt.Sprintf("%d", parent.prev_block[counter].block.BlockID)

		} else {
			block.prev_block = append(block.prev_block, parent)
			temphash += fmt.Sprintf("%d", parent.block.BlockID)
		}
		counter++
	}
	self.blocks[block.hash_value] = block
	for _, t := range block.block.Transactions {
		var tmp_uuid string
		tmp_uuid = *t.UUID
		if _, ok := self.list_of_uuid[tmp_uuid]; !ok {
			self.list_of_uuid[tmp_uuid] = []string{block.hash_value}
		} else {
			// if not exist, create one
			self.list_of_uuid[tmp_uuid] = append(self.list_of_uuid[tmp_uuid], block.hash_value) // else append to existing
		}
	}
	var tmp_long_id int32
	tmp_long_id = *self.long_chain.block.BlockID
	if (tmp_long_id == tmp1) || (tmp_long_id < tmp1) {
		self.long_chain = block
		shuffle = true
	}
	return shuffle, nil
}

// Fetches the UUID , the unique identification
func (self *BlockChain) FetchUUID(uuid string, block *MyChain) (int32, string) {
	hashval, ok := self.list_of_uuid[uuid]
	if ok {
		for _, val := range hashval {
			parent, ok := self.blocks[val]
			if !ok {
				loggerB.Errorf("Errooor again!!! %v", ok)
			}
			if self.IsParentUUIDFound(block, parent) {
				var tmp int32
				tmp = *parent.block.BlockID
				return tmp, parent.hash_value
			}
		}
	}
	return OK, INITIALIZE_STRING
}
// Check for parent, return a bool flag true or false
func (self *BlockChain) IsParentUUIDFound(b *MyChain, p *MyChain) bool {
	if p.hash_value == INITIAL_HASH {
		return true
	} // first block base case
	var tmp1, tmp2 int32
	tmp1 = *b.block.BlockID
	tmp2 = *p.block.BlockID
	if tmp1 <= tmp2 {
		return b.hash_value == p.hash_value
	}
	target := tmp2
	curr := b
	// Keep checking
	for {
		var tmp1, tmp2 int32
		tmp1 = *curr.block.BlockID
		tmp2 = *p.block.BlockID
		if tmp1 <= tmp2 {
			loggerB.Error("No No No!!!! parent error")
		}
		for _, tmp := range curr.prev_block {
			var tp int32
			tp = *tmp.block.BlockID
			if tp == target {
				// found, return
				return tmp.hash_value == p.hash_value
			} else if tp < target {
				// reach extreme break
				break
			}// traverse the doubly link list
			curr = tmp
		}
	}

}