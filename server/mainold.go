package main
// Learning module, with rest api code
//
//import (
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os"
//	// "strings"
//	"time"
//	"encoding/json"
//  	"github.com/gorilla/mux"
//)
//type nodesList struct {
//	Nodes []string `json:"nodes"`
//}
//var nodesDetails nodesList
//
//type transaction struct {
//	Sender string `json:"sender"`
//	Recipient string `json:"recipient"`
//	Amount int `json:"amount"`
//}
//var newTrans []transaction
//
//type block struct {
//	index  int
//	timestamp string
//	transactions []transaction
//	proof int
//	previous_hash string
//}
//
///*func loadConfig() {
//	conf, err := ioutil.ReadFile("config.json")
//
//	if err != nil {
//		fmt.Print("Can't open config file")
//	}
//
//	var dat map[string]interface{}
//	_ = json.Unmarshal(conf, &dat)
//
//	nserver = int(dat["nservers"].(float64))
//
//	loglv, ok := dat["loglevel"].(string)
//	if !ok {
//		loadLogger(loglevel)
//	} else if lv, err := logging.LogLevel(loglv); err == nil {
//		loadLogger(lv)
//	} else {
//		loadLogger(loglevel)
//	}
//
//	log.Infof("nserver=%d", nserver)
//
//	for i := 1; i <= nserver; i++ {
//		id := fmt.Sprintf("%d", i)
//		server, _ := dat[id].(map[string]interface{})
//		servermap = append(servermap, Server{
//			ip:      server["ip"].(string),
//			port:    server["port"].(string),
//			dataDir: server["dataDir"].(string),
//			id:      i,
//		})
//		log.Infof("%s:%s, %s", servermap[i-1].ip, servermap[i-1].port, servermap[i-1].dataDir)
//	}
//	log.Debugf("len(servermap)=%d", len(servermap))
//}*/
//
//
//// type nodesList  []string
//
//func newBlock(proof int,prevHash string) (block){
//	var temp_block block
//	temp_block.index = 1
//	temp_block.timestamp = time.Now().Format(time.RFC3339)
//	temp_block.transactions = newTrans
//	temp_block.proof = proof
//	temp_block.previous_hash = prevHash
//	return temp_block
//}
//
//func mine(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("inside Mine!!!")
//}
//func getChain(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("inside getChain!!!")
//}
//func newTransaction(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("inside new transaction!!!")
//
//	w.Header().Set("Content-Type", "application/json")
//	b, _ := ioutil.ReadAll(r.Body)
//
//	fmt.Println("json value",string(b))
//	var temp_trans transaction
//	json.Unmarshal(b,&temp_trans)
//
//	newTrans = append(newTrans,temp_trans)
//	fmt.Println("new Transaction is ->",newTrans)
//}
//
//func registerNodes(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	// var n nodesList
//	b, _ := ioutil.ReadAll(r.Body)
//	fmt.Println("json value",string(b))
//	json.Unmarshal(b,&nodesDetails)
//	// nodesDetails := n
//	fmt.Println("nodes is ->",nodesDetails)
//}
//
//func getNodes(w http.ResponseWriter, r *http.Request) {
//	fmt.Println("inside Get Nodes!!!!",nodesDetails)
//	// fmt.Println("")
//
//	// for i:= range len(nodesDetails.Nodes) {
//	// 	fmt.Println("node ->",nodesList[i])
//	// }
//}
//// func DeletePerson(w http.ResponseWriter, r *http.Request) {}
//
//func main() {
//	arg := os.Args[1:]
//	port_number := arg[1]
//	if arg[0] != "-p" {
//		fmt.Println("Incorrect flag variable, exiting....")
//	}
//
//	router := mux.NewRouter()
//	router.HandleFunc("/mine", mine).Methods("GET")
//	// router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
//	router.HandleFunc("/chain", getChain).Methods("GET")
//	router.HandleFunc("/nodes", getNodes).Methods("GET")
//	router.HandleFunc("/transactions/new", newTransaction).Methods("POST")
//	router.HandleFunc("/nodes/register", registerNodes).Methods("POST")
//    // router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
//	string_port := ":" + string(port_number)
//
//	fmt.Println("what->",newBlock(100,"woah"))
//	fmt.Println("Yo, the server runs at ",string_port)
//  log.Fatal(http.ListenAndServe(string_port, router))
//	}