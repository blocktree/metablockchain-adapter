package metablockchainTransaction

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"testing"
)

func Test_MMUI_transaction(t *testing.T) {
	senderPubKey := "xxxx"

	//pubkey, _ := hex.DecodeString(senderPubKey)
	//edpub, _ := owcrypt.CURVE25519_convert_Ed_to_X(pubkey)
	//
	//senderPubKey = hex.EncodeToString( edpub )

	tx := TxStruct{
		//发送方公钥
		SenderPubkey:    senderPubKey,
		//接收方公钥
		RecipientPubkey: "xxxx",
		//发送金额（最小单位）
		Amount:         3087655,
		//资产ID
		Memo:         "1234",
		//nonce
		Nonce:           1,
		//手续费（最小单位）
		Fee:             0,
		Tip:             0,
		//当前高度
		BlockHeight:     3191328,
		//当前高度区块哈希
		BlockHash:       "xxxx",
		//创世块哈希
		GenesisHash:     "xxxx",
		//spec版本
		SpecVersion:     5,
		//Transaction版本
		TxVersion: 1,
	}

	// 创建空交易单和待签消息
	emptyTrans, message, err := tx.CreateEmptyTransactionAndMessage()
	if err != nil {
		t.Error("create failed : ", err)
		return
	}
	fmt.Println("空交易单 ： ", emptyTrans)
	fmt.Println("待签消息 ： ",message)

	// 签名
	prikey, _ := hex.DecodeString("xxxx")
	signature, err := SignTransaction(message, prikey)
	if err != nil {
		t.Error("sign failed")
		return
	}
	fmt.Println("签名结果 ： ", hex.EncodeToString(signature))

	// 验签与交易单合并
	signedTrans, pass := VerifyAndCombineTransaction(emptyTrans, hex.EncodeToString(signature))
	if pass {
		fmt.Println("验签成功")
		fmt.Println("签名交易单 ： ", signedTrans)
	} else {
		t.Error("验签失败")
	}
}


func Test_json(t *testing.T)  {
	ts := TxStruct{
		SenderPubkey:    "123",
		RecipientPubkey: "",
		Amount:          0,
		Nonce:           0,
		Fee:             0,
		BlockHeight:     0,
		BlockHash:       "234",
		GenesisHash:     "345",
		SpecVersion:     0,
	}

	js, _ := json.Marshal(ts)

	fmt.Println(string(js))
}

func Test_decode(t *testing.T) {
	//en, _ := codec.Encode(Compact_U32, uint64(139))
	//fmt.Println(en)
}

func Test_regist_vc(t *testing.T){
	IDENTIFIER_PREFIX := "did:ssid:";
	uid := "x3";

	ssidVc := SsidVc{
		Did : IDENTIFIER_PREFIX+uid,
		PublicKey: "xxxx",
	}

	ssidVcJSONBytes, _ := json.Marshal(ssidVc)

	fmt.Println( string(ssidVcJSONBytes) )
	fmt.Println( ssidVcJSONBytes )

	ssidVcHash := owcrypt.Hash(ssidVcJSONBytes, 0, owcrypt.HASH_ALG_SHA256)
	fmt.Println( hex.EncodeToString(ssidVcHash) )

	ssidVcNormalHash := GetNormalHash( hex.EncodeToString(ssidVcHash) )
	fmt.Println( "0x"+ssidVcNormalHash )

	privateKey := []byte{}
	fmt.Println(privateKey)

	hashBytes, _ := hex.DecodeString(ssidVcNormalHash)

	signature, _, retCode := owcrypt.Signature(privateKey, nil, hashBytes, owcrypt.ECC_CURVE_ED25519_NORMAL)
	if retCode != owcrypt.SUCCESS {
		fmt.Println("sign failed")
	}
	fmt.Println( "0x" + hex.EncodeToString(signature) )
	fmt.Println( "0xxxxx")
}