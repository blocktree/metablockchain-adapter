package openwtester

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/blocktree/go-owcrypt"
	"github.com/blocktree/metablockchain-adapter/metablockchain"
	"github.com/blocktree/metablockchain-adapter/metablockchainTransaction"
	"github.com/blocktree/openwallet/v2/log"
	"github.com/blocktree/openwallet/v2/openwallet"
	"testing"
	"time"
)

func GetRegistVcHash(fromPub, uid string) (string, string, error) {
	ssidVc := metablockchainTransaction.SsidVc{
		Did : metablockchain.IDENTIFIER_PREFIX+uid,
		PublicKey: fromPub,
	}

	ssidVcJSONBytes, _ := json.Marshal(ssidVc)

	ssidVcHash := owcrypt.Hash(ssidVcJSONBytes, 0, owcrypt.HASH_ALG_SHA256)

	ssidVcNormalHash := metablockchainTransaction.GetNormalHash( hex.EncodeToString(ssidVcHash) )
	result := "0x"+ssidVcNormalHash

	message := string(ssidVcJSONBytes)

	return result, message, nil
}

func SignRegistVcHash(privkey, hash string) (string, error) {

	log.Info("wait sign message : ", hash)

	keyBytes, err := hex.DecodeString(privkey)
	if err!=nil {
		return "", errors.New("wrong private key")
	}

	//签名
	ssidVcNormalHash := hash[2:]
	hashBytes, _ := hex.DecodeString(ssidVcNormalHash)

	signature, _, retCode := owcrypt.Signature(keyBytes, nil, hashBytes, metablockchain.CurveType)
	if retCode != owcrypt.SUCCESS {
		return "", fmt.Errorf("transaction hash sign failed, unexpected error: %v", err)
	}
	return "0x" + hex.EncodeToString(signature), nil
}

func testRegistVcSubmit(uid, publickey, hash, signature string) (error) {
	fmt.Println("uid : ", uid, ", publickey : ", publickey, ", signature : ", signature, ", hash : ", hash)

	//tw := metablockchain.NewClient("http://127.0.0.1:12523", true, "MMUI")

	params := "?uid=" + uid
	params = params + "&publickey=" + publickey
	params = params + "&signature=" + signature
	params = params + "&hash=" + hash

	//if r, err := tw.GetCall("/account/registvc"+params ); err != nil {
	//	log.Errorf("Get Call Result failed: %v\n", err)
	//} else {
	//	log.Info(r.String())
	//}

	return nil
}

func TestRegistVc(t *testing.T) {
	tm := testInitWalletManager()

	walletID := "xxxx"
	accountID := "xxxx"

	uid := "xxxx"

	fromPub := ""
	privkey := ""

	wrapper, err := tm.NewWalletWrapper(testApp, walletID)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	_, err = wrapper.GetAssetsAccountInfo(accountID)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	addresses, err := wrapper.GetAddressList(0, -1, "AccountID", accountID)

	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	if len(addresses) == 0 {
		fmt.Println( openwallet.Errorf(openwallet.ErrAccountNotAddress, "[%s] have not addresses", accountID).Error() )
		return
	}

	address := addresses[0]
	fromPub = "0x" + address.PublicKey

	//解锁钱包
	err = wrapper.UnlockWallet("xxxx", 5*time.Second)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	key, err := wrapper.HDKey()
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	childKey, err := key.DerivedKeyWithPath(address.HDPath, metablockchain.CurveType)
	privKeyBytes, err := childKey.GetPrivateKeyBytes()
	if err != nil {
		fmt.Println( err.Error() )
		return
	}
	privkey = hex.EncodeToString( privKeyBytes )

	hash, _, err := GetRegistVcHash(fromPub, uid)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	signature, err := SignRegistVcHash(privkey, hash)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	err = testRegistVcSubmit(uid, fromPub, hash, signature)
	if err != nil {
		return
	}
}


func TestRegistVcOut(t *testing.T) {
	uid := "xxxx"

	fromPub := "xxxx"

	signature := ""

	hash, message, err := GetRegistVcHash(fromPub, uid)
	if err != nil {
		fmt.Println( err.Error() )
		return
	}

	fmt.Println("fromPub : ", fromPub, "hash : ", hash, "message : ", message, "signature : ", signature)

	err = testRegistVcSubmit(uid, fromPub, hash, signature)
	if err != nil {
		return
	}
}