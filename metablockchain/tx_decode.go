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
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blocktree/metablockchain-adapter/metablockchainTransaction"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TransactionDecoder struct {
	openwallet.TransactionDecoderBase
	wm *WalletManager //钱包管理者
}

//NewTransactionDecoder 交易单解析器
func NewTransactionDecoder(wm *WalletManager) *TransactionDecoder {
	decoder := TransactionDecoder{}
	decoder.wm = wm
	return &decoder
}

//CreateRawTransaction 创建交易单
func (decoder *TransactionDecoder) CreateRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.CreateMMUIRawTransaction(wrapper, rawTx)
}

//SignRawTransaction 签名交易单
func (decoder *TransactionDecoder) SignRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.SignMmuiRawTransaction(wrapper, rawTx)
}

//VerifyRawTransaction 验证交易单，验证交易单并返回加入签名后的交易单
func (decoder *TransactionDecoder) VerifyRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	return decoder.VerifyMMUIRawTransaction(wrapper, rawTx)
}

func (decoder *TransactionDecoder) SubmitRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) (*openwallet.Transaction, error) {
	if len(rawTx.RawHex) == 0 {
		return nil, fmt.Errorf("transaction hex is empty")
	}

	if !rawTx.IsCompleted {
		return nil, fmt.Errorf("transaction is not completed validation")
	}

	from := rawTx.Signatures[rawTx.Account.AccountID][0].Address.Address
	nonce := rawTx.Signatures[rawTx.Account.AccountID][0].Nonce
	nonceUint, _ := strconv.ParseUint(nonce[2:], 16, 64)

	decoder.wm.Log.Info("nonce : ", nonceUint, " update from : ", from)

	txid, err := decoder.wm.ApiClient.sendTransaction(rawTx.RawHex)
	if err != nil {
		decoder.wm.UpdateAddressNonce(wrapper, from, 0)
		decoder.wm.Log.Error("Error Tx to send: ", rawTx.RawHex)
		return nil, err
	}

	//交易成功，地址nonce+1并记录到缓存
	newNonce, _ := math.SafeAdd(nonceUint, uint64(1)) //nonce+1
	decoder.wm.UpdateAddressNonce(wrapper, from, newNonce)

	rawTx.TxID = txid
	rawTx.IsSubmit = true

	decimals := int32(6)

	tx := openwallet.Transaction{
		From:       rawTx.TxFrom,
		To:         rawTx.TxTo,
		Amount:     rawTx.TxAmount,
		Coin:       rawTx.Coin,
		TxID:       rawTx.TxID,
		Decimal:    decimals,
		AccountID:  rawTx.Account.AccountID,
		Fees:       rawTx.Fees,
		SubmitTime: time.Now().Unix(),
	}

	tx.WxID = openwallet.GenTransactionWxID(&tx)

	return &tx, nil
}

func (decoder *TransactionDecoder) SignMmuiRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {
	key, err := wrapper.HDKey()
	if err != nil {
		return nil
	}

	keySignatures := rawTx.Signatures[rawTx.Account.AccountID]

	if keySignatures != nil {
		for _, keySignature := range keySignatures {

			childKey, err := key.DerivedKeyWithPath(keySignature.Address.HDPath, keySignature.EccType)
			keyBytes, err := childKey.GetPrivateKeyBytes()
			if err != nil {
				return err
			}

			//签名交易
			///////交易单哈希签名
			signature, err := metablockchainTransaction.SignTransaction(keySignature.Message, keyBytes)
			if err != nil {
				return fmt.Errorf("transaction hash sign failed, unexpected error: %v", err)
			}
			keySignature.Signature = hex.EncodeToString(signature)
		}
	}

	rawTx.Signatures[rawTx.Account.AccountID] = keySignatures

	return nil
}

func (decoder *TransactionDecoder) VerifyMMUIRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	var (
		emptyTrans = rawTx.RawHex
		signature  = ""
	)

	for accountID, keySignatures := range rawTx.Signatures {
		log.Debug("accountID Signatures:", accountID)
		for _, keySignature := range keySignatures {

			signature = keySignature.Signature

			log.Debug("Signature:", keySignature.Signature)
			log.Debug("PublicKey:", keySignature.Address.PublicKey)
		}
	}

	signedTrans, pass := metablockchainTransaction.VerifyAndCombineTransaction(emptyTrans, signature)

	if pass {
		log.Debug("transaction verify passed")
		rawTx.IsCompleted = true
		rawTx.RawHex = signedTrans
	} else {
		log.Debug("transaction verify failed")
		rawTx.IsCompleted = false
	}

	return nil
}

