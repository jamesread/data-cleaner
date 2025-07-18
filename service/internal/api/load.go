package api

import (
	pb "github.com/jamesread/data-cleaner/gen/grpc/data_cleaner/api/v1"
	"github.com/jamesread/data-cleaner/internal/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Connector interface {
	Connect() error
	Load(dataRows []DataRow, columnMap map[int]string) error
}

type MySQLConnector struct {
	Properties map[string]string

	conn *sql.DB
}

func (c *MySQLConnector) Connect() error {
	var err error

	c.conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.Properties["user"], c.Properties["pass"], c.Properties["host"], c.Properties["port"], c.Properties["database"]))

	if err != nil {
		log.Errorf("Failed to connect to MySQL: %v", err)
	}

	return err
}

func (c *MySQLConnector) Load(dataRows []DataRow, columnMap map[int]string) error {
	stmt, err := c.conn.Prepare("TRUNCATE TABLE " + c.Properties["table"])
	stmt.Exec()

	stmt, err = c.conn.Prepare(prepareStatement(columnMap, c.Properties["table"]))

	if err != nil {
		log.Errorf("Failed to prepare statement: %v", err)
		return err
	}

	for i, row := range dataRows {
		_, err = stmt.Exec(
			row.Get("date"),
			row.Get("description"),
			row.Get("category"),
			row.Get("value"),
			row.Get("balance"),
		)

		log.Infof("Executing row %d: %v", i, row.ToSlice())
	}

	if err != nil {
		log.Errorf("Failed to execute statement: %v", err)
		return err
	}

	return nil
}

func prepareStatement(columnMap map[int]string, tableName string) string {
	sql := fmt.Sprintf("INSERT INTO %v (", tableName)

	cols := make([]string, 0, len(columnMap))

	for _, colName := range columnMap {
		cols = append(cols, colName)
	}

	for i := 0; i < len(columnMap); i++ {
		sql += columnMap[i]

		if i < len(columnMap)-1 {
			sql += ", "
		}
	}

	sql += ") VALUES (?, ?, ?, ?, ?)"
	
	log.Infof("Executing SQL:", sql)

	return sql
}


func (api *EtlApi) Load() *pb.LoadResponse {
	res := &pb.LoadResponse{
	}

	ldconfig := config.GetConfig().Load

	connector := initConnector(ldconfig.Destination, config.GetConfig())
	connector.Connect()

	connector.Load(api.Transform(), ldconfig.ColumnMap)

	return res;
}

func initConnector(connectorName string, cfg *config.Config) Connector {
	props := cfg.Connections[connectorName]

	switch props["type"] {
	case "mysql":
		return NewMySQLConnector(props)
	default:
		 log.Warnf("Unknown connector type: %s", props["type"])
		 return nil
	}
}

func NewMySQLConnector(props map[string]string) *MySQLConnector {
	return &MySQLConnector{
		Properties: props,
	}
}

