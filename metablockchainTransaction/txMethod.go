package metablockchainTransaction

import (
	"encoding/hex"
	"errors"
)

type MethodTransfer struct {
	DestPubkey []byte
	Amount     []byte
	Memo    []byte
}

func NewMethodTransfer(pubkey string, amount uint64, memo string) (*MethodTransfer, error) {
	pubBytes, err := hex.DecodeString(pubkey)
	if  err != nil || len(pubBytes) != 32 {
		return nil, errors.New("invalid dest public key")
	}

	if amount == 0 {
		return nil, errors.New("zero amount")
	}
	amountStr := Encode( uint64(amount) )
	if err != nil {
		return nil, errors.New("invalid amount")
	}
	amountBytes, _ := hex.DecodeString(amountStr)

	memoBytes := []byte(memo)

	return &MethodTransfer{
		DestPubkey: pubBytes,
		Amount:     amountBytes,
		Memo:    memoBytes,
	}, nil
}

func (mt MethodTransfer) ToBytes(transferCode string) ([]byte, error) {

	if mt.DestPubkey == nil || len(mt.DestPubkey) != 32 || mt.Amount == nil || len(mt.Amount) == 0 {
		return nil, errors.New("invalid method")
	}

	ret, _ := hex.DecodeString(transferCode)
	if AccounntIDFollow {
		ret = append(ret, 0x00)
	}

	ret = append(ret, mt.DestPubkey...)
	ret = append(ret, mt.Amount...)

	memoLengthBytes := []byte{0x00}

	memoLength := uint64(len(mt.Memo) )
	if memoLength>0 {
		memoLengthStr := Encode( uint64(len(mt.Memo) ) )
		memoLengthBytes, _ = hex.DecodeString(memoLengthStr)
	}

	ret = append(ret, memoLengthBytes...)

	ret = append(ret, mt.Memo...)

	return ret, nil
}