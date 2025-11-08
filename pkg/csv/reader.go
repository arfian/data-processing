// ============================================
// pkg/csv/reader.go
// ============================================
package csv

import (
	"data-processing/internal/domain"
	"encoding/csv"
	"os"
)

type Reader struct{}

func NewReader() *Reader {
	return &Reader{}
}

func (r *Reader) ReadCSV(filePath string) ([]*domain.CSVRecord, error) {
	file, err := os.Open(filePath)
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
			SKU:         record[0],
			Name:        record[1],
			Description: record[2],
			Price:       record[3],
			Stock:       record[4],
			RowNumber:   i + 1,
		})
	}

	return csvRecords, nil
}
