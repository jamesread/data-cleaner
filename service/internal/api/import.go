package api

import (
	"encoding/csv"
	"fmt"
	pb "github.com/jamesread/data-cleaner/gen/grpc"
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

func Import() *pb.ImportResponse {
	log.Infof("Starting new import")

	res := &pb.ImportResponse{
		Issues: make([]*pb.Issue, 0),
	}

	dataRows = []StatementRow{}
	globalIndex = 0

	dir := config.GetConfig().ImportDirectory

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

		res.SourceFiles = append(res.SourceFiles, entry.Name())

		parseFile(path.Join(dir, entry.Name()))
	}

	sort.Sort(ByGlobalIndex(dataRows))
	checkRollingTotal(res)

	res.TotalLines = int64(len(dataRows))
	res.CompletedDate = time.Now().Format(time.RFC3339)

	return res
}

func checkRollingTotal(res *pb.ImportResponse) {
	lastBalance := 0.0
	lastDate := time.Now()

	for _, row := range dataRows {
		if row.Index == 0 {
			lastDate = row.Date
			lastBalance = row.Balance
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
			/*
				issue.LastBalance = lastBalance
				issue.ExpectedBalance = newBalance
				issue.Value = row.Value
				issue.ActualBalance = row.Balance
				issue.LastDate = lastDate.Format("02 January 2006")
				issue.RowDate = row.Date.Format("02 January 2006")
				issue.DateDelta = row.Date.Sub(lastDate).String()
				issue.Diff = diff
			*/

			issue.LocationLineNumber = row.LineNumber
			issue.LocationFilename = row.Filename
			res.Issues = append(res.Issues, issue)
		}

		lastBalance = row.Balance
		lastDate = row.Date
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

func parseLines(lines [][]string, filename string) {
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

				break
			case 1: // Transaction type
			case 2: // Sort Code
			case 3: // Account Number
				break
			case 4: // Transaction Description
				rec.Description = value
			case 5:
				rec.Value = -parseMoney(value)
				break
			case 6:
				if rec.Value != 0 {
					break
				}

				rec.Value = parseMoney(value)
				break
			case 7:
				rec.Balance = parseMoney(value)
				break
			default:
				log.Infof("Field %d: %s", column, value)
			}
		}

		rec.Index = globalIndex
		rec.LineNumber = int64(lineNumber)
		rec.Filename = filename

		dataRows = append(dataRows, rec)
		globalIndex++
	}
}

var globalIndex = 0

var dataRows []StatementRow

type ByGlobalIndex []StatementRow

func (a ByGlobalIndex) Len() int           { return len(a) }
func (a ByGlobalIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGlobalIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }

func parseFile(filepath string) {
	log.Infof("Parsing file: %s", filepath)

	contents, err := os.Open(filepath)

	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	defer contents.Close()

	csvReader := csv.NewReader(contents)
	lines, err := csvReader.ReadAll()

	if err != nil {
		log.Fatalf("failed to read csv: %v", err)
	}

	parseLines(lines, filepath)
}
