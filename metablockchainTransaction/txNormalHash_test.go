package metablockchainTransaction

import (
	"fmt"
	"testing"
)

func TestGetNormalHash(t *testing.T) {
	hash := "xxxx"

	normalHash := GetNormalHash(hash)

	fmt.Println("hash : ", hash, " => normalHash : 0x"+normalHash)
}
