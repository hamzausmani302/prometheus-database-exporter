/*
*
Contains reader and all the reader implementations for different databases or other data sources.
All implementations must implement the Reader interface.

reader := PosgresReader{...}
reader.Read()
*
*/
package reader

import (
	"database/sql"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/lib/pq"
)

type Reader interface {
	Connect() (*sql.DB, error)
	Read(query string) (dataframe.DataFrame, error)
	Close() error
}
