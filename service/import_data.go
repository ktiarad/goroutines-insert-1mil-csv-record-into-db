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

var dataHeaders = make([]string, 0)

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

func (i *ImportDataServices) DispatchWorkers(jobs <-chan model.Domain) {
	for workerIndex := 0; workerIndex <= config.TotalWorker; workerIndex++ {
		go func(workerIndex int, jobs <-chan model.Domain) {
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

func (i *ImportDataServices) ReadCsvFilePerLineThenSendToWorker(csvReader *csv.Reader, jobs chan<- model.Domain) {
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

		// rowOrdered := make([]interface{}, 0)
		// rowOrdered := make([]model.Domain, 0)
		// for _, data := range row {
		// 	rowOrdered = append(rowOrdered, data)
		// }
		rowData := model.Domain{
			GlobalRank:     ToInt(row[0]),
			TldRank:        ToInt(row[1]),
			Domain:         row[2],
			TLD:            row[3],
			RefSubNets:     ToInt(row[4]),
			RefIPs:         ToInt(row[5]),
			IDN_Domain:     row[6],
			IDN_TLD:        row[7],
			PrevGlobalRank: ToInt(row[8]),
			PrevTldRank:    ToInt(row[9]),
			PrevRefSubNets: ToInt(row[10]),
			PrevRefIPs:     ToInt(row[11]),
		}

		i.Wg.Add(1)
		jobs <- rowData
	}

	close(jobs)
}
