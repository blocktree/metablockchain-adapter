package metablockchainTransaction

import (
	"fmt"
	"testing"
)

func TestGetNormalHash(t *testing.T) {
	hash := "ed997bada4af7bb1902f2d5bcdf1bbd13cbdf437c003e3789a29d8af8630ddc0"

	normalHash := GetNormalHash(hash)

	fmt.Println("hash : ", hash, " => normalHash : 0x"+normalHash)
}
