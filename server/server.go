package main

import (
	"errors"
	"fmt"
	"time"
	p_b "../protobuf/go"
	"github.com/op/go-logging"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct {
	server_id           int
	server_ip           string
	server_port         string
	data_directory_temp string
}

type P2PQuery struct {
	api_type uint8
	argument interface{}
}

type P2PAnswer struct {
	api_type  uint8
	response  interface{}
	error_msg error
}

/*
	Push transactions to all the peer in the network
	Citation: Stackoverflow and open community for golang on github
*/
func PushTransactions(transaction *p_b.Transaction) bool {
	// We define the api types for peer to peer communication based on numbering as below
	tasks := []P2PQuery{{api_type: PUSHT_API, argument: transaction}}
	// One for each server to broadcast the info, except for self
	manager := make(chan *Server, count_servers+1)
	// Collect all P2P response
	respsChan := make(chan []*P2PAnswer, count_servers+1)

	for i := 0; i < count_servers; i++ {
		// exclude self chan
		if i != identifier-1 {
			manager <- &(server_list[i])
		}
	}
	close(manager)
	for i := 0; i < number_of_push_transactions_threads; i++ {
		go InformSlaves( tasks, manager, respsChan)
	}
	for i := 0; i < count_servers-1; i++ {
		getrepsonse := <-respsChan
		ans := getrepsonse[0]
		if ans.error_msg == nil {
			// If correct
			return true
		}
	}
	// failed while pushing
	return false
}
// Trigger apis and store the block, remove from blacklist or the dishonest track of blocks
func QueryAllBlocks(ids []string) {
	logger1.Infof("********enters the query block function!!!!******************")
	logger := logging.MustGetLogger("QueryBlocks")
	// One for each server to broadcast the info, except for self
	manager := make(chan *Server, count_servers+1)
	// Collect all P2P response
	respsChan := make(chan []*P2PAnswer, count_servers+1)
	// Keeps a track of all blacklists, if any in the P2P system
	blacklist := make(map[string]bool)
	// List of activities
	tasks := make([]P2PQuery, 0, len(ids))
	for _, id := range ids {
		blacklist[id] = false
		// get Blocks
		tasks = append(tasks, P2PQuery{api_type: BLOCK_API, argument: &p_b.GetBlockRequest{CreateString(id), []byte("")}})
	}
	for i := 0; i < count_servers; i++ {
		if i != identifier-1 {
			manager <- &(server_list[i])
		}
	}
	close(manager)
	for i := 0; i < number_of_block_query_threads; i++ {
		go InformSlaves(tasks, manager, respsChan)
	}
	for i := 0; i < count_servers-1; i++ {
		resp := <- respsChan
		for _, res := range resp {
			if res == nil { continue } // found nothing
			if res.error_msg != nil { continue } // no error
			getJson, success := res.response.(*p_b.JsonBlockString)
			if !success { continue }
			if b, error := block_chain.ValidateBlock(getJson); error == nil {
				block_chain.AddBlock(b, true) // add block
				delete(blacklist, b.hash_value) // remove from black/ dishonest list
			} else { logger1.Errorf(" Some error %v" , error) }
		}
		if len(blacklist) == 0 { break }
	}
	if len(blacklist) > 0 {
		block_chain.lockFraud.Lock()
		for key,_ := range blacklist { block_chain.fraud_if_any[key] = true }
		block_chain.lockFraud.Unlock()
	}

	for {
		if _, success := <- manager; !success { break }
	}

	logger.Infof("@@@@@@@@@@@@@@@@@@@@@All APIs done@@@@@@@@@@@@@@@@@@@@@@@@")
	block_chain.msg.wake()
}
// Slave handler, invoked in Miner
func InformSlaves(l []P2PQuery, in chan *Server, out chan []*P2PAnswer) {
	loggerSlave.Infof("Slaves!")
	for {
		server, success := <- in
		if !success { break }
		resp , error := TriggerAPIS(server, l)
		if error != nil {
			loggerSlave.Errorf("Error!! %v ", error)
			out <- nil // either empty based on error
		} else {
			out <- resp // or the response reecived
		}
	}
}
// Iterate for the array of queries received, and trigger those grpc endpoints
func TriggerAPIS(server *Server, l []P2PQuery)([]*P2PAnswer, error){
	client, conn, error := CreateNewClient(server)
	if error != nil {
		return nil, error
	} else { defer conn.Close() }
	resp := make([]*P2PAnswer, 0, len(l))
	for _, task := range l {
		answer := TriggerAPI(client, task)
		resp = append(resp, &answer)
	}
	return resp, nil
}
/*
	Refer grpc connection, refer grpc documentation
*/
func CreateNewClient(server *Server) (p_b.BlockChainMinerClient, *grpc.ClientConn, error) {
	conn, error := grpc.Dial(fmt.Sprintf("%s:%s", server.server_ip, server.server_port),
		grpc.WithTimeout(time.Duration(TIMEOUT_IN_SEC)*time.Second), grpc.WithInsecure())
	if error != nil {
		return nil, nil, error
	}
	client := p_b.NewBlockChainMinerClient(conn)
	return client, conn, error
}
/*

rpc Get(GetRequest) returns (GetResponse) {}  -> 1
rpc Transfer(Transaction) returns (BooleanResponse) {} -> 2
rpc Verify(Transaction) returns (VerifyResponse) {} -> 3
rpc GetHeight(Null) returns (GetHeightResponse) {} -> 4
rpc GetBlock(GetBlockRequest) returns (JsonBlockString) {} -> 5
rpc PushBlock(JsonBlockString) returns (Null) {} -> 6
rpc PushTransaction(Transaction) returns (Null) {} -> 7
const HEIGHT_API = 4
const BLOCK_API = 5
const PUSHB_API = 6
const PUSHT_API = 7
*/
func TriggerAPI(client p_b.BlockChainMinerClient, q P2PQuery) P2PAnswer {
	switch q.api_type {
	case HEIGHT_API:
		argument, success := q.argument.(*p_b.Null)
		if !success {return P2PAnswer{api_type: HEIGHT_API,  response: nil, error_msg: errors.New("Error Height API")}}
		if resp, error := client.GetHeight(context.Background(), argument); error != nil {
			return P2PAnswer{api_type: HEIGHT_API, response: nil, error_msg: error}
		} else {
			return P2PAnswer{api_type: HEIGHT_API, response: resp,error_msg: errors.New("nil")}
		}
	case BLOCK_API:
		argument, success := q.argument.(*p_b.GetBlockRequest)
		if !success {return P2PAnswer{api_type: BLOCK_API, response: nil, error_msg: errors.New("Error Block API")}}
		if resp, error:= client.GetBlock(context.Background(), argument); error != nil {
			return P2PAnswer{api_type: BLOCK_API, response: nil, error_msg:error}
		} else {
			return P2PAnswer{api_type: BLOCK_API,response:  resp, error_msg: errors.New("nil")}
		}
	case PUSHB_API:
		argument, success := q.argument.(*p_b.JsonBlockString)
		if !success {return P2PAnswer{api_type: PUSHB_API, response: nil, error_msg: errors.New("Error Push Block")}}
		if resp, error := client.PushBlock(context.Background(), argument); error != nil {
			return P2PAnswer{api_type: PUSHB_API, response: nil, error_msg:error}
		} else {
			return P2PAnswer{api_type: PUSHB_API, response: resp, error_msg: errors.New("nil")}
		}
	case PUSHT_API:
		argument, success := q.argument.(*p_b.Transaction)
		if !success { return P2PAnswer{api_type: PUSHT_API, response: nil, error_msg: errors.New("Error Push Transaction")}}
		if resp, error := client.PushTransaction(context.Background(), argument); error != nil {
			return P2PAnswer{api_type: PUSHT_API, response: nil, error_msg: error}
		} else {
			return P2PAnswer{api_type: PUSHT_API, response: resp, error_msg: errors.New("nil")}
		}
	}
	return P2PAnswer{api_type: DEFAULT,  response: nil, error_msg: errors.New("API error")}
}
