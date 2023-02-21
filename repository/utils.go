package repository

import (
	"fmt"
	"strings"
)

func GenerateDollarsMark(data []interface{}) string {
	s := make([]string, 0)

	for i := 1; i <= len(data); i++ {
		s = append(s, fmt.Sprintf("$%d", i))
	}

	return strings.Join(s, ",")
}
