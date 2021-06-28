package metablockchain

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestAddressDecoder_AddressEncode(t *testing.T) {
	result := "5H6hdmSfgSSc9fdt9SHijPEmHX5K8fkMyDdvwoFDaQuzgLyr"
	p2pk, _ := hex.DecodeString("86377c388ec1afc558ef40c5edb3b4f7bba1a697b1bb711ece23fc7cdbfe2e1f")
	p2pkAddr, _ := tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )

	result = "5HQGHHrz1DMhCP8UkoqJieTRw7KHTjRjLyCujDvu8fPJm5gu"
	p2pk, _ = hex.DecodeString("ec17fb0bc229cbf6c157632a7a25490dc85e1e9bd2398a00bf619c254429c266")
	p2pkAddr, _ = tw.Decoder.AddressEncode(p2pk)
	fmt.Println("p2pkAddr: ", p2pkAddr, strings.EqualFold(result, p2pkAddr) )
}

func TestAddressDecoder_AddressDecode(t *testing.T) {
	result := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "5GE3ifH1FuRDqk9rJFZ57XF9bkmpZrWk2bZi5ywNaFeKRaLP"
	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )

	result = "ec17fb0bc229cbf6c157632a7a25490dc85e1e9bd2398a00bf619c254429c266"
	p2pkAddr = "5EEmkVpixBg3YdgmUcLj7ScDTr125AuT97BGCbb9aw8tTtk7"
	p2pk, _ = tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr = hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(result, p2pkStr) )
}

func TestAddressDecoder(t *testing.T) {
	pubkey := "deb2e8c58b53dbe18cf9b103d06f8e3a83bff78443ba8750ecc1cb9a31532363"
	p2pkAddr := "162zn6hjYDi5bCeQ75LisY4v994xpyJW3iNR76Ea8VwWrnYJ"

	testP2pkAddr, _ := tw.Decoder.AddressEncode( hex.DecodeString(pubkey) )
	fmt.Println("testP2pkAddr: ", testP2pkAddr, strings.EqualFold(testP2pkAddr, p2pkAddr) )

	p2pk, _ := tw.Decoder.AddressDecode(p2pkAddr)
	p2pkStr := hex.EncodeToString(p2pk)
	fmt.Println("p2pkStr: ", p2pkStr, strings.EqualFold(pubkey, p2pkStr) )
}

func Test_ed25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "5HQGHHrz1DMhCP8UkoqJieTRw7KHTjRjLyCujDvu8fPJm5gu")	//正确
	addressArr = append(addressArr, "5HQGHHrz1DMhCP8UkoqJieTRw7KHTjRjLyCujDvu8fPJm5g4")	//改了最后一位
	addressArr = append(addressArr, "2HQGHHrz1DMhCP8UkoqJieTRw7KHTjRjLyCujDvu8fPJm5gu")	//改了第一位
	addressArr = append(addressArr, "5HQGHHrz1DMhCP8UkoqJheTRw7KHTjRjLyCujDvu8fPJm5gu")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}

func Test_sr25519_AddressVerify_Valid(t *testing.T) {
	addressArr := make([]string, 0)
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGac")	//正确
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGa6")	//改了最后一位
	addressArr = append(addressArr, "h4iKEitKAGicKHJprSmFjfveM6FkoVNoSEiL1bC2tykLUGac")	//改了第一位
	addressArr = append(addressArr, "14iKEitKAGicKHJprSmFjfveM6FkoVNoSdiL1bC2tykLUGac")	//改了中间

	for i := 0; i < len(addressArr); i++ {
		address := addressArr[i]
		valid := tw.Decoder.AddressVerify(address)

		fmt.Println(address, " isvalid : ", valid)
	}
}