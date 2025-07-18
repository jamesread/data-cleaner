package api

import (
	"regexp"
	"fmt"
	"github.com/jamesread/data-cleaner/internal/config"
)

type DataRow struct {
	contents map[string]string
}

func (r *DataRow) Set(key, value string) {
	r.contents[key] = value
}

func (r *DataRow) Get(key string) string {
	if val, ok := r.contents[key]; ok {
		return val
	}
	return ""
}

func (r *DataRow) ToSlice() []string {
	ret := make([]string, 0)
	ret = append(ret, r.contents["date"])
	ret = append(ret, r.contents["description"])
	ret = append(ret, r.contents["category"])
	ret = append(ret, r.contents["value"])
	ret = append(ret, r.contents["balance"])

	return ret
}

func (api *EtlApi) Transform() []DataRow {
	cfg := config.GetConfig()

	ret := make([]DataRow, 0)

	for _, rec := range api.dataRows {
		category := findCategory(cfg.Transform.Replacements, rec.Description) 

		row := DataRow{
			contents: make(map[string]string),
		}
		row.Set("date", rec.Date.Format("2006-01-02"))
		row.Set("description", rec.Description)
		row.Set("category", category)
		row.Set("value", fmt.Sprintf("%.2f", rec.Value))
		row.Set("balance", fmt.Sprintf("%.2f", rec.Balance))
		ret = append(ret, row)

	}

	return ret
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
