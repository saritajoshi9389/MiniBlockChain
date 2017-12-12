package main

import (
	"sync"
	hashapi "./hash"
	p_b "../protobuf/go"
	// "fmt"
	"container/list"
	"fmt"
	"time"
)
type Miner struct {
	longest_chain              *MyChain
	list_of_mined_transactions TransactionsList
	list_of_uuid               map[string]int32
	msg                        Messenger
	lockLongest                sync.RWMutex
}

func InitializeMiner() *Miner {
	loggerMiner.Noticef("About to initialize Miner")
	var newMiner Miner

	// Initiat newMiner with existing longest Chain
	newMiner.longest_chain = block_chain.long_chain

	// Existing List of mined transactions and uuid will be empty 
	newMiner.list_of_mined_transactions = TransactionsList{}
	newMiner.list_of_uuid = make(map[string]int32)

	// Initiating Channel for this miner
	newMiner.msg.channel = make(chan int, 1)
	return &newMiner
}

// func create_cust(types p_b.Transaction_Types) *p_b.Transaction_Types {
// 	return &types
// }

func (self *Miner) StartMining() {
	// logger := logging.MustGetLogger("miner")
	loggerMiner.Noticef("Inside Start Mining")

	// An incoming channel which keeps track of all 
	// New transaction based on miner availability
	incomingChannel := make([]chan string, number_of_miner)

	// Maleficient(acts like sleep) provides a spell which makes the miner
	// to sleep until a time when miner is hit with another transaction
	// Then Maleficient wakes it back up, and they live happily ever after
	Maleficient := make([]Messenger, number_of_miner)

	// Outgoing channel which ontains the json after mining, doesn't get 
	// activated till mining successfully finishes.
	outgoingChannel := make(chan string)

	//As many number of channels to be created as the number of miners
	for i := 0; i < number_of_miner; i++ {
		incomingChannel[i] = make(chan string, 1)
		Maleficient[i].channel = make(chan int, 1)

		// Creates a new GoRoutine for mining
		go mineTransaction(incomingChannel[i], outgoingChannel, Maleficient[i], i)
	}

	//
	go self.blockPreparer(outgoingChannel)

	//Update longest chain initially
	block_chain.lockBlocks.RLock()
	self.longest_chain = block_chain.blocks[INITIAL_HASH]
	block_chain.lockBlocks.RUnlock()
	// With this, the miner is ready to mine
	loggerMiner.Info("Miner Ready to Mine, gimme some action:")

	// This for loop cotinuously goes on till the server is killed
	for {
		//Sine miner doesn't find any action, it goes to sleep.
		loggerMiner.Info("Miner sleep state")
		fmt.Println(time.Now().Format(time.RFC850))
		self.msg.sleep()

		// Once the miner wakes up, it updates its data 
		// with existing longest-chain from blocks.go
		self.lockLongest.Lock()
		self.refreshMiner()

		// Checks if list of mined transaction is greater than zero
		// Creates a new temporary block
		if len(self.list_of_mined_transactions) > OK {

			tmp := p_b.Block{
				BlockID:      CreateInt(*self.longest_chain.block.BlockID + int32(1)),
				PrevHash:     CreateString(self.longest_chain.hash_value),
				Transactions: self.list_of_mined_transactions,
				MinerID:      CreateString(fmt.Sprintf("Server%02d", identifier)),
				Nonce:        CreateString(INITIAL_NONCE),
			}

			// Now convert the Block to a json string
			js, error := pbMarshal.MarshalToString(&tmp)
			if error != nil {
				loggerMiner.Infof("Awwww Error!!!!: %v", error)
			}

			// Send this details to incoming channel, to start mining 
			in := string(js[:len(js) - BUFFER_BLOCKS])

			//Based on the number of miners, the transactions will be mined
			for i := 0; i < number_of_miner; i++ {
				select {
				case <-incomingChannel[i]:
				default:
				}
				// Here the actual mining should start
				incomingChannel[i] <- in

				// This is where the for loop logic in start_mining wakes up
				// Because maleficient feels it is the right time
				// Maleficient knows the best
				Maleficient[i].wake()
			}
		} else {

			for i := 0; i < number_of_miner; i++ {
				select {
				case <-incomingChannel[i]:
				default:
				}
				incomingChannel[i] <- INITIALIZE_STRING
			}
		}
		self.lockLongest.Unlock()
	}
}

