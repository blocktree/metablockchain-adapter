package metablockchainTransaction

import "github.com/blocktree/go-owcrypt"

const (
	Generic_Asset_Transfer = "0501"
	SigningBitV4 = byte(0x84)
	AccounntIDFollow = true

	//CurveType = owcrypt.ECC_CURVE_ED25519_NORMAL	 // ed25519_normal

	CurveType = owcrypt.ECC_CURVE_ED25519	 // ed25519
	SuffixOf25519 = byte(0x00)	 // ed25519

	//CurveType = owcrypt.ECC_CURVE_X25519	 // sr25519
	//SuffixOf25519 = byte(0x01)	 // sr25519
)

const  (
	modeBits = 2
	singleMode   byte = 0
	twoByteMode  byte = 1
	fourByteMode byte = 2
	bigIntMode   byte = 3
	singleModeMaxValue   = 63
	twoByteModeMaxValue  = 16383
	fourByteModeMaxValue = 1073741823
)
var modeToNumOfBytes = map[byte]uint{
	singleMode:   1,
	twoByteMode:  2,
	fourByteMode: 4,
}
