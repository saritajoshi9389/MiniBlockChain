package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"time"
	p_b "../protobuf/go"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	fmt.Println("Let's start the mock block chain implementation....!")
	ConnectToServers()
}

//var server_ips = func() []string {
//	input_config, error := ioutil.ReadFile("config.json")
//	if error != nil {
//		panic(error)
//	}
//	var data map[string]interface{}
//	error = json.Unmarshal(input_config, &data)
//	if error != nil {
//		panic(error)
//	}
//	count_of_servers := int(data["count"].(float64))
//	fmt.Println("Total servers up and running", count_of_servers)
//	server_ips := make([]string, count_of_servers)
//	for i := 0; i < count_of_servers; i++ {
//		d := data[strconv.Itoa(i+1)].(map[string]interface{})
//		server_ips[i] = fmt.Sprintf("%s:%s", d["server_ip"], d["server_port"])
//		fmt.Println("Server IP::", server_ips[i])
//	}
//	return server_ips
//}()

var (
	Operation = flag.String("T", "", `DB transactions to perform:
	1 (or GET): Show the balance of a given UserID. Require option -user.
	5 (or TRANSFER): Transfer some money from one account to another.Require option -from, -to, -value.
	6 (or VERIFY): Check if the transaction is correct or not
	`)
	UserID = flag.String("user", "abcd", "UserID")
	FromID = flag.String("from", "xyz", "FromID")
	ToID   = flag.String("to", "yolo", "Account")
	Value  = flag.Int("value", 10, "Amount")
	Fee    = flag.Int("fee", 5, "MiningFee")
	port =  flag.String("p", "5001", "specify port to use.  defaults to 5001")
	ip = flag.String("ip", "127.0.0.1", "IP on DS system, default 127.0.0.1")
)

func ConnectToServers() {
	fmt.Println("Connecting to Server in Distributed Environment...")
	// count_of_servers := len(server_ips)
	// clients := make([]p_b.BlockChainMinerClient, count_of_servers)
	// fmt.Println("My structure to store client..", clients)
	// conns := make([]*grpc.ClientConn, count_of_servers)
	// fmt.Println("this is conn", conns)
	// for i := 0; i < 1; i++ {
	// 	clients[i], conns[i] = CreateNewClient(server_ips[i])
	// 	defer conns[i].Close()
	//fmt.Println(*ip, *port)
	s:= *ip +":"+ *port
	//fmt.Println(s)
	conn, error := grpc.Dial(s, grpc.WithInsecure())
	if error != nil {
		log.Fatalf("Failed to Establish Connection: %v", error)
	}
	defer conn.Close()
	client := p_b.NewBlockChainMinerClient(conn)
	// return client, conn
	switch *Operation {
	default:
		fmt.Println("value is", *Operation)
		log.Fatal("Unknown operation.")
	case "GET":
		var id string
		id = *UserID
		fmt.Println("The user entered is", id)
		if r, err := client.Get(context.Background(), &p_b.GetRequest{UserID: UserID}); err != nil {
			log.Printf("GET Error: %v", err)
		} else {
			var val int32
			val = *r.Value
			log.Printf("%s", r)
			log.Printf("Balance Amount ::: %d", val)
		}
	case "TRANSFER":
		log.Printf("Starting to transfer your money hey !!!!")
		fmt.Println(time.Now().Format(time.RFC850))
		var i32 interface{}
		var v int
		v = *Value
		i32 = int32(v)
		var i int32
		i32_tmp := i32.(int32)
		i = int32(i32_tmp)
		/// Same for fee
		var f_i32 interface{}
		var f int
		f = *Fee
		f_i32 = int32(f)
		var i1 int32
		f_i32_tmp := f_i32.(int32)
		i1 = int32(f_i32_tmp)
		if r, err := client.Transfer(context.Background(), &p_b.Transaction{
			Type:   create_cust(p_b.Transaction_TRANSFER),
			UUID:   create_string(UUID128bit()),
			FromID: FromID, ToID: ToID, Value: create(i), MiningFee: create(i1)}); err != nil {
			log.Printf("TRANSFER Error: %v", err)
		} else {
			log.Printf("TRANSFER Return: %s", r)
		}
	case "VERIFY":
		log.Printf("Starting to transfer your money hey !!!!")
		fmt.Println(time.Now().Format(time.RFC850))
		var i32 interface{}
		var v int
		v = *Value
		i32 = int32(v)
		var i int32
		i32_tmp := i32.(int32)
		i = int32(i32_tmp)
		/// Same for fee
		var f_i32 interface{}
		var f int
		f = *Fee
		f_i32 = int32(f)
		var i1 int32
		f_i32_tmp := f_i32.(int32)
		i1 = int32(f_i32_tmp)
		uuid_test := create_string(UUID128bit())
		if r, err := client.Transfer(context.Background(), &p_b.Transaction{
			Type:   create_cust(p_b.Transaction_TRANSFER),
			UUID:   uuid_test,
			FromID: FromID, ToID: ToID, Value: create(i), MiningFee: create(i1)}); err != nil {
			log.Printf("Transfer Error: %v", err)
		} else {
			log.Printf("Transfer Return: %s", r)
		}
		time.Sleep(time.Second * 20)
		if r, err := client.Verify(context.Background(), &p_b.Transaction{
			Type:   create_cust(p_b.Transaction_TRANSFER),
			UUID:   uuid_test,
			FromID: FromID, ToID: ToID, Value: create(i), MiningFee: create(i1)}); err != nil {
			log.Printf("Verification Error: %v", err)
		} else {
			log.Printf("Verification Return: %s", r)
		}


	}
}

/*

./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=6 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00000 --to=12345 --value=10 -fee=2
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=10 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=8 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=6 -fee=1
sleep 17
sleep 7
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=4 -fee=1
./client --p=$1 --ip=localhost -T=TRANSFER --from=00001 --to=12346 --value=10 -fee=2


*/

// }
func UUID128bit() string {
	// u := uuid.NewV4()
	// fmt.Println("hi u", u)
	// return fmt.Sprintf("%x",u)

	u := make([]byte, 16)
	_, _ = rand.Read(u)

	u[6] = (u[6] | 0x40) & 0x4F

	u[8] = (u[8] | 0x80) & 0xBF
	return fmt.Sprintf("%x", u)
}
func CreateNewClient(i string) (p_b.BlockChainMinerClient, *grpc.ClientConn) {
	fmt.Println("Establishing connection for every client on DS...!")
	conn, error := grpc.Dial(i, grpc.WithInsecure())
	if error != nil {
		log.Fatalf("Failed to Establish Connection: %v", error)
	}
	client := p_b.NewBlockChainMinerClient(conn)
	return client, conn
}

func create_cust(types p_b.Transaction_Types) *p_b.Transaction_Types {
	return &types
}

func create(x int32) *int32 {
	return &x
}
func create_string(x string) *string {
	return &x
}
