package metablockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

const (
	testNodeAPI = "http://127.0.0.1:12523" //local
	symbol = "MMUI"
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

	if r, err := tw.GetCall("/blocks/head" ); err != nil {
		t.Errorf("Get Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

func TestPostCall(t *testing.T) {
	tw := NewClient(testNodeAPI, true, symbol)

	body := map[string]interface{}{
		"txraw" : "0xabc123",
	}

	if r, err := tw.PostCall("/transaction", body); err != nil {
		t.Errorf("Post Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

//func Test_getBlockHeight(t *testing.T) {
//
//	c := NewClient(testNodeAPI, true, symbol)
//
//	r, err := c.getBlockHeight()
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("height:", r)
//	}
//}
//
func Test_getBalance(t *testing.T) {

	c := NewClient(testNodeAPI, true, symbol)

	address := IDENTIFIER_PREFIX+"xxxx"

	r, err := c.getBalance(address)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockByHeight(t *testing.T) {
	c := NewClient(testNodeAPI, true, symbol)
	r, err := c.getBlockByHeight(2812510)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
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

func Test_getDidByAddress(t *testing.T) {
	c := NewClient(testNodeAPI, true, symbol)
	r, err := c.GetDidByAddress("xxxx")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}