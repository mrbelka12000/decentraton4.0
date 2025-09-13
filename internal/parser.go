package internal

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"decentraton/internal/models"
)

func ReadCsvFile(dirPath string) ([]models.CSVData, error) {
	var data []models.CSVData
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("read directory: %w", err)
	}

	for _, file := range dir {

		f, err := os.Open(strings.Join([]string{dirPath, file.Name()}, "/"))
		if err != nil {
			return nil, err
		}

		csvReader := csv.NewReader(f)

		records, err := csvReader.ReadAll()
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV: %w , IN  %s ", err, file.Name())
		}

		if err := f.Close(); err != nil {
			return nil, fmt.Errorf("failed to close file: %w", err)
		}

		if len(records) > 0 {
			records = records[1:]
		}

		data = append(data, models.CSVData{
			Data:     records,
			FileName: file.Name(),
		})
	}

	return data, nil
}
