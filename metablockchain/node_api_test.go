package cennz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	testNodeAPI = "https://service.eks.centralityapp.com/cennznet-explorer-api" //官方
	//testNodeAPI = "http://8.210.21.7:20003" //local
	symbol = "CENNZ"
)

func PrintJsonLog(t *testing.T, logCont string){
	if strings.HasPrefix(logCont, "{") {
		var str bytes.Buffer
		_ = json.Indent(&str, []byte(logCont), "", "    ")
		t.Logf("Get Call Result return: \n\t%+v\n", str.String())
	}else{
		t.Logf("Get Call Result return: \n\t%+v\n", logCont)
	}
}

func TestGetCall(t *testing.T) {
	tw := NewClient(testNodeAPI, true, symbol)

	if r, err := tw.GetCall("/api/scan/blocks?row=1&page=1" ); err != nil {
		t.Errorf("Get Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

func TestPostCall(t *testing.T) {
	tw := NewClient(testNodeAPI, true, symbol)

	body := map[string]interface{}{
		"address" : "5FHg8oRaXYHWZU5qX5WsKEDjm3QDW8rA1nFS91LwS2bWbvVu",
	}

	if r, err := tw.PostCall("/api/scan/account", body); err != nil {
		t.Errorf("Post Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

func Test_getBlockHeight(t *testing.T) {

	c := NewClient(testNodeAPI, true, symbol)

	r, err := c.getBlockHeight()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("height:", r)
	}
}

func Test_getBalance(t *testing.T) {

	c := NewClient(testNodeAPI, true, symbol)

	address := "xxxx"

	r, err := c.getBalance(address, "")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}

	address = "xxxx"

	r, err = c.getBalance(address, "")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockByHeight(t *testing.T) {
	c := NewClient(testNodeAPI, true, symbol)
	r, err := c.getBlockByHeight(4605510)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockByHeight_1(t *testing.T) {
	c := NewClient(testNodeAPI, false, symbol)
	currentHeight := 4610336
	for i := currentHeight; i < currentHeight+6; i++ {
		block, err := c.getBlockByHeight( uint64(i) )
		if err != nil {
			t.Errorf("GetBlockByHeight failed, err=%v", err)
			return
		}
		t.Log("高度 : ", i, ", 哈希 : ", block.Hash, ", 父哈希 : ", block.PrevBlockHash )
	}
}

func Test_getMostHeightBlock(t *testing.T) {
	c := NewClient(testNodeAPI, true, symbol)
	r, err := c.getMostHeightBlock()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}