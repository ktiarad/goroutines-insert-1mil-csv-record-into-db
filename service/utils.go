package service

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

func OpenCsvFile(fileName string) (*csv.Reader, *os.File, error) {
	log.Println("=> Open CSV file")

	file, err := os.Open(fileName)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	return reader, file, nil
}

func ToInt(data string) int {
	dataInt, err := strconv.Atoi(data)
	if err != nil {
		log.Fatal(err)
	}

	return dataInt
}
