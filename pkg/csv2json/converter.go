package csv2json

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

func Convert(r io.Reader, w io.Writer) error {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("cannot read csv: %w", err)
	}

	if len(records) == 0 || len(records[0]) == 0 {
		return fmt.Errorf("no CSV headers found")
	}

	headers := records[0]

	encoder := json.NewEncoder(w)
	for i := 1; i < len(records); i++ {
		item := rowToMap(headers, records[i])
		err := encoder.Encode(item)
		if err != nil {
			return fmt.Errorf("cannot encode row %d: %w", i, err)
		}
	}

	return nil
}

func rowToMap(headers, row []string) map[string]string {
	if len(headers) != len(row) {
		log.Fatalf("header size does not match row size\n")
	}

	rowMap := make(map[string]string)

	for i, h := range headers {
		rowMap[h] = row[i]
	}
	return rowMap
}
