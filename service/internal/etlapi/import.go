package api

import (
	"encoding/csv"
	"fmt"
	pb "github.com/jamesread/data-cleaner/gen/data_cleaner/api/v1"
	"github.com/jamesread/data-cleaner/internal/config"
	log "github.com/sirupsen/logrus"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (api *EtlApi) Import() *pb.ImportResponse {
	log.Infof("Starting new import")

	res := &pb.ImportResponse{
		Issues: make([]*pb.Issue, 0),
	}

	api.dataRows = make([]StatementRow, 0)
	api.globalIndex = 0

	dir := config.GetConfig().Extract.ImportDirectory

	entries, err := os.ReadDir(dir)

	if err != nil {
		log.Warnf("failed to read directory: %v", err)

		res.Issues = append(res.Issues, &pb.Issue{
			Description: "Failed to read directory: " + err.Error(),
		})
	}

	for _, entry := range entries {
		if entry.Name() == "Export.csv" {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".csv") {
			continue
		}

		res.SourceFiles = append(res.SourceFiles, api.parseFile(dir, entry.Name()))
	}

	sort.Sort(ByGlobalIndex(api.dataRows))
	api.checkRollingTotal(res)

	res.TotalLines = int64(len(api.dataRows))
	res.CompletedDate = time.Now().Format(time.RFC3339)

	res.Transformations = append(res.Transformations, &pb.Transformation{ Description: fmt.Sprintf("Exact replacements: %v", len(config.GetConfig().Transform.Replacements.Exact))})
	res.Transformations = append(res.Transformations, &pb.Transformation{ Description: fmt.Sprintf("Regex replacements: %v", len(config.GetConfig().Transform.Replacements.Regex))})

	return res
}

func (api *EtlApi) checkRollingTotal(res *pb.ImportResponse) {
	lastBalance := 0.0
	lastDate := time.Now()
	lastFile := "?"
	lastLineNumber := int64(0)

	for _, row := range api.dataRows {
		if row.Index == 0 {
			lastDate = row.Date
			lastBalance = row.Balance
			lastFile = row.Filename
			lastLineNumber = row.LineNumber
			continue
		}

		newBalance := lastBalance + row.Value

		diff := math.Abs(newBalance - row.Balance)

		if diff > 0.001 {
			issue := &pb.Issue{}
			issue.Description = "Rolling total mismatch, possible missing data"

			issue.Expected = append(issue.Expected, &pb.RowAttribute{
				Key: "Balance",
				Val: fmt.Sprintf("%v", newBalance),
			})

			issue.Intermediate = append(issue.Intermediate, &pb.RowAttribute{
				Key: "Value",
				Val: fmt.Sprintf("%v", row.Value),
			})

			issue.Intermediate = append(issue.Intermediate, &pb.RowAttribute{
				Key: "Diff",
				Val: fmt.Sprintf("%v", diff),
			})

			issue.Intermediate = append(issue.Intermediate, &pb.RowAttribute{
				Key: "DateDiff",
				Val: row.Date.Sub(lastDate).String(),
			})

			issue.Actual = append(issue.Actual, &pb.RowAttribute{
				Key: "Balance",
				Val: fmt.Sprintf("%v", row.Balance),
			})

			issue.CurrentLocationLineNumber = row.LineNumber
			issue.CurrentLocationFilename = row.Filename
			issue.LastLocationFilename = lastFile
			issue.LastLocationLineNumber = lastLineNumber
			res.Issues = append(res.Issues, issue)
		}

		lastFile = row.Filename
		lastBalance = row.Balance
		lastDate = row.Date
		lastLineNumber = row.LineNumber
		//		log.Infof("Row: %+v", row)
	}
}

type StatementRow struct {
	Date        time.Time
	Description string
	Value       float64
	Balance     float64
	Index       int
	LineNumber  int64
	Filename    string
}

func parseMoney(value string) float64 {
	if value == "" {
		return 0
	}

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		log.Warnf("failed to parse float: %v", err)
	}

	return v
}

func (api *EtlApi) parseLines(lines [][]string, filename string) {
	var err error

	// Skip the first line and iterate in reverse
	for lineNumber := len(lines) - 1; lineNumber >= 0; lineNumber-- {
		if lineNumber == 0 {
			continue
		}

		line := lines[lineNumber]

		var rec StatementRow

		for column, value := range line {
			switch column {
			case 0:
				rec.Date, err = time.Parse("2006-01-02", value)

				if err != nil {
					log.Warnf("failed to parse date: %v", err)
				}

				continue
			case 1: // Transaction type
			case 2: // Sort Code
			case 3: // Account Number
				continue
			case 4: // Transaction Description
				rec.Description = value
			case 5:
				rec.Value = -parseMoney(value)
			case 6:
				if rec.Value != 0 {
					break
				}

				rec.Value = parseMoney(value)
			case 7:
				rec.Balance = parseMoney(value)
			default:
				log.Infof("Field %d: %s", column, value)
			}
		}

		rec.Index = api.globalIndex
		rec.LineNumber = int64(lineNumber)
		rec.Filename = filename

		api.dataRows = append(api.dataRows, rec)
		api.globalIndex++
	}
}

type ByGlobalIndex []StatementRow

func (a ByGlobalIndex) Len() int           { return len(a) }
func (a ByGlobalIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGlobalIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func (api *EtlApi) parseFile(directory string, filename string) *pb.SourceFile {
	filepath := path.Join(directory, filename)

	log.Infof("Parsing file: %s", filepath)

	contents, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}


	csvReader := csv.NewReader(contents)
	lines, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalf("failed to read csv: %v", err)
	}

	err = contents.Close()

	if err != nil {
		log.Fatalf("failed to close file: %v", err)
	}

	api.parseLines(lines, filepath)

	return &pb.SourceFile{
		Filename: filename,
		LineCount:     int64(len(lines)),
	}
}
