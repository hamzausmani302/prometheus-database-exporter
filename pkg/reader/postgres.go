package reader

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Reader interface {
	Connect() (*sql.DB, error)
	Read(query string) (dataframe.DataFrame, error)
	Close() error
}

type PostgresReader struct {
	ctx              context.Context
	conn             *sql.DB
	Logger           *logrus.Logger
	Host             string
	Port             int
	Username         string
	Password         string
	ConnectionString string
}

// Read data from query from the database engine
func (reader *PostgresReader) Read(query string) (dataframe.DataFrame, error) {
	reader.Logger.Infof("Reading the query - %s", query)
	rows, err := reader.conn.Query(query)
	if err != nil {
		reader.Logger.Error(err)
		return dataframe.DataFrame{}, err
	}
	cols, cerr := rows.Columns()
	if cerr != nil {
		reader.Logger.Error(cerr)
		return dataframe.DataFrame{}, cerr
	}
	defer rows.Close()
	var results []map[string]interface{}

	for rows.Next() {
		// Make a slice of interface{}'s to represent each column
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))

		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err)
		}

		rowMap := make(map[string]interface{})
		for i, col := range cols {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			rowMap[col] = v
		}

		results = append(results, rowMap)
	}
	df := dataframe.LoadMaps(results)
	return df, nil
}

// Connect to the database engine
func (reader *PostgresReader) Connect() (*sql.DB, error) {
	if reader.conn != nil {
		return reader.conn, nil
	}
	conn_string := reader.ConnectionString
	if conn_string == "" {
		reader.Logger.Warn("connection string is empty ,so generating one from the info provided")
		conn_string = fmt.Sprintf("postgres://%s:%s@%s:%d/id3_dashboard", reader.Host, reader.Password, reader.Host, reader.Port)
	}
	reader.Logger.Info("conn_string - ", conn_string)
	reader.ctx = context.Background()
	conn, err := sql.Open("postgres", conn_string)
	if err != nil {
		return nil, err
	}
	reader.conn = conn
	return conn, nil
}

// Close Database connection
func (reader *PostgresReader) Close() error {
	if reader.conn != nil {
		err := reader.conn.Close()
		return err
	}
	return nil
}
