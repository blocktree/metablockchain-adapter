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
	"github.com/blocktree/openwallet/v2/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Elastos RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	Symbol      string
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}

func NewClient(url string /*token string,*/, debug bool, symbol string) *Client {
	c := Client{
		BaseURL: url,
		//	AccessToken: token,
		Debug: debug,
	}

	log.Debug("BaseURL : ", url)

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api
	c.Symbol = symbol

	return &c
}

// 用get方法获取内容
func (c *Client) PostCall(path string, v map[string]interface{}) (*gjson.Result, error) {
	if c.Debug {
		log.Debug("Start Request API...")
	}

	r, err := req.Post(c.BaseURL+path, req.BodyJSON(&v))

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Debugf("%+v\n", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	result := resp

	return &result, nil
}

// 用get方法获取内容
func (c *Client) GetCall(path string) (*gjson.Result, error) {

	if c.Debug {
		log.Debug("Start Request API : ", c.BaseURL + path)
	}

	r, err := req.Get(c.BaseURL + path)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Debugf("%+v\n", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())

	result := resp

	return &result, nil
}

// 获取当前最高区块
func (c *Client) getBlockHeight() (uint64, error) {
	block, err := c.getMostHeightBlock()
	if err != nil {
		return 0, err
	}
	return block.Height, nil
}

// 获取地址余额
func (c *Client) getBalance(account string) (*AddrBalance, error) {

	if !strings.HasPrefix(account, IDENTIFIER_PREFIX) {
		return nil, errors.New("wrong did format : " + account)
	}

	resp, err := c.GetCall("/account/balance?account="+account)

	if err != nil {
		return nil, err
	}

	if gjson.Get(resp.Raw, "message").Exists() {
		return nil, errors.New( account + " getBalance api error : " + gjson.Get(resp.Raw, "message").String() )
	}

	address := gjson.Get(resp.Raw, "address").String()

	if gjson.Get(resp.Raw, "nonce").Exists()==false {
		return nil, errors.New(account + " getBalance api error : nonce not found")
	}
	nonce := gjson.Get(resp.Raw, "nonce").Uint()

	addrBalance := AddrBalance{
		Account: account,
		Address: address,
		Nonce: nonce,
		Actived: true,
	}

	if gjson.Get(resp.Raw, "balance").Exists()==false {
		return nil, errors.New(account + " getBalance api error : balance not found")
	}
	balanceJSON := gjson.Get(resp.Raw, "balance")

	if gjson.Get(balanceJSON.Raw, "data").Exists()==false {
		return nil, errors.New(account + " getBalance api error : balance -> data not found")
	}
	dataJSON := gjson.Get(balanceJSON.Raw, "data")

	ok := true
	free, ok := big.NewInt(0).SetString(dataJSON.Get("free").String(), 10)
	if !ok {
		return nil, errors.New(account + " getBalance api error : balance -> data -> free wrong data : " + dataJSON.Get("free").String() )
	}
	addrBalance.Free = free

	feeFrozen, ok := big.NewInt(0).SetString(dataJSON.Get("feeFrozen").String(), 10)
	if !ok {
		return nil, errors.New(account + " getBalance api error : balance -> data -> feeFrozen wrong data : " + dataJSON.Get("feeFrozen").String() )
	}
	miscFrozen, ok := big.NewInt(0).SetString(dataJSON.Get("miscFrozen").String(), 10)
	if !ok {
		return nil, errors.New(account + " getBalance api error : balance -> data -> miscFrozen wrong data : " + dataJSON.Get("miscFrozen").String() )
	}
	freeze := new(big.Int).Add(feeFrozen, miscFrozen)
	addrBalance.Freeze = freeze

	balanceBigInt := new(big.Int)
	addrBalance.Balance = balanceBigInt.Sub(free, freeze)

	return &addrBalance, nil
}

func (c *Client) getBlockByHeight(height uint64) (*Block, error) {
	resp, err := c.GetCall("/blocks/getblock?number=" + strconv.FormatUint(height, 10))

	if err != nil {
		return nil, err
	}

	block, err := NewBlock(resp, c.Symbol)
	if err != nil {
		return nil, err
	}

	return block, nil
}

//获取当前最新高度
func (c *Client) getMostHeightBlock() (*Block, error) {
	resp, err := c.GetCall("/blocks/head")

	if err != nil {
		return nil, err
	}

	block, err := NewBlock(resp, c.Symbol)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (c *Client) GetDidByAddress(address string) (string, error) {
	resp, err := c.GetCall("/account/getDidByAddress?address=" + address)

	if err != nil {
		return "", err
	}

	did := gjson.Get(resp.Raw, "did").String()

	did = RemoveSpecialChar(did)

	return did, nil
}

func (c *Client) GetEra(toaddress string, amount uint64, memo string, blocknumber uint64) (string, error) {
	resp, err := c.GetCall("/transaction/getera?toaddress=" + toaddress + "&amount=" + strconv.FormatUint(amount, 10) + "&memo=" + memo + "&blocknumber=" + strconv.FormatUint(blocknumber, 10) )

	if err != nil {
		return "", err
	}

	era := gjson.Get(resp.Raw, "era").String()

	return era, nil
}

func (c *Client) sendTransaction(rawTx string) (string, error) {
	body := map[string]interface{}{
		"txraw": rawTx,
	}

	//resp, err := c.PostCall("/tx", body)
	resp, err := c.PostCall("/transaction", body)
	if err != nil {
		return "", err
	}

	time.Sleep(time.Duration(1) * time.Second)

	log.Debug("sendTransaction result : ", resp)

	if resp.Get("message").String() != "" && resp.Get("message").String() != "" {
		return "", errors.New("Submit transaction with error: " + resp.Get("message").String() + "," + resp.Get("message").String())
	}

	return resp.Get("hash").String(), nil
}

// 获取当前最高区块
func (c *Client) getTxMaterial() (*TxMaterial, error) {
	resp, err := c.GetCall("/transaction/material")

	if err != nil {
		return nil, err
	}
	return GetTxMaterial(resp), nil
}
