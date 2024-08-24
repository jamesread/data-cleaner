package api

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/jamesread/data-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type TranslationFile struct {
	Exact map[string]string
	Regex map[string]string
}

func Export() []byte {
	log.Infof("Exporting data")

	config.ReloadConfig()

	cfg := config.GetConfig()

	buf := &bytes.Buffer{}

	writer := csv.NewWriter(buf)

	for _, rec := range dataRows {
		row := []string{
			rec.Date.Format("2006-01-02"),
			rec.Description,
			findCategory(cfg.Replacements, rec.Description),
			fmt.Sprintf("%.2f", rec.Value),
			fmt.Sprintf("%.2f", rec.Balance), // To avoid duplicate lines
		}

		writer.Write(row)
	}

	writer.Flush()

	return buf.Bytes()
}

func findCategory(replacements *config.ReplacementsConfig, description string) string {
	if val, ok := replacements.Exact[description]; ok {
		return val
	}

	for pattern, val := range replacements.Regex {
		if match, _ := regexp.MatchString(pattern, description); match {
			return val
		}
	}

	return ""
}
