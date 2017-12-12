 /*
Usage
MiniBlockChain 12345 --getbalance
Working CLI for mini block chain
*/

package main

import (
	"fmt"
	"os"
	"github.com/urfave/cli"
	"time"
	"google.golang.org/grpc"
	p_b "../protobuf/go"
	"golang.org/x/net/context"
	"math/rand"
	"strconv"
)

func main() {
	app := cli.NewApp()
	app.Name = "MiniBlockChain"
	app.Usage = "See how user friendly it is!! Just focus on the demo! We will help you out"
	app.Version = "1.00.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Sarita Joshi",
			Email: "joshi.sar@husky.neu.edu",
		},
		{
			Name:  "Akshaya Khare",
			Email: "khare.ak@husky.neu.edu",
		},
		{
			Name:  "Rishab Khandelwal",
			Email: "khandelwal.r@husky.neu.edu",
		},
	}
	app.Copyright = "(c) 2017 SJAKRK"
	app.HelpName = "Mini BlockChain"
	app.Usage = "DEMO Mini BlockChain Implementation"
	app.UsageText = "Demonstrating the available API"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "getbalance",
			Value: "12345",
			Usage: "Which get call to be done",
		},
		cli.StringFlag{
			Name: "transfer",
			//Value : "00000",
			//Value : "12345",
			//Value : 10,
			//Value : 2,
			Usage: "Need 4 arguments",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("Welcome to Mini BlockChain!!! Let's Explore!")
		var Value, Fee int32
		userid := "00000"
		FromID := "00000"
		ToID := "12345"
		Value = 10
		Fee = 2
		fmt.Println("Arg leb", c.NArg())
		if c.NArg() > 0 && c.NArg() == 2 {
			userid = c.Args().Get(0)
			fmt.Println("Fetching balance for user :: ", userid)
			conn, error := grpc.Dial("localhost:5002", grpc.WithInsecure())
			if error != nil {
				fmt.Println("Failed to Establish Connection: ", error)
			}
			defer conn.Close()
			client := p_b.NewBlockChainMinerClient(conn)
			if r, err := client.Get(context.Background(), &p_b.GetRequest{UserID: create_string(userid)}); err != nil {
				fmt.Println("GET Error: ", err)
			} else {
				var val int32
				val = *r.Value
				fmt.Println("The balance for user", userid, "is !", val)
			}
		}
		if (c.NArg() > 0 && c.NArg() == 5) {
			FromID = c.Args().Get(0)
			ToID = c.Args().Get(1)
			Value1 := c.Args().Get(2)
			Fee1 := c.Args().Get(3)
			i64_v, _ := strconv.ParseInt(Value1, 10, 32)
			Value = int32(i64_v)
			i64_f, _ := strconv.ParseInt(Fee1, 10, 32)
			Fee = int32(i64_f)
			fmt.Println("Transfer Invoke from :: ", FromID, "To Id ::", ToID)
			conn, error := grpc.Dial("localhost:5002", grpc.WithInsecure())
			if error != nil {
				fmt.Println("Failed to Establish Connection: ", error)
			}
			defer conn.Close()
			client := p_b.NewBlockChainMinerClient(conn)
			if r, err := client.Transfer(context.Background(), &p_b.Transaction{
				Type:   create_cust(p_b.Transaction_TRANSFER),
				UUID:   create_string(UUID128bit()),
				FromID: create_string(FromID), ToID: create_string(ToID), Value: create(Value), MiningFee: create(Fee)}); err != nil {
				fmt.Println("TRANSFER Error: ", err)
			} else {
				fmt.Println("TRANSFER Return: ", r)
			}

		}

		return nil
	}

	app.Run(os.Args)
}

func create(x int32) *int32 {
	return &x
}

func create_string(x string) *string {
	return &x
}

func create_cust(types p_b.Transaction_Types) *p_b.Transaction_Types {
	return &types
}



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
