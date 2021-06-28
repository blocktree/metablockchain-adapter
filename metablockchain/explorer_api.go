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
	"strconv"
)

// A Client is a Elastos RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type ExplorerApiClient struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	Symbol      string
}

func NewExplorerApiClient(url string /*token string,*/, debug bool, symbol string) *ExplorerApiClient {
	c := ExplorerApiClient{
		BaseURL: url,
		//	AccessToken: token,
		Debug: debug,
	}

	log.Debug("Explorer BaseURL : ", url)

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api
	c.Symbol = symbol

	return &c
}

// 用get方法获取内容
func (c *ExplorerApiClient) ExplorerApiGetCall(path string) (*gjson.Result, error) {

	if c.Debug {
		log.Debug("Start Request API...")
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

// 获取地址余额
func (c *ExplorerApiClient) setTransactionStatus(transaction *Transaction) (error) {
	eventIndex := strconv.FormatUint(transaction.BlockHeight, 10) + "-" + strconv.FormatUint(transaction.TxIndex, 10)
	url := "/api/v1/extrinsic/" + eventIndex

	r, err := c.ExplorerApiGetCall(url)
	if err != nil {
		return err
	}

	if gjson.Get(r.Raw, "data").Exists()==false {
		return errors.New(eventIndex + " can not found data")
	}
	dataJSON := gjson.Get(r.Raw, "data")

	if gjson.Get(dataJSON.Raw, "attributes").Exists()==false {
		return errors.New(eventIndex + " can not found data->attributes")
	}
	attributesJSON := gjson.Get(dataJSON.Raw, "attributes")

	if gjson.Get(attributesJSON.Raw, "success").Exists()==false {
		return errors.New(eventIndex + " can not found data->attributes->success")
	}
	success := gjson.Get(attributesJSON.Raw, "success").Int()

	if success==1 {
		transaction.Status = "1"
	}

	return nil
}