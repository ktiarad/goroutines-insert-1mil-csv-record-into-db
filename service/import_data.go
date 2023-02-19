package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"goinsertmil/config"
	"goinsertmil/model"
	"goinsertmil/repository"
	"io"
	"log"
	"sync"
)

func NewImportDataServices(domainRepo repository.DomainRepository, wg *sync.WaitGroup) *ImportDataServices {
	return &ImportDataServices{
		DomainRepo: domainRepo,
		Wg:         wg,
	}
}

type ImportDataServices struct {
	DomainRepo repository.DomainRepository
	Wg         *sync.WaitGroup
}

func (i *ImportDataServices) DispatchWorkers(jobs <-chan []interface{}) {
	for workerIndex := 0; workerIndex <= config.TotalWorker; workerIndex++ {
		go func(workerIndex int, jobs <-chan []interface{}) {
			counter := 0

			for job := range jobs {
				i.ImportData(workerIndex, counter, job)
				i.Wg.Done()
				counter++
			}
		}(workerIndex, jobs)
	}
}

func (i *ImportDataServices) ImportData(workerIndex, counter int, request model.Domain) {
	var outerError error
	for {
		func(outerError *error) {
			defer func() {
				if err := recover(); err != nil {
					*outerError = fmt.Errorf("%v", err)
				}
			}()

			ctx := context.Background()

			err := i.DomainRepo.InsertDomain(ctx, request)
			if err != nil {
				log.Fatal(err.Error())
			}

		}(&outerError)
		if outerError == nil {
			break
		}
	}

	if counter%100 == 0 {
		log.Println("=> worker", workerIndex, "inserted", counter, "data")
	}
}

func (i *ImportDataServices) ReadCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- []interface{}) {
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}

		if len(dataHeaders) == 0 {
			dataHeaders = row
			continue
		}

		rowOrdered := make([]interface{}, 0)
		for _, data := range row {
			rowOrdered = append(rowOrdered, data)
		}

		i.Wg.Add(1)
		jobs <- rowOrdered
	}

	close(jobs)
}
