/*
    Citation: Blockchain  open community projects and documentation
    Final Project CS5600
    Authors: Akshaya Khare, Rishab Khandelwal, Sarita Joshi
*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/op/go-logging"

)

/*
	Reference:: https://github.com/op/go-logging
	Snippet taken from examples
*/
func loadLogger(level logging.Level) {
	formatter := logging.MustStringFormatter(
		"%{color}%{time:15:04:05.0000} [%{level:.4s}] %{shortfile: 19.19s} %{shortfunc: 15.15s} %{module: 5.5s}▶ %{color:reset}%{message} ",
	)
	backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), formatter))
	backend.SetLevel(level, "")
	logging.SetBackend(backend)
}

func CreateInt(x int32) *int32 {
	return &x
}
func CreateString(x string) *string {
	return &x
}

func CreateBool(x bool) *bool {
	return &x
}

func main() {
	var pid = flag.Int("s_id", 1, "enter the unique id for the server")
	flag.Parse()
	//The load configuration functions that reads the json config file
	cwd, _ := os.Getwd()
	x := cwd + "/*.*"
	fmt.Println("heyya", cwd, x)
	input_config, err := ioutil.ReadFile("config.json")
	//fmt.Println("heyllo", string(input_config))
	//_:= []byte(input_config)
	if err != nil {
		fmt.Print("Error!!!!")
	}
	var data map[string]interface{}
	_ = json.Unmarshal(input_config, &data)

	count_servers = int(data["count"].(float64))
	loadLogger(log_level)
	for i := 1; i <= count_servers; i++ {
		// Unique server id available
		id := fmt.Sprintf("%d", i)
		server, _ := data[id].(map[string]interface{})
		server_list = append(server_list, Server{
			server_id:           i,
			server_ip:           server["server_ip"].(string),
			server_port:         server["server_port"].(string),
			data_directory_temp: server["data_directory_temp"].(string),
		})
	}
	logger.Infof("------------Total servers in the system :: %d--------------", len(server_list))
	identifier = *pid
	if identifier <= 0 {
		logger.Fatal("Fatal!!!!!")
		os.Exit(2)
	}
	if identifier > len(server_list){
		logger.Fatal("Fatal!!!!!")
		os.Exit(2)
	}
	_, err_os := os.Stat(server_list[identifier-1].data_directory_temp)
	if err_os != nil {
		os.MkdirAll(server_list[identifier-1].data_directory_temp, os.ModePerm)
	}
	CreateEcosystem()
}