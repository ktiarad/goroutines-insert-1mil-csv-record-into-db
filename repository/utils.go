package repository

func GenerateDollars(data []interface{}) []string {
	dollarStr := make([]string, 0)

	for i := 1; i <= len(data); i++ {
		dollarStr = append(dollarStr, "?")
	}

	return dollarStr
}
