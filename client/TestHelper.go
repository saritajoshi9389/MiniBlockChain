package main

import (
	"strings"
	"fmt"
	"strconv"
	"golang.org/x/net/context"
	hash "../server/hash"
	p_b "../protobuf/go"
)

func MakeBlock(id int, ph string, t []*p_b.Transaction, miner string) (b *p_b.Block, js string, h string, e error){
    b = &p_b.Block{
        BlockID: createI(int32(id)),
        PrevHash: create_stringI(ph),
        Transactions: t,
        MinerID: create_stringI(miner),
        Nonce: create_stringI("00000000"),
    }
    json, e := pbMarshal.MarshalToString(b)
    if e != nil {return
    }
    presuf := strings.Split(json, "\"Nonce\":\"00000000\"")
    prefix, suffix := presuf[0], presuf[1]
    prefix = prefix + "\"Nonce\":\""
    suffix = "\"" + suffix
    for i := 0; true; i++ {
        nonce := fmt.Sprintf("%08x", i)
        str := strings.Join([]string{prefix, nonce, suffix}, "")
        h = hash.GetHashString(str)
        succ := hash.CheckHash(h)
	_, err  := strconv.Atoi(nonce)
        if succ && err == nil {
            b.Nonce = create_stringI(nonce)
            js = str
            break
        }
    }
    return
}

func MakeTransaction(from string, to string, val int, fee int) (t *p_b.Transaction) {
    return &p_b.Transaction{
            Type:create_custI(p_b.Transaction_TRANSFER),
            UUID:create_stringI(UUID128bitI()),
            FromID: create_stringI(from), ToID: create_stringI(to), Value: createI(int32(val)), MiningFee: createI(int32(fee))}
}

func Transfer(c p_b.BlockChainMinerClient, from string, to string, val int, fee int) (trans *p_b.Transaction, e error) {
    loggerTester.Noticef(" Transfer from ->  %s  to ::  %s, value is :: %d  and fee is :: %d", from, to, val, fee)
    trans = MakeTransaction(from, to, val, fee)
    r, e := c.Transfer(context.Background(), trans)
    if e != nil {
        loggerTester.Errorf(" Transfer Error!!!!!!! %v", e)
        return
    } else {
        loggerTester.Infof("Transfer Response :: %v", r)
    }
    return
}

func Get(c p_b.BlockChainMinerClient, id string) (e error) {
    loggerTester.Noticef("[GET] %s", id)
    r, e := c.Get(context.Background(), &p_b.GetRequest{UserID: create_stringI(id)})
    if e != nil {
        loggerTester.Errorf("Awww ! here error!  %v", e)
        return
    } else {
        loggerTester.Infof(" Yay! Get response received : %v", r)
    }
    return
}

func Verify(c p_b.BlockChainMinerClient, t *p_b.Transaction) ( e error) {
    loggerTester.Noticef("Verify the given transaction in the Ecosystem :: -> %v", t)
    r, e := c.Verify(context.Background(), t)
    if e != nil {
        loggerTester.Errorf("Oops !!! Error %v", e)
        return
    } else {
        loggerTester.Infof("Verify transaction response :: ->  %v", r)
    }
    return
}

func PushTransaction(c p_b.BlockChainMinerClient, t *p_b.Transaction) (e error) {
    loggerTester.Noticef("Push the received transaction :: -> %v", t)
    _, e = c.PushTransaction(context.Background(), t)
    if e != nil {
        loggerTester.Errorf("Push Transaction Error: %v", e)
        return
    }
    return
}

func PushBlock(c p_b.BlockChainMinerClient, js string) (e error) {
     loggerTester.Noticef(" Push block request received , block value ::  %s", js)
    _, e = c.PushBlock(context.Background(), &p_b.JsonBlockString{Json: create_stringI(js)})
    if e != nil {
        loggerTester.Errorf("Push block Error: %v", e)
        return
    }
    return
}

func GetBlock(c p_b.BlockChainMinerClient, hash string) (js string, e error) {
    loggerTester.Noticef("Getting block with the hash value  %s", hash)
    res, e := c.GetBlock(context.Background(), &p_b.GetBlockRequest{BlockHash: create_stringI(hash)})
    if e != nil {
	    loggerTester.Error("Error!!!!!")
	    return
    } else { var json1 string
	    json1 = *res.Json
	    js = json1
	    loggerTester.Infof("Response received ::  %v", js)
	    return
    }
    return
}
