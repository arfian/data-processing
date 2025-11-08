// ============================================
// pkg/csv/reader.go
// ============================================
package csv

import (
	"data-processing/internal/domain"
	"encoding/csv"
	"log"
	"os"
)

type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadCSV(filePath string) ([]*domain.CSVRecord, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(currentDir + filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var csvRecords []*domain.CSVRecord
	// Skip header row
	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 5 {
			continue
		}

		csvRecords = append(csvRecords, &domain.CSVRecord{
			ID:           record[0],
			Name:         record[1],
			Description:  record[2],
			Brand:        record[3],
			Category:     record[4],
			Price:        record[5],
			Currency:     record[6],
			Stock:        record[7],
			Ean:          record[8],
			Color:        record[9],
			Size:         record[10],
			Availability: record[11],
			InternalId:   record[12],
			RowNumber:    i + 1,
		})
	}

	return csvRecords, nil
}