func (decoder *TransactionDecoder) CreateMMUIRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction) error {

	addresses, err := wrapper.GetAddressList(0, -1, "AccountID", rawTx.Account.AccountID)

	if err != nil {
		return err
	}

	if len(addresses) == 0 {
		return openwallet.Errorf(openwallet.ErrAccountNotAddress, "[%s] have not addresses", rawTx.Account.AccountID)
	}

	addressesBalanceList := make([]AddrBalance, 0, len(addresses))

	for i, addr := range addresses {
		address := addr.Address
		if strings.HasPrefix(address, "5") {
			_, err := decoder.wm.Decoder.AddressDecode(address)
			if err != nil {
				continue
			}

			address, err = decoder.wm.ApiClient.getDidByAddress(address)
			if err != nil {
				continue
			}
		}

		balance, err := decoder.wm.ApiClient.getBalance(address)
		if err != nil {
			return err
		}
		nonce, err := decoder.wm.GetAddressNonce(wrapper, balance)
		if err != nil {
			return err
		}
		balance.Nonce = nonce

		balance.index = i
		addressesBalanceList = append(addressesBalanceList, *balance)
	}

	sort.Slice(addressesBalanceList, func(i int, j int) bool {
		return addressesBalanceList[i].Balance.Cmp(addressesBalanceList[j].Balance) >= 0
	})

	fee := uint64(0)

	var amountStr, to string
	for k, v := range rawTx.To {
		to = k
		amountStr = v
		break
	}

	amount := uint64(int64(convertFromAmount(amountStr, decoder.wm.GetDecimal())))

	from := ""
	fromPub := ""
	nonce := uint64(0)
	for _, a := range addressesBalanceList {
		from = a.Address
		fromPub = addresses[a.index].PublicKey
		nonce = a.Nonce
		break
	}

	if from == "" {
		return openwallet.Errorf(openwallet.ErrInsufficientBalanceOfAccount, "the balance: %s is not enough", amountStr)
	}

	nonceMap := map[string]uint64{
		from: nonce,
	}

	rawTx.TxFrom = []string{from}
	rawTx.TxTo = []string{to}
	rawTx.SetExtParam("nonce", nonceMap)
	rawTx.TxAmount = amountStr
	rawTx.Fees = "0" //strconv.FormatUint(fee, 10)	//链上实际收取的，加上0.01的固定消耗
	rawTx.FeeRate = "0"     //strconv.FormatUint(fee, 10)

	memo := rawTx.GetExtParam().Get("memo").String()

	if len(memo)>140 {
		return errors.New("memo length too long")
	}

	mostHeightBlock, err := decoder.wm.ApiClient.getMostHeightBlock()
	if err != nil {
		return err
	}

	toBalance, err := decoder.wm.ApiClient.getBalance( to )
	if err != nil {
		return err
	}

	toPub, err := decoder.wm.Decoder.AddressDecode( toBalance.Address )
	if err != nil {
		return err
	}

	decoder.wm.Log.Debugf("nonce: %d", nonce)

	emptyTrans, message, err := decoder.CreateEmptyRawTransactionAndMessage(fromPub, hex.EncodeToString(toPub), memo, amount, nonce, fee, mostHeightBlock)
	if err != nil {
		return err
	}
	rawTx.RawHex = emptyTrans

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	keySigs := make([]*openwallet.KeySignature, 0)

	addr, err := wrapper.GetAddress(from)
	if err != nil {
		return err
	}
	signature := openwallet.KeySignature{
		EccType: decoder.wm.Config.CurveType,
		Nonce:   "0x" + strconv.FormatUint(nonce, 16),
		Address: addr,
		Message: message,
	}

	keySigs = append(keySigs, &signature)

	rawTx.Signatures[rawTx.Account.AccountID] = keySigs

	rawTx.FeeRate = big.NewInt(int64(fee)).String()

	rawTx.IsBuilt = true

	return nil
}

func (decoder *TransactionDecoder) GetRawTransactionFeeRate() (feeRate string, unit string, err error) {
	rate := uint64(decoder.wm.Config.FixedFee)
	return convertToAmount(rate, decoder.wm.GetDecimal()), "TX", nil
}

//CreateSummaryRawTransaction 创建汇总交易，返回原始交易单数组
func (decoder *TransactionDecoder) CreateSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {
	if sumRawTx.Coin.IsContract {
		return nil, nil
	} else {
		return decoder.CreateSimpleSummaryRawTransaction(wrapper, sumRawTx)
	}
}

