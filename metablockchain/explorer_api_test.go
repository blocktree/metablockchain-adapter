package cennz

import (
	"testing"
)

const (
	testBalanceAPI = "http://47.57.239.181:20083"
)

func TestGetBalanceApiCall(t *testing.T) {
	tw := NewBalanceClient(testBalanceAPI, true, symbol)

	address := "5E2Y17D8wr59XPoUq5jp6syqP92on6cyjyv9MnX1KewuKTEi"
	blockHash := "0xbff23e0423944a3e82f45ef4f81cd0530feb9e4c49fbde0a9e1d20dd00d11c0f"
	if r, err := tw.BalanceApiGetCall("/account/balance?address=" + address + "&&assetid=1&blockhash=" + blockHash ); err != nil {
		t.Errorf("Get Call Result failed: %v\n", err)
	} else {
		PrintJsonLog(t, r.String())
	}
}

//func Test_getApiBalance(t *testing.T) {
//
//	tw := NewBalanceClient(testBalanceAPI, true, symbol)
//
//	address := ""
//	blockHash := ""
//
//	address = "5E2Y17D8wr59XPoUq5jp6syqP92on6cyjyv9MnX1KewuKTEi"
//	blockHash = "0xbff23e0423944a3e82f45ef4f81cd0530feb9e4c49fbde0a9e1d20dd00d11c0f"
//
//	r, err := tw.getApiBalance(address, blockHash)
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(r)
//	}
//
//	address = "xxxx"
//
//	r, err = c.getBalance(address, "")
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(r)
//	}
//
//	r, err := c.GetFinalizedHead()
//
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println("r:", r)
//	}
//}