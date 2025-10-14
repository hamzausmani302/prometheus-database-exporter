package datasource

import (
	"github.com/go-gota/gota/dataframe"
)

type QueryType string

const (
	SQLQueryType QueryType = "SQL"
)

/* Query arguments can be interchanged for different datasources
 */
type IQuery interface {
	Kind() QueryType
}

// IDataSource interface defines the methods that any data source implementation must have.
// Can be mocked for testing purposes as well.
type IDataSource interface {
	GetData(query IQuery) dataframe.DataFrame
	Connect() error
	Close() error
}
