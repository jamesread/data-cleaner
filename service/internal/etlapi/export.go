package api

import (
	"bytes"
	"encoding/csv"
	"github.com/jamesread/data-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
)

type TranslationFile struct {
	Exact map[string]string
	Regex map[string]string
}

func (api *EtlApi) Export(notNullColumns int64) []byte {
	log.Infof("Exporting data, rows: %v", len(api.dataRows))

	config.ReloadConfig()

	buf := &bytes.Buffer{}

	writer := csv.NewWriter(buf)

	for _, row := range api.Transform() {
		err := writer.Write(row.ToSlice())

		if err != nil {
			log.Errorf("Error writing row to CSV: %v", err)
			continue
		}
	}

	writer.Flush()

	return buf.Bytes()
}
