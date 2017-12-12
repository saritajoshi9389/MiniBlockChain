package main

import (
       "container/list"
       "sync"
       p_b "../protobuf/go"

)
//This messenger is used as a channel for waking and sleeping the blocks
type Messenger struct {
       channel chan int
}

// Block the channel;; needed for blocking call
func (self *Messenger) sleep() {
       <-self.channel
}

// channel buffer val to 0, unblocks and wake it up
func (self *Messenger) wake() {
       select {
       case self.channel <- 0:
       default:
       }
}
//Strcture for a blockchain. Here we are maintaining separate locks for separate blocks block 
//types such as fraud , pending and normal blocks.

type BlockChain struct {
       msg Messenger        
       blocks         map[string]*MyChain
       long_chain     *MyChain // refers to the longest chain in the blockchain among all the chains.
       list_of_uuid   map[string][]string
       pending_blocks *list.List
       list_processed map[string]bool
       fraud_if_any   map[string]bool//variable for keeping track of fraud blocks.
       //The fraud blocks are eventually removed from the list.

       lockBlocks        sync.RWMutex
       lockPendingBlocks sync.RWMutex
       lockFraud         sync.RWMutex
}
//An empty BlockChain is initiailised here 
func InitializeBlockChain() *BlockChain {
       var bc BlockChain
       new := MyChain{
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
       new.data = make(map[string]int32)
       bc.blocks = make(map[string]*MyChain)
       bc.blocks[new.hash_value] = &new
       bc.long_chain = &new
       bc.list_of_uuid = make(map[string][]string)
       bc.pending_blocks = list.New()
       bc.list_processed = make(map[string]bool)
       bc.fraud_if_any = make(map[string]bool)
       /*
              creates a buffered channel with a capacity of 1.Normally channels are synchronous;
               both sides of the channel will wait until the other side is ready. A buffered channel
               is asynchronous; sending or receiving a message will not wait unless the channel is already full.*/
       bc.msg.channel = make(chan int, 1)
       return &bc
}
/*
https://blog.golang.org/go-maps-in-action - Citation
Concurrency
Maps are not safe for concurrent use:
it's not defined what happens when you read and write to them simultaneously.
If you need to read from and write to a map from concurrently executing goroutines,
 the accesses must be mediated by some kind of synchronization mechanism.
 One common way to protect maps is with sync.RWMutex.
This statement declares a counter variable that is an anonymous struct
containing a map and an embedded sync.RWMutex.

Package list implements a doubly linked list.
To iterate over a list (where l is a *List): https://golang.org/pkg/container/list/

*/

func (self *BlockChain) StartBlocks() {
	loggerB.Notice("***********************Blockchain starts pushing the blocks to form MyChain********************")
	current_status_of_blocks := make(map[string]int8)
	var flip_flag bool
	current_block_list := list.New()
	for {
		self.msg.sleep() // block this thread
              	self.lockPendingBlocks.Lock()
              	// push the entire pending list to the current blocks to consider
              	current_block_list.PushBackList(self.pending_blocks)
              	self.pending_blocks = list.New()
              	self.lockPendingBlocks.Unlock()
              	flip_flag = false
              	var ele *list.Element
              	// Traverse now through the complete list
		for e := current_block_list.Front(); e != nil; e = ele {
                     ele = e.Next()
                     block := e.Value.(*MyChain) // get the current block i.e. list.Value()
                     val := current_status_of_blocks[block.hash_value]
                     // Block for that specific hash value not found, val = 0
                     if val == OK {
                            self.lockBlocks.RLock()
                            _, ok := self.blocks[block.hash_value] // key found
                            self.lockBlocks.RUnlock()
                            // If curr found with the hash value, set the key val with different numbers
                            if ok {
                                   // When the same block happens to be visited make it back as 1
                                   current_status_of_blocks[block.hash_value] = NOUNCE_FOUND
                            } else {
                                   // Nouce not found, first time enter, may it two
                                   current_status_of_blocks[block.hash_value] = NOUNCE_NOTFOUND

                            }
                            val = current_status_of_blocks[block.hash_value]
                     }
                     // remove duplicacy
			if (val == INCORRECT_PID) {
				current_block_list.Remove(e)
				continue
			}else if (val == NOUNCE_FOUND){
				current_block_list.Remove(e)
				continue
			}else if (val >= NOUNCE_NOTFOUND){
				current_status_of_blocks[block.hash_value] = NOUNCE_NOTFOUND
			}
			// swap the current and prev whenever new block added
                     prevHash := block.block.PrevHash
                     var tmp string
                     tmp = "<nil>"
                     tmp = *prevHash
                     prev_val := current_status_of_blocks[tmp]
                     if prev_val == OK {
                            self.lockBlocks.RLock()
                            _, ok := self.blocks[tmp]
                            self.lockBlocks.RUnlock()
                            if ok {
                                   current_status_of_blocks[tmp] = NOUNCE_FOUND
                            } else { // twice of prev with nounce not found
                                   current_status_of_blocks[tmp] = 2 * NOUNCE_NOTFOUND
                            }
                            prev_val = current_status_of_blocks[tmp]
                     }
                     switch prev_val {
                     case INCORRECT_PID:
                            current_status_of_blocks[block.hash_value] = INCORRECT_PID
                            self.msg.wake()
                            current_block_list.Remove(e)
                     case NOUNCE_FOUND:
                            flag, error := self.InsertBlock(block)
                            if error == nil {
                                   current_status_of_blocks[block.hash_value] = NOUNCE_FOUND
                                   self.msg.wake()
                                   flip_flag = flip_flag || flag // either add or insert success
                            } else {
                                   current_status_of_blocks[block.hash_value] = INCORRECT_PID
                                   self.msg.wake()
                            }
                            current_block_list.Remove(e)
                     case 3: // blacklist if any
                            self.lockFraud.RLock()
                            _, black := self.fraud_if_any[tmp]
                            self.lockFraud.RUnlock()
                            if black {
                                   current_status_of_blocks[tmp] = INCORRECT_PID
                                   current_status_of_blocks[block.hash_value] = INCORRECT_PID
                                   self.msg.wake()
                                   current_block_list.Remove(e)
                            }
                     default:
                            // Wait for next turn
                     }
              }
              if flip_flag { // received, awake miner
                     miner.msg.wake()
              }
              hash_list := make([]string, 0)
              for k, val := range current_status_of_blocks {
                     if val == 4 {
                            hash_list = append(hash_list, k)
                            current_status_of_blocks[k] = 3
                     }
              }
              if len(hash_list) > 0 { go QueryAllBlocks(hash_list) }
       }
       loggerB.Info("Blockchain start fn completed......")

}

// Transaction valid or not
func (self *BlockChain) VerifyTransaction(in *p_b.Transaction, block *MyChain) (int32, string) {
	inUUID := *in.UUID
	p_id, hash_val := self.VerifyUUID(inUUID, block)
	loggerP2P.Warningf("this is enter %d :: %s", p_id, hash_val)
	if p_id == OK { return p_id, hash_val } else {
		loggerP2P.Warningf("this is enter %d :: %s", p_id, hash_val)
		block := self.blocks[hash_val]
		for _, t := range block.block.Transactions {
			loggerP2P.Warningf("this is enter %v %s %s ", t, *t.UUID,  *in.UUID )
			if *t.UUID == *in.UUID {
				// same
				loggerP2P.Warning("hrer")
				if (*in.FromID == *t.FromID) &&
					(*in.ToID == *t.ToID) &&
					(*in.Value == *t.Value) &&
					(*in.MiningFee == *t.MiningFee) &&
					(*in.Type == *t.Type) { return p_id, hash_val} else {
					loggerP2P.Warning("hrer %d", INCORRECT_PID)
					return INCORRECT_PID, "" // -1 as nothing and null hash
				}
			}
		}
		loggerP2P.Error("Incorrect!!!!!!!!")
		return INCORRECT_PID, INITIALIZE_STRING
	}

}

func (self *BlockChain) VerifyUUID(uuid string, block *MyChain) (int32, string) {
	val, success := self.list_of_uuid[uuid]
	if success {
		for _, hash := range val {
			parent, success := self.blocks[hash]
			if !success {
				loggerP2P.Error("Incorrect!!!!!!!!")
			}
			if self.VerifyParent(block, parent) {
				var p_id int32
				p_id = *parent.block.BlockID
				return p_id, parent.hash_value
			}
		}
	}
	return OK, ""
}
/*
	Check if already found
*/
func (self *BlockChain) VerifyParent(block *MyChain, parent *MyChain) bool {
	if parent.hash_value == INITIAL_HASH { return true }
	var tmp1, tmp2 int32
	tmp1 = *block.block.BlockID
	tmp2 = *parent.block.BlockID
	if tmp1 <= tmp2 { return block.hash_value == parent.hash_value } // same only
	var aim int32
	aim = *parent.block.BlockID
	currentBlock := block
	for {
		// check the list entirely
		for _, tmp := range currentBlock.prev_block {
			var id int32
			id = *tmp.block.BlockID // return when found
			if id == aim { return tmp.hash_value == parent.hash_value } else if id < aim { break }
			currentBlock = tmp
		}
	}

}
