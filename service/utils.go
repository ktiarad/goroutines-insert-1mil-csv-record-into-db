package service

import (
	"encoding/csv"
	"goinsertmil/config"
	"log"
	"os"
)

var dataHeaders = make([]string, 0)

func OpenCsvFile() (*csv.Reader, *os.File, error) {
	log.Println("=> Open CSV file")

	file, err := os.Open(config.CsvFile)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(file)
	return reader, file, nil
}
