package main

import (
	"fmt"
	"goinsertmil/config"
	"goinsertmil/model"
	"goinsertmil/repository"
	"goinsertmil/service"
	"log"
	"math"
	"sync"
	"time"
)

func main() {
	fmt.Println("start")

	start := time.Now()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	wg := new(sync.WaitGroup)
	domainRepo := repository.NewDomainRepository(db)
	importDataService := service.NewImportDataServices(domainRepo, wg)

	csvReader, csvFile, err := service.OpenCsvFile(config.CsvFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer csvFile.Close()

	jobs := make(chan model.Domain)
	go importDataService.DispatchWorkers(jobs)
	importDataService.ReadCsvFilePerLineThenSendToWorker(csvReader, jobs)

	wg.Wait()

	duration := time.Since(start)
	fmt.Println("done in", int(math.Ceil(duration.Seconds())), "seconds")
}
