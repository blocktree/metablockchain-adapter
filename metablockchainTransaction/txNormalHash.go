package metablockchainTransaction

func GetNormalHash(hash string) string {
	alphabetTable := map[string]string {
		"0": "00",
		"1": "01",
		"2": "02",
		"3": "03",
		"4": "04",
		"5": "05",
		"6": "06",
		"7": "07",
		"8": "08",
		"9": "09",
	}

	result := ""

	for i := 0; i < len(hash); i++ {
		ch := string(hash[i])

		newItem, found := alphabetTable[ch]
		if found{
			result = result + newItem
		}
	}

	return result
}