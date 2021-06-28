package metablockchainTransaction

import (
	"encoding/hex"
	"errors"
	"github.com/blocktree/go-owcrypt"
)



type SignaturePubkey struct {
	Signature []byte
	PublicKey []byte
}

func (ts TxStruct) CreateEmptyTransactionAndMessage() (string, string, error) {

	tp, err := ts.NewTxPayLoad()
	if err != nil {
		return "", "", err
	}

	return ts.ToJSONString(), tp.ToBytesString(), nil
}

func SignTransaction(msgStr string, prikey []byte) ([]byte, error) {
	msg, err := hex.DecodeString(msgStr)
	if err != nil || len(msg) == 0 {
		return nil, errors.New("invalid message to sign")
	}

	if prikey == nil || len(prikey) != 32 {
		return nil, errors.New("invalid private key")
	}

	signature, _, retCode := owcrypt.Signature(prikey, nil, msg, CurveType)
	if retCode != owcrypt.SUCCESS {
		return nil, errors.New("sign failed")
	}

	return signature, nil
}

func VerifyAndCombineTransaction(emptyTrans, signature string) (string, bool) {
	ts, err := NewTxStructFromJSON(emptyTrans)
	if err != nil {
		return "", false
	}

	tp, err := ts.NewTxPayLoad()
	if err != nil {
		return "", false
	}

	msg, _ := hex.DecodeString(tp.ToBytesString())

	pubkey, _ := hex.DecodeString(ts.SenderPubkey)

	sig, err := hex.DecodeString(signature)
	if err != nil || len(sig) != 64{
		return "", false
	}

	if owcrypt.SUCCESS != owcrypt.Verify(pubkey, nil, msg, sig, CurveType) {
		return "", false
	}

	signned, err := ts.GetSignedTransaction(Generic_Asset_Transfer, signature)
	if err != nil {
		return "", false
	}

	return signned, true
}