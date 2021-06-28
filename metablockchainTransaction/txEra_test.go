package metablockchainTransaction

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestGetEra(t *testing.T) {
	height := uint64(3181112)

	era := GetEra(height)

	fmt.Println(hex.EncodeToString(era))

	era = GetMortalEra(height)

	fmt.Println(hex.EncodeToString(era))
}
