package data_read

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type CsvParser interface {
	Parse(fileName string) []CSVRecord
}

type CsvParserImpl struct {
}

type CSVRecord struct {
	Sku            string
	Query          string
	MasterTradeURL string
	Datasheet      string
	CataloguePage  string
}

func (c CsvParserImpl) Parse(fileName string) []CSVRecord {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	reader := csv.NewReader(file)

	var records []CSVRecord

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		record := CSVRecord{
			Sku:            row[0],
			Query:          row[1],
			MasterTradeURL: row[2],
			Datasheet:      row[3],
			CataloguePage:  row[4],
		}

		records = append(records, record)
	}

	fmt.Println("Total of " + strconv.Itoa(len(records)) + " records to process")
	return records
}

func NewCsvParser() CsvParser {
	return &CsvParserImpl{}
}