// Verifies Hash and prepares a block for the list of transactions
// and then pushes it to all 
func (self *Miner) blockPreparer(outchannel chan string) {
	loggerMiner.Info("Inside run miner ")
	for {
		js := <-outchannel
		// Compare the hash to verify the transaction
		if hashapi.CheckHash(hashapi.GetHashString(js)) {
			//Update the P2P data type for sending list of transactions to everyone
			queries := []P2PQuery{{api_type: PUSHB_API, argument: &p_b.JsonBlockString{CreateString(js), []byte(INITIALIZE_STRING)}}}
			// channel with all servers
			serverChannel := make(chan *Server, count_servers + 1)
			// Another cannel with P2P data and number of servers
			resps := make(chan []*P2PAnswer, count_servers + 1)
			// Update channel with servernames
			for i := 0; i < count_servers; i++ {
				if i != identifier - 1 {
					serverChannel <- &(server_list[i])
				}
			}
			close(serverChannel)
			// Gossip across the network 
			for i := 0; i < number_of_push_blocks_threads; i++ {
				go InformSlaves(queries, serverChannel, resps)
			}
			//Add block to the block chain
			if b, error := block_chain.ValidateBlock(&p_b.JsonBlockString{CreateString(js), []byte(INITIALIZE_STRING)}); error == nil {
				block_chain.AddBlock(b, true)
			} else {
				loggerMiner.Error("Oops check error for now")
			}
		} else {
			loggerMiner.Error("Oops check error for now")
		}
	}
}

// Contains the core logic for mining
func mineTransaction(in chan string, out chan string, msg Messenger, minerId int) {
	loggerMiner.Info("Miner Slave start")
	js := INITIALIZE_STRING
	a := OK
	for {
		flip := true
		for flip {
			select {
			// Will wait for inchannel to be heard
			case tmp := <-in:
				js = tmp
			default:
				flip = false
			}
		}
		// waits for js to be filled, meaning a proper mine call
		if js == INITIALIZE_STRING {
			loggerMiner.Info("Miner Sleeping... mine transaction!")
			fmt.Println(time.Now().Format(time.RFC850))
			msg.sleep()
			loggerMiner.Info("Miner awake... mined transaction!")

		} else {
			// Start mining, to find the right hash
			for b := 0; b < MAX_TRANSACTIONS_CHANNELS; b++ {
				var dummy string
				dummy = hashapi.GetHashString(fmt.Sprintf("%s%01d%03d%04d\"}", js, minerId, a, b))
				if hashapi.CheckHash(dummy) {
					// If we find the right hash, add it to the channel, nonce length 8
					ans := js + fmt.Sprintf("%01d%03d%04d\"}", minerId, a, b)
					out <- ans
					js = ""
					break
				}
			}
		}
		a++
		a = a % 1000
	}
}


// Refreshes the current miner with latest changes
func (self *Miner) refreshMiner() bool {
	block_chain.lockBlocks.RLock()
	// If existing chain is same as the new blockchain, well that's it then 
	if self.longest_chain == block_chain.long_chain {
		block_chain.lockBlocks.RUnlock() //CHANGE to defer
		return false
	}
	// temporary long chain
	tmpLongChain := self.longest_chain
	// Verify the longest chains and main longest chain have same parent
	for !block_chain.VerifyParent(block_chain.long_chain, tmpLongChain) {
		// If not, remove the UUID map for tmpLongChain, and update with existing 
		// tmpLongChain with original Blockchain's 
		self.unMapUUID(tmpLongChain)
		tmpLongChain = block_chain.blocks[*tmpLongChain.block.PrevHash]
	}
	tmpBlocklongChain := block_chain.long_chain
	// Compare both the long chains, until we find the latest block
	for tmpLongChain != tmpBlocklongChain {
		self.mapUUID(tmpBlocklongChain)
		tmpBlocklongChain = block_chain.blocks[*tmpBlocklongChain.block.PrevHash]
	}
	// Then update longest chain
	self.longest_chain = block_chain.long_chain
	block_chain.lockBlocks.RUnlock()
	transaction.lockData.Lock()
	// For each transaction in the miner block block, it sees if the transaction
	// is actually there or not, if it exists, it deletes it.
	var next *list.Element
	for ele := transaction.data.Front(); ele != nil; ele = next {
		next = ele.Next()
		t, ok := ele.Value.(*p_b.Transaction)
		if !ok {
			loggerMiner.Error("Errorrrrr!!!!!")
		}
		if val, ok := self.list_of_uuid[*t.UUID]; ok && *self.longest_chain.block.BlockID - val >= BUFFER_BLOCKS {
			transaction.data.Remove(ele)
			delete(transaction.list_of_uuid, *t.UUID)
		}
	}
	//If transaction data too big, it filers it
	if transaction.data.Len() > MAX_TRANSACTIONS_CHANNELS {
		self.longest_chain.MinerTransactionFilter()
	} else {
		loggerMiner.Info(" Wait for the miner to filter the transactions!!!!!!!!")
	}

	// Everything else is added to the list of mined transaction
	self.list_of_mined_transactions = make([]*p_b.Transaction, OK, BATCH_TRANSACTIONS)
	self.InsertNewTransaction(transaction.data) //Rebuild miningTx using pending
	transaction.lockData.Unlock()
	return true
}

// Maps the UUID to the right block
func (self *Miner) mapUUID(block *MyChain) {
	for _, t := range block.block.Transactions {
		self.list_of_uuid[*t.UUID] = *block.block.BlockID
	}
}

//UnMaps the existing UUID from the block
func (self *Miner) unMapUUID(block *MyChain) {
	for _, t := range block.block.Transactions {
		delete(self.list_of_uuid, *t.UUID)
	}
}
