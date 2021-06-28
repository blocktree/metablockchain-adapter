package metablockchainTransaction

import (
	"github.com/blocktree/go-owcrypt"
)

//const calPeriod = 128
const calPeriod = 64

func GetEra(height uint64) []byte {
	calPeriod := uint64(64)
	if CurveType==owcrypt.ECC_CURVE_X25519 {
		calPeriod = uint64(128)
	}

	index := uint64(6)
	if CurveType==owcrypt.ECC_CURVE_X25519 {
		index = uint64(7)
	}

	phase := height % calPeriod
	trailingZero := index - 1

	var encoded uint64
	if trailingZero > 1 {
		encoded = trailingZero
	} else {
		encoded = 1
	}

	if trailingZero < 15 {
		encoded = trailingZero
	} else {
		encoded = 15
	}

	encoded += phase / 1 << 4

	first := byte(encoded >> 8)
	second := byte(encoded & 0xff)

	return []byte{second, first}
}

func GetMortalEra(height uint64) []byte {

	//let calPeriod = Math.pow(2, Math.ceil(Math.log2(period)));
	//calPeriod = Math.min(Math.max(calPeriod, 4), 1 << 16);
	//const phase = current % calPeriod;
	//const quantizeFactor = Math.max(calPeriod >> 12, 1);
	//const quantizedPhase = phase / quantizeFactor * quantizeFactor;

	//function getTrailingZeros(period) {
	//	//	const binary = period.toString(2);
	//	//	let index = 0;
	//	//
	//	//	while (binary[binary.length - 1 - index] === '0') {
	//	//		index++;
	//	//	}
	//	//
	//	//	return index;
	//	//}

	//const period = this.period.toNumber();
	//const phase = this.phase.toNumber();
	//const quantizeFactor = Math.max(period >> 12, 1);
	//const trailingZeros = getTrailingZeros(period);
	//const encoded = Math.min(15, Math.max(1, trailingZeros - 1)) + (phase / quantizeFactor << 4);
	//const first = encoded >> 8;
	//const second = encoded & 0xff;
	//return new Uint8Array([second, first]);

	//calPeriod := 128
	//
	//phase := float64( height % uint64(calPeriod) )
	//quantizeFactor := math.Max(float64(calPeriod>>12), 1)
	//quantizedPhase := phase / quantizeFactor * quantizeFactor
	//
	//quantizeFactor = math.Max(float64(calPeriod >> 12), 1);
	//const trailingZeros = getTrailingZeros(calPeriod);
	//const encoded = Math.min(15, Math.max(1, trailingZeros - 1)) + (phase / quantizeFactor << 4);
	//const first = encoded >> 8;
	//const second = encoded & 0xff;
	//
	//fmt.Println(calPeriod, ",", quantizedPhase )

	return []byte{}
}