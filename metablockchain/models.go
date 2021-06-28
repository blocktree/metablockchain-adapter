/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */
package metablockchain

import (
	"errors"
	"fmt"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/shopspring/decimal"
	"github.com/tidwall/gjson"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const IDENTIFIER_PREFIX = "did:ssid:";

type txFeeInfo struct {
	GasLimit *big.Int
	GasPrice *big.Int
	Fee      *big.Int
}

type Metadata struct {
	BlockNum        uint64
	NetworkNode     string
	SpecVersion     uint32
}

type RuntimeVersion struct {
	SpecVersion         uint32
	TransactionVersion  uint32
}

type AddrBalance struct {
	Address string
	Account string
	Balance *big.Int
	Free    *big.Int
	Freeze  *big.Int
	Nonce   uint64
	index   int
	Actived bool
}

type Block struct {
	Hash          string        `json:"block"`         // actually block signature in MMUI chain
	PrevBlockHash string        `json:"previousBlock"` // actually block signature in MMUI chain
	Timestamp     uint64        `json:"timestamp"`
	Height        uint64        `json:"height"`
	Transactions  []Transaction `json:"transactions"`
	Finalized      bool          `json:"finalized"`
}

type Extrinsic struct {
	BlockNum uint64
	BlockTimestamp uint64
	ExtrinsicHash string
	ExtrinsicIndex uint64
	Section string
	Method string
	ToArr       []string //@required 格式："地址":"数量"
	MemoArr        []string //@required 备注
	ToDecArr    []string //@required 格式："地址":"数量(带小数)"
	From        string
	Nonce		uint64
	Fee         string
	Status      string
}

type Transaction struct {
	TxID        string
	TxIndex     uint64
	Fee         uint64
	TimeStamp   uint64
	From        string
	To          string
	Amount      uint64
	BlockHeight uint64
	BlockHash   string
	Status      string
	//ToArr       []string //@required 格式："地址":"数量":资产id
	//ToDecArr    []string //@required 格式："地址":"数量(带小数)":资产id
	//FromArr     []string //@required 格式："地址":"数量(带小数)":资产id
	ToTrxDetailArr       []TrxDetail
	FromTrxDetailArr     []TrxDetail
}

type TrxDetail struct {
	Addr        string
	Amount      string
	AmountDec   string
	Memo     string
}

type TxMaterial struct {
	GenesisHash  string
	SpecVersion  uint32
	TransactionVersion  uint32
}

func GetTxMaterial(json *gjson.Result) *TxMaterial {
	obj := &TxMaterial{}

	obj.GenesisHash = gjson.Get(json.Raw, "genesisHash").String()
	obj.SpecVersion = uint32(gjson.Get(json.Raw, "specVersion").Uint())
	obj.TransactionVersion = uint32(gjson.Get(json.Raw, "transactionVersion").Uint())

	return obj
}

func (trx *Transaction) ToDecArr() []string {
	toDecArr := make([]string, 0)

	for _, trxDetail := range trx.ToTrxDetailArr {
		toDecStr := trxDetail.Addr + ":" + trxDetail.AmountDec
		toDecArr = append(toDecArr, toDecStr)
	}

	return toDecArr
}

func GetApiData(json *gjson.Result) (*gjson.Result, error){
	if gjson.Get(json.Raw, "code").Exists()==false {
		return nil, errors.New("api code not found")
	}

	code := gjson.Get(json.Raw, "code").Uint()
	if code!=0 {
		return nil, errors.New("wrong api code " + strconv.FormatUint(code, 10) )
	}

	if gjson.Get(json.Raw, "data").Exists()==false {
		return nil, errors.New("api data not found")
	}

	dataJSON := gjson.Get(json.Raw, "data")

	return &dataJSON, nil
}

func GetMetadata(json *gjson.Result) (*Metadata, error) {
	obj := &Metadata{}

	dataJSON, err :=  GetApiData(json)
	if err!=nil {
		return nil, err
	}

	obj.BlockNum = gjson.Get(dataJSON.Raw, "blockNum").Uint()
	obj.NetworkNode = gjson.Get(dataJSON.Raw, "networkNode").String()
	obj.SpecVersion = uint32(gjson.Get(dataJSON.Raw, "specVersion").Uint())

	return obj, nil
}

func GetRuntimeVersion(json *gjson.Result) (*RuntimeVersion, error) {
	obj := &RuntimeVersion{}

	obj.SpecVersion = uint32( gjson.Get(json.Raw, "specVersion").Uint() )
	obj.TransactionVersion = uint32( gjson.Get(json.Raw, "transactionVersion").Uint() )

	return obj, nil
}

func NewBlock(json *gjson.Result, symbol string) (*Block, error) {
	obj := &Block{}
	// 解析
	obj.Hash = gjson.Get(json.Raw, "hash").String()
	obj.PrevBlockHash = gjson.Get(json.Raw, "parentHash").String()
	obj.Height = gjson.Get(json.Raw, "number").Uint()
	//obj.Timestamp = gjson.Get(json.Raw, "block_timestamp").Uint()
	//obj.Finalized = gjson.Get(json.Raw, "finalized").Bool()
	transactions, err := GetTransactionInBlock(json, symbol)
	if err != nil {
		return nil, err
	}
	obj.Transactions = transactions

	if obj.Hash == "" {
		time.Sleep(5 * time.Second)
	}
	return obj, nil
}

//BlockHeader 区块链头
func (b *Block) BlockHeader() *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	obj.Previousblockhash = b.PrevBlockHash
	obj.Height = b.Height

	return &obj
}