func (decoder *TransactionDecoder) CreateSimpleSummaryRawTransaction(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransaction, error) {

	var (
		rawTxArray      = make([]*openwallet.RawTransaction, 0)
		accountID       = sumRawTx.Account.AccountID
		minTransfer     = big.NewInt(int64(convertFromAmount(sumRawTx.MinTransfer, decoder.wm.GetDecimal())))
		retainedBalance = big.NewInt(int64(convertFromAmount(sumRawTx.RetainedBalance, decoder.wm.GetDecimal())))
	)

	if minTransfer.Cmp(retainedBalance) < 0 {
		return nil, fmt.Errorf("mini transfer amount must be greater than address retained balance")
	}

	if !decoder.wm.Config.IgnoreReserve {
		retainedBalance = retainedBalance.Add(retainedBalance, big.NewInt(decoder.wm.Config.ReserveAmount))
	}

	//获取wallet
	addresses, err := wrapper.GetAddressList(sumRawTx.AddressStartIndex, sumRawTx.AddressLimit,
		"AccountID", sumRawTx.Account.AccountID)
	if err != nil {
		return nil, err
	}

	if len(addresses) == 0 {
		return nil, fmt.Errorf("[%s] have not addresses", accountID)
	}

	searchAddrs := make([]string, 0)
	for _, address := range addresses {
		searchAddrs = append(searchAddrs, address.Address)
	}

	addrBalanceArray, err := decoder.wm.Blockscanner.GetBalanceByAddress(searchAddrs...)
	if err != nil {
		return nil, err
	}

	for _, addrBalance := range addrBalanceArray {

		//检查余额是否超过最低转账
		addrBalance_BI := big.NewInt(int64(convertFromAmount(addrBalance.Balance, decoder.wm.GetDecimal())))

		if addrBalance_BI.Cmp(minTransfer) < 0 {
			continue
		}
		//计算汇总数量 = 余额 - 保留余额
		sumAmount_BI := new(big.Int)
		sumAmount_BI.Sub(addrBalance_BI, retainedBalance)

		//this.wm.Log.Debug("sumAmount:", sumAmount)
		//计算手续费
		feeInt := uint64(0)
		fee := big.NewInt(int64(feeInt))

		//减去手续费
		sumAmount_BI.Sub(sumAmount_BI, fee)
		if sumAmount_BI.Cmp(big.NewInt(0)) <= 0 {
			continue
		}
		if sumAmount_BI.Cmp(big.NewInt(decoder.wm.Config.ReserveAmount)) < 0 {
			return nil, errors.New("The summary address [" + sumRawTx.SummaryAddress + "] 保留余额不足!")
		}

		sumAmount := convertToAmount(sumAmount_BI.Uint64(), decoder.wm.GetDecimal())
		fees := convertToAmount(fee.Uint64(), decoder.wm.GetDecimal())

		decoder.wm.Log.Debug(
			"address : ", addrBalance.Address,
			" balance : ", addrBalance.Balance,
			" fees : ", fees,
			" sumAmount : ", sumAmount)

		//创建一笔交易单
		rawTx := &openwallet.RawTransaction{
			Coin:     sumRawTx.Coin,
			Account:  sumRawTx.Account,
			ExtParam: sumRawTx.ExtParam,
			To: map[string]string{
				sumRawTx.SummaryAddress: sumAmount,
			},
			Required: 1,
			FeeRate:  sumRawTx.FeeRate,
		}

		createErr := decoder.createRawTransaction(
			wrapper,
			rawTx,
			addrBalance)
		if createErr != nil {
			return nil, createErr
		}

		//创建成功，添加到队列
		rawTxArray = append(rawTxArray, rawTx)
	}
	return rawTxArray, nil
}

