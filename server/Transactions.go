package main

import (
	"container/list"
	"sync"
	p_b "../protobuf/go"
	"sort"
	"errors"
)

/*
	Holds the list of UUIDs for pending transactions
	List of pending transactions
	List of UUID for transactions
	The transaction data
	The appropriate channel from p_b.transaction
	Notification for transaction pending
	Notification for transaction acted upon
	temp and data locks

*/
type Transactions struct {
	list_of_uuid_temp   map[string]bool
	temp                *list.List
	list_of_uuid        map[string]bool
	data                *list.List
	transaction_channel chan *p_b.Transaction
	msg                 Messenger
	msg_final           Messenger
	lockTemp sync.RWMutex
	lockData sync.RWMutex
}

func InitializeTransaction() *Transactions {
	var new_transaction Transactions
	new_transaction.temp = list.New()
	new_transaction.data = list.New()
	new_transaction.list_of_uuid_temp = make(map[string]bool)
	new_transaction.list_of_uuid = make(map[string]bool)
	new_transaction.msg.channel = make(chan int, 1)
	new_transaction.transaction_channel = make(chan *p_b.Transaction, MAX_TRANSACTIONS_CHANNELS)
	new_transaction.msg_final.channel = make(chan int, 1)
	return &new_transaction
}

func (self *Transactions) StartTransactions() {
	loggerT.Infof("^^^^^^^^^^^^^^^^^^Transaction starting^^^^^^^^^^^^^^^^^^^^^")
	for {
		// sleep the thread until found
		self.msg.sleep()
		self.lockTemp.Lock()
		tmp := self.temp
		self.temp = list.New()
		all_uuids := make([]string, 0, len(self.list_of_uuid_temp))
		for key,_ := range self.list_of_uuid_temp {
			all_uuids = append(all_uuids, key)
		}
		self.list_of_uuid_temp = make(map[string]bool)
		self.lockTemp.Unlock()
		self.lockData.Lock()
		self.data.PushBackList(tmp)
		for _, key := range all_uuids {
			self.list_of_uuid[key] = true
		}
		self.lockData.Unlock()
		miner.lockLongest.Lock()
		// Passing the current transaction along with the logger to contibue inserting
		miner.InsertNewTransaction(tmp)
		miner.lockLongest.Unlock()
	}
	loggerT.Infof("^^^^^^^^^^^^^^^^^^Transaction found^^^^^^^^^^^^^^^^^^^^^")
}
/*
	Helper function to add the new received transactions, miner task
*/
func (self *Miner) InsertNewTransaction(tmp *list.List) {
	// Simple list iteration
	loggerT.Infof("**************Insert new Transaction********************")
	for element := tmp.Front(); element != nil; element = element.Next() {
		t := element.Value.(*p_b.Transaction)
		// Ignore if found
		if _, exist := self.list_of_uuid[*t.UUID]; exist {
			continue
		}
		self.list_of_mined_transactions = append(self.list_of_mined_transactions, t)
		//logger.Infof("Mined until now ::: %v", self.list_of_mined_transactions)
	}
	//https://golang.org/pkg/sort/ -> simple sorting based by list_of_mined_transactions, considering mining fee
	sort.Sort(self.list_of_mined_transactions)
	self.list_of_mined_transactions = self.longest_chain.ProcessTransactionsList(MAX_TRANSACTIONS_IN_BLOCK, self.list_of_mined_transactions)
	// wake up the thread to proceed ahead
	self.msg.wake()
}

/*
	Helper functions to enable sorting based on struct values
*/

func (self TransactionsList) Len() int { return len(self) }

func (self TransactionsList) Swap(i, j int) { self[i], self[j] = self[j], self[i] }

func (self TransactionsList) Less(i, j int) bool {
	// consider mining fee as the factor for mining fee
	if *self[i].MiningFee != *self[j].MiningFee { return *self[i].MiningFee > *self[j].MiningFee }
	return *self[i].UUID > *self[j].UUID
}

/*
	MyChain appends the new processed transaction, adding up to the
	available blocks, max limit for the number of transactions in a block is default to 40
*/

func (self *MyChain) ProcessTransactionsList(max_trans int, trans_list []*p_b.Transaction) []*p_b.Transaction {
	loggerT.Infof("@@@@@@@@@@@@@@@@@Now the transaction needs to be appended to the MyChain list@@@@@@@@@@@@@")
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
	ifExists := make(map[string]bool)
	// iterate all the data list k, v
	for key, value := range self.data {
		tmp.data[key] = value
	}
	self.lockData.RUnlock()
	ans := make([]*p_b.Transaction, 0, BATCH_TRANSACTIONS)
	for _, t := range trans_list {
		if _, success := ifExists[*t.UUID]; tmp.VerifyTransactionValue(*t.FromID, *t.Value) && !success {
			ifExists[*t.UUID] = true
			ans = append(ans, t)
			if  MAX_TRANSACTIONS_IN_BLOCK == len(ans){
				return ans
			}
			tmp.ModifyBalance(*t.FromID, -*t.Value)
			tmp.ModifyBalance(*t.ToID, *t.Value - *t.MiningFee)
		}
	}
	loggerT.Infof("@@@@@@@@@@@@@@@@@ sending processed list of transactions @@@@@@@@@@@@@")
	return ans
}

/*
	Reduce money from sender, add to receiver after deducting the mining fee
*/
func (self *MyChain) ModifyBalance(id string, amount int32) (int32, error) {
	if _, ok := self.data[id]; !ok { self.data[id] = INITIAL_BALANCE }
	if amount + self.data[id] > 0 {
		self.data[id] += amount
		return self.data[id], nil
	} else {
		err := errors.New("Less Bank Balance")
		return 0.0, err
	}
}
/*
	Helper, checks for sufficient money in balance
*/
func (self *MyChain) VerifyTransactionValue(id string, amount int32) bool {
	loggerChain.Infof("Id %s", id, "amt %d", amount, "json %v", self.data)
	// If already present
	if _, ok := self.data[id]; ok{
		loggerChain.Infof("baby here %d", self.data[id])
		return self.data[id] >= amount
	}
	loggerChain.Infof("baby here %b", INITIAL_BALANCE >= amount)
	return INITIAL_BALANCE >= amount
}

/*
	Adds new transaction to Transactions List, distributed in the system
*/

func (self *Transactions) AddNewTransaction(t *p_b.Transaction) bool {
	loggerP2P.Info("############To add new transaction %v##############", t)
	self.lockTemp.RLock()
	txUUID := *t.UUID
	// First check for the pending list
	_, ok_temp := self.list_of_uuid_temp[txUUID]
	self.lockTemp.RUnlock()
	// Second check for the processed list
	self.lockData.RLock()
	_, ok := self.list_of_uuid[txUUID]
	self.lockData.RUnlock()
	// Both fails, then pushback to the list of pending transaction
	if !(ok_temp || ok) {
		self.lockTemp.Lock()
		self.list_of_uuid_temp[txUUID] = true
		self.temp.PushBack(t)
		self.lockTemp.Unlock()
		select {
		default:
			// do nothing
		case self.transaction_channel <- t:
			loggerP2P.Debug("Thread wakes for transaction")
			self.msg_final.wake() // Transaction successfully added
		}
		self.msg.wake()
		return true
	} else {
		return false
	}
}
