package datasource

import (
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/hamzausmani302/prometheus-database-exporter/config"

	"github.com/hamzausmani302/prometheus-database-exporter/pkg/reader"
	"github.com/sirupsen/logrus"
)

type SQLQuery struct {
	Query string
	Args  map[string]interface{}
}

func (s SQLQuery) Kind() QueryType {
	return SQLQueryType
}

/*
PostgresDataSource implements the IDataSource interface for PostgreSQL databases.
*/
type PostgresDataSource struct {
	logger *logrus.Logger
	cfg    *config.ApplicationConfig

	dataSourceconfig config.DataSourceConfig
	Reader           *reader.Reader
}

// GetData fetches data from the PostgreSQL database and returns it as a DataFrame.
func (p *PostgresDataSource) GetData(query IQuery) dataframe.DataFrame {
	queryOptions := query.(SQLQuery)
	p.logger.Info("Reading data from Postgres Database")
	df, err := (*p.Reader).Read(queryOptions.Query)
	if err != nil {
		p.logger.Errorf("%s faied with err %s", query, err)
	}
	p.logger.Debug(df)
	return df
}

// Connect establishes a connection to the PostgreSQL database.
func (p *PostgresDataSource) Connect() error {
	p.logger.Info("Connecting to Postgres Database")
	if _, err := (*p.Reader).Connect(); err != nil {
		panic(err)
	}
	return nil
}

// Close closes the connection to the PostgreSQL database.
func (p *PostgresDataSource) Close() error {
	p.logger.Info("Closing connection to Postgres Database")
	(*p.Reader).Close()
	return nil
}

// New creates a new instance of PostgresDataSource.
func NewPostgresDatasource(logger *logrus.Logger, configuration *config.ApplicationConfig, dataSourceConfig config.DataSourceConfig) *PostgresDataSource {
	port, err := strconv.Atoi(dataSourceConfig.Metadata.ConnectionDetails["port"])
	if err != nil {
		port = 5432
	}
	var reader reader.Reader = &reader.PostgresReader{
		Logger:           logger,
		Host:             dataSourceConfig.Metadata.ConnectionDetails["host"],
		Port:             port,
		Username:         dataSourceConfig.Metadata.ConnectionDetails["username"],
		Password:         dataSourceConfig.Metadata.ConnectionDetails["password"],
		ConnectionString: dataSourceConfig.Metadata.ConnectionDetails["connectionString"],
	}
	logger.Info(reader, dataSourceConfig.Metadata.ConnectionDetails)
	ds := PostgresDataSource{logger: logger, cfg: configuration, dataSourceconfig: dataSourceConfig, Reader: &reader}
	if err := ds.Connect(); err != nil {
		logger.Error("Failed to connect to Postgres database", err)
	}
	return &ds
}