func (decoder *TransactionDecoder) createRawTransaction(wrapper openwallet.WalletDAI, rawTx *openwallet.RawTransaction, addrBalance *openwallet.Balance) error {

	fee := uint64(0)

	var amountStr, to string
	for k, v := range rawTx.To {
		to = k
		amountStr = v
		break
	}

	amount := uint64(convertFromAmount(amountStr, decoder.wm.GetDecimal()))
	from := addrBalance.Address
	fromAddr, err := wrapper.GetAddress(from)
	if err != nil {
		return err
	}

	rawTx.TxFrom = []string{from}
	rawTx.TxTo = []string{to}
	rawTx.TxAmount = amountStr
	rawTx.Fees = "0"
	rawTx.FeeRate = "0"

	address := from
	if strings.HasPrefix(address, "5") {
		_, err := decoder.wm.Decoder.AddressDecode(address)
		if err != nil {
			return err
		}

		address, err = decoder.wm.ApiClient.getDidByAddress(address)
		if err != nil {
			return err
		}
	}

	addrNodeBalance, err := decoder.wm.ApiClient.getBalance(from)
	if err != nil {
		return errors.New("Failed to get nonce when create summay transaction!")
	}
	nonce, err := decoder.wm.GetAddressNonce(wrapper, addrNodeBalance)
	if err != nil {
		return errors.New("Failed to get nonce when create summay transaction!")
	}

	nonceJSON := map[string]interface{}{}
	if len(rawTx.ExtParam) > 0 {
		err = json.Unmarshal([]byte(rawTx.ExtParam), &nonceJSON)
		if err != nil {
			return err
		}
	}
	nonceJSON[from] = nonce

	rawTx.SetExtParam("nonce", nonceJSON)

	mostHeightBlock, err := decoder.wm.ApiClient.getMostHeightBlock()
	if err != nil {
		return errors.New("Failed to get block height when create summay transaction!")
	}

	toPub, err := decoder.wm.Decoder.AddressDecode(to)
	if err != nil {
		return err
	}

	emptyTrans, hash, err := decoder.CreateEmptyRawTransactionAndMessage(fromAddr.PublicKey, hex.EncodeToString(toPub), "", amount, nonce, fee, mostHeightBlock)

	if err != nil {
		return err
	}
	rawTx.RawHex = emptyTrans

	if rawTx.Signatures == nil {
		rawTx.Signatures = make(map[string][]*openwallet.KeySignature)
	}

	keySigs := make([]*openwallet.KeySignature, 0)

	signature := openwallet.KeySignature{
		EccType: decoder.wm.Config.CurveType,
		Nonce:   "0x" + strconv.FormatUint(nonce, 16),
		Address: fromAddr,
		Message: hash,
	}

	keySigs = append(keySigs, &signature)

	rawTx.Signatures[rawTx.Account.AccountID] = keySigs

	rawTx.FeeRate = big.NewInt(int64(fee)).String()

	rawTx.IsBuilt = true

	return nil
}

//CreateSummaryRawTransactionWithError 创建汇总交易，返回能原始交易单数组（包含带错误的原始交易单）
func (decoder *TransactionDecoder) CreateSummaryRawTransactionWithError(wrapper openwallet.WalletDAI, sumRawTx *openwallet.SummaryRawTransaction) ([]*openwallet.RawTransactionWithError, error) {
	raTxWithErr := make([]*openwallet.RawTransactionWithError, 0)
	rawTxs, err := decoder.CreateSummaryRawTransaction(wrapper, sumRawTx)
	if err != nil {
		return nil, err
	}
	for _, tx := range rawTxs {
		raTxWithErr = append(raTxWithErr, &openwallet.RawTransactionWithError{
			RawTx: tx,
			Error: nil,
		})
	}
	return raTxWithErr, nil
}

func (decoder *TransactionDecoder) CreateEmptyRawTransactionAndMessage(fromPub, toPub, memo string, amount uint64, nonce uint64, fee uint64, mostHeightBlock *Block) (string, string, error) {

	txMaterial, err := decoder.wm.ApiClient.getTxMaterial()
	if err != nil {
		return "", "", err
	}
	genesisHash := txMaterial.GenesisHash
	specVersion := txMaterial.SpecVersion
	transactionVersion := txMaterial.TransactionVersion

	//pubkey, _ := hex.DecodeString(fromPub)
	//edpub, _ := owcrypt.CURVE25519_convert_Ed_to_X(pubkey)
	//
	//fromPub = hex.EncodeToString( edpub )

	tx := metablockchainTransaction.TxStruct{
		//发送方公钥
		SenderPubkey: fromPub,
		//接收方公钥
		RecipientPubkey: toPub,
		//发送金额（最小单位）
		Amount: amount,
		//nonce
		Nonce: nonce,
		//手续费（最小单位）
		Fee: 0,
		//备注
		Memo: memo,
		//当前高度
		BlockHeight: mostHeightBlock.Height,
		//当前高度区块哈希
		BlockHash: RemoveOxToAddress(mostHeightBlock.Hash),
		//创世块哈希
		GenesisHash: RemoveOxToAddress(genesisHash),
		//spec版本
		SpecVersion: specVersion,
		//TransactionVersion
		TxVersion : transactionVersion,
	}

	//toAddr, err := decoder.wm.Decoder.AddressEncode( hex.DecodeString(toPub) )
	//if err != nil {
	//	return "", "", err
	//}
	//eraJSONStr, err := decoder.wm.ApiClient.getEra(toAddr, amount, memo, mostHeightBlock.Height)
	//if err != nil {
	//	return "", "", err
	//}
	//eraJSON := gjson.Parse(eraJSONStr)
	//eraFirst := gjson.Get(eraJSON.Raw, "0").Int()
	//eraSecond := gjson.Get(eraJSON.Raw, "1").Int()
	//
	//tx.EraFirst = eraFirst
	//tx.EraSecond = eraSecond

	return tx.CreateEmptyTransactionAndMessage()
}