func GetTransactionInBlock(json *gjson.Result, symbol string) ([]Transaction, error) {
	transactions := make([]Transaction, 0)

	blockHash := gjson.Get(json.Raw, "hash").String()
	blockHeight := gjson.Get(json.Raw, "number").Uint()

	blockTime := uint64(time.Now().Unix())

	extrinsicMap := make(map[string]Extrinsic)

	for extrinsicIndex, extrinsicJSON := range gjson.Get(json.Raw, "extrinsics").Array() {
		section := gjson.Get(extrinsicJSON.Raw, "section").String()
		method := gjson.Get(extrinsicJSON.Raw, "method").String()
		//success := gjson.Get(extrinsicJSON.Raw, "success").Bool()
		//finalized := gjson.Get(extrinsicJSON.Raw, "finalized").Bool()
		//paramsStr := gjson.Get(extrinsicJSON.Raw, "params").String()

		txid := gjson.Get(extrinsicJSON.Raw, "hash").String()

		//log.Debug("section : ", section, "method : ", method, ", txid : ", txid)

		//if !success{
		//	continue
		//}

		//获取这个区块的时间
		if section == "timestamp" && method=="set" {
			args := gjson.Get(extrinsicJSON.Raw, "args")
			if len(args.Raw) >0 {
				blockTime = args.Array()[0].Uint()
			}
		}

		isSimpleTransfer := section=="balances" && method=="transferWithMemo"
		if isSimpleTransfer {
			from := gjson.Get(extrinsicJSON.Raw, "senderDid").String()
			to := gjson.Get(extrinsicJSON.Raw, "receiverDid").String()
			amount := gjson.Get(extrinsicJSON.Raw, "amount").String()
			memo := gjson.Get(extrinsicJSON.Raw, "memo").String()
			nonce := gjson.Get(extrinsicJSON.Raw, "nonce").Uint()

			toStr := to + ":" + amount

			toArr := make([]string, 0)
			toArr = append(toArr, toStr)

			memoArr := make([]string, 0)
			memoArr = append(memoArr, memo)
			//fee := gjson.Get(extrinsicJSON.Raw, "fee").String()

			extrinsic := Extrinsic{
				ExtrinsicHash:       txid,
				Section:          	  section,
				Method:			 	  method,
				BlockNum:			  blockHeight,
				ExtrinsicIndex:       uint64(extrinsicIndex),
				BlockTimestamp:		  blockTime,
				ToArr:                toArr,
				ToDecArr:             nil,
				From:                 from,
				Fee:                  "0",
				Status:               "0",
				Nonce:                nonce,
				MemoArr:			  memoArr,
			}

			extrinsicMap[txid] = extrinsic
		}
	}

	for txid, extrinsic := range extrinsicMap {
		toStrArr := strings.Split(extrinsic.ToArr[0], ":")
		if len(toStrArr) !=4 {
			return make([]Transaction, 0), errors.New("wrong txid : " + txid)
		}

		to := IDENTIFIER_PREFIX+toStrArr[ len(toStrArr)-2 ]
		amount := toStrArr[ len(toStrArr)-1 ]
		memo := extrinsic.MemoArr[0]

		amountUint, err := strconv.ParseUint(amount, 10, 64)
		if err != nil {
			return make([]Transaction, 0), errors.New("wrong txid : " + txid + ", wrong amount : " + amount)
		}

		toTrxDetailArr := make([]TrxDetail, 0)
		toTrxDetail := TrxDetail{
			Addr:      to,
			Amount:    amount,
			AmountDec: "",
			Memo:   memo,
		}
		toTrxDetailArr = append(toTrxDetailArr, toTrxDetail)

		fromTrxDetailArr := make([]TrxDetail, 0)
		fromTrxDetail := TrxDetail{
			Addr:      extrinsic.From,
			Amount:    amount,
			AmountDec: "",
			Memo:   memo,
		}
		fromTrxDetailArr = append(fromTrxDetailArr, fromTrxDetail)

		transaction := Transaction{
			TxID:             txid,
			TimeStamp:        extrinsic.BlockTimestamp,
			From:             extrinsic.From,
			To:               to,
			Amount:           amountUint,
			BlockHeight:      blockHeight,
			BlockHash:        blockHash,
			Status:           "0",
			ToTrxDetailArr:   toTrxDetailArr,
			FromTrxDetailArr: fromTrxDetailArr,
			Fee :             0,
			TxIndex:          extrinsic.ExtrinsicIndex,
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

// 从最小单位的 amount 转为带小数点的表示
func convertToAmount(amount uint64, amountDecimal uint64) string {
	amountStr := fmt.Sprintf("%d", amount)
	d, _ := decimal.NewFromString(amountStr)
	ten := math.BigPow(10, int64(amountDecimal) )
	w, _ := decimal.NewFromString(ten.String())

	d = d.Div(w)
	return d.String()
}

// amount 字符串转为最小单位的表示
func convertFromAmount(amountStr string, amountDecimal uint64) uint64 {
	d, _ := decimal.NewFromString(amountStr)
	ten := math.BigPow(10, int64(amountDecimal) )
	w, _ := decimal.NewFromString(ten.String())
	d = d.Mul(w)
	r, _ := strconv.ParseInt(d.String(), 10, 64)
	return uint64(r)
}

func RemoveOxToAddress(addr string) string {
	if strings.Index(addr, "0x") == 0 {
		return addr[2:]
	}
	return addr
}