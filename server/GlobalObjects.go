package main

import (
	"github.com/op/go-logging"
	"net"
	"fmt"
	"os"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
	p_b "../protobuf/go"
	"runtime"
	"time"
	"github.com/golang/protobuf/jsonpb"
)
/* Needed for the start of mini-blockchain ecosystem*/
var logger = logging.MustGetLogger("main")
var identifier int
var server_list []Server
var count_servers int
var block_chain *BlockChain
var transaction *Transactions
var miner *Miner
// Below three can be tuned for benchmarking
var number_of_block_query_threads = 2
var number_of_push_transactions_threads = 1
var number_of_push_blocks_threads = 1
const TIMEOUT_IN_SEC = 20
var logger1 = logging.MustGetLogger("Server")
var loggerT = logging.MustGetLogger("Start Transaction")
var loggerP2P = logging.MustGetLogger("RPC Communication")
var loggerChain = logging.MustGetLogger("MyChain")
var loggerMiner = logging.MustGetLogger("Miner")
var loggerSlave = logging.MustGetLogger("InformSlaves")
var loggerB = logging.MustGetLogger("Blocks")

const log_level = logging.INFO
const HEIGHT_API = 4
const BLOCK_API = 5
const PUSHB_API = 6
const PUSHT_API = 7
const DEFAULT = 0
const INCORRECT_PID = -1
const OK = 0
const NOUNCE_NOTFOUND = 2
const NOUNCE_FOUND = 1
/*rpc GetHeight(Null) returns (GetHeightResponse) {} -> 4
rpc GetBlock(GetBlockRequest) returns (JsonBlockString) {} -> 5
rpc PushBlock(JsonBlockString) returns (Null) {} -> 6
rpc PushTransaction(Transaction) returns (Null) {} -> 7*/

/*
	Based on the number of CPU,
	we determine the number of threads to be spawn
*/
var number_of_miner = func() int {
	// Mostly half, we tune the system by changing this value
	if runtime.NumCPU() >= 5 {
		return 4
	} else if runtime.NumCPU() >= 4 {
		return 2
	} else {
		return 1 // poor miner, just one thread
	}
}()

/*
 	Initial Block specific constants
*/
const PREV_HASH = "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
const INITIAL_HASH = "0000000000000000000000000000000000000000000000000000000000000000"
const INITIAL_BLOCK_ID = 0
const INITIAL_MINER_ID = ""
const INITIAL_NONCE = "00000000"
const INITIALIZE_STRING = ""
const MAX_TRANSACTIONS_CHANNELS = 10000
const MAX_TRANSACTIONS_IN_BLOCK = 40
const BATCH_TRANSACTIONS = 5
const INITIAL_BALANCE = 5000
const BUFFER_BLOCKS = 10

type TransactionsList []*p_b.Transaction

var pbMarshal = jsonpb.Marshaler{}


func Initialization() {
	/*
		To start the system, we need object of each type, a block, transaction and a miner
	*/
	logger.Infof("---------Start Initialization of Distributed Mini BlockChain!!!------------")
	fmt.Println(time.Now().Format(time.RFC850))
	time.Sleep(10000 * time.Millisecond)

	block_chain = InitializeBlockChain()
	transaction = InitializeTransaction()
	miner = InitializeMiner()
	go block_chain.StartBlocks()
	go transaction.StartTransactions()
	go miner.StartMining()

}

func CreateEcosystem() {
	logger.Infof("-----------Creating the ecosystem, based on the config file.. Wait!!!!------------")
	Initialization()
	StartServers()
}

func StartServers() {

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server_list[identifier-1].server_ip, server_list[identifier-1].server_port))
	if err != nil {
		logger.Fatalf("Fatal!!!!!: %v", err)
		os.Exit(2)
	}
	logger.Infof("-------Up and running %s::%s!!!!!!----------", server_list[identifier-1].server_ip, server_list[identifier-1].server_port)
	// Create gRPC server
	fmt.Println(time.Now().Format(time.RFC850))
	s := grpc.NewServer()
	p_b.RegisterBlockChainMinerServer(s, &P2PServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	// Start server
	if err := s.Serve(lis); err != nil {
		logger.Fatalf("Fatal!!!!!: %v", err)
	}
}


