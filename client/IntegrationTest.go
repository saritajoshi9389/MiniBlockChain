package main
/*
Citation: Open source community for BlockChain
*/
import (
	"crypto/rand"
	"flag"
	"fmt"
	p_b "../protobuf/go"
	"google.golang.org/grpc"
	"strings"
	"github.com/golang/protobuf/jsonpb"
	"time"
	"github.com/op/go-logging"
	"os"
)
var pbMarshal = jsonpb.Marshaler{}
var loggerTester = logging.MustGetLogger("IntegrationTest")
func loadLogger(level logging.Level) {
	formatter := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.0000} [%{level:.4s}] %{shortfile: 19.19s} %{shortfunc: 15.15s} %{module: 5.5s}▶ %{color:reset}%{message} ",
	)
	backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), formatter))
	backend.SetLevel(level, "")
	logging.SetBackend(backend)
}
func create_custI(types p_b.Transaction_Types) *p_b.Transaction_Types {
	return &types
}

func createI(x int32) *int32 {
	return &x
}
func create_stringI(x string) *string {
	return &x
}

func main() {
	flag.Parse()
	fmt.Println("Let's start the mock block chain implementation....!")
	loadLogger(logging.INFO)
	ConnectToServersI()
}

func ConnectToServersI() {
	fmt.Println("Connecting to Server in Distributed Environment...")
	conn, error := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if error != nil {
		loggerTester.Errorf("Failed to Establish Connection: %v", error)
	}
	defer conn.Close()
	client := p_b.NewBlockChainMinerClient(conn)
	loggerTester.Notice("======================Integration Testing======================================")
	error = IntegrationTest(client)
	loggerTester.Notice("======================Integration Testing Completed ===========================")
	if error == nil{
		loggerTester.Info(" Mini BlockChain Demo completed!!")
	}
}
func IntegrationTest(c p_b.BlockChainMinerClient) ( e error){
	blockID := 0
        t1 := MakeTransaction("00000", "12345", 1000, 1)
        t2 := MakeTransaction("00000", "12345", 10, 2)
        t3 := MakeTransaction("00000", "12345", 8, 2)
        t4 := MakeTransaction("12345", "00000", 5, 1)
        t2list := []*p_b.Transaction{ t2 }
        t3list := []*p_b.Transaction{ t3 }
        initialHash := strings.Repeat("0", 64)
        _, js, getHash, makeerr := MakeBlock(blockID + 1, initialHash, t2list, "Server01")
	loggerTester.Infof("==============Block 1 created:: Block id ->  %d and server node %s============ " , blockID + 1, "Server01")
        if makeerr != nil { return makeerr }
	blockID += 1
	// next hash value
        initialHash = getHash
        queryHash := getHash
        js2 := js
        _, js2, getHash, makeerr = MakeBlock(blockID + 1, initialHash, t3list, "Server02")
	loggerTester.Infof("=============Block 2 created:: Block id ->  %d and server node %s============= " , blockID + 1, "Server02")
        if makeerr != nil { return makeerr }
        blockID += 1
        initialHash = getHash
        makeerr = PushTransaction(c, t1)
	loggerTester.Infof("===============Push Transaction invoked :: for transaction -> %v=================", t1)
        if makeerr != nil { return makeerr }
	// manual sleep for wait
	time.Sleep(5000 * time.Millisecond)
        makeerr = PushBlock(c, js)
	loggerTester.Infof("========================= Push Block invoked :: for block -> %v=====================", js)
        if makeerr != nil { return makeerr }
        makeerr = PushBlock(c, js2)
	loggerTester.Infof("=========================== Push Block invoked :: for block -> %v======================", js2)
        if makeerr != nil { return makeerr }
        time.Sleep(5000 * time.Millisecond)
	i:= 0
	// verify twice
	loggerTester.Info("=========== Verification Invoked ========================!!!")
	for i < 2{
		time.Sleep(5000 * time.Millisecond)
		makeerr = Verify(c, t2) // t2 check
        if makeerr != nil { return  makeerr }
		i += 1
	}
	loggerTester.Info("=========== Verification completed==========================!!!")
        // And the block should exist.
	loggerTester.Infof("=========== Get block  Invoked for block hash %s ============!!!", queryHash)
	_, makeerr = GetBlock(c, queryHash)
        if makeerr != nil { return makeerr }
	loggerTester.Infof(" Push Transaction invoked :: for transaction -> %v", t4)
        makeerr = PushTransaction(c, t4)
        if makeerr != nil {
            return makeerr
        }
	return e
}


func UUID128bitI() string {
	// u := uuid.NewV4()
	// fmt.Println("hi u", u)
	// return fmt.Sprintf("%x",u)
	u := make([]byte, 16)
	_, _ = rand.Read(u)
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	return fmt.Sprintf("%x", u)
}
