package factories

import (
	"fmt"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/sirupsen/logrus"
)

/*
Factory class to initiate the configurations and provide abstraction over the creation of
datasource objects
*/
type DatasourceFactory struct {
	logger *logrus.Logger
	cfg    *config.ApplicationConfig
}

func (dsf *DatasourceFactory) Create(dataSourceConfig config.DataSourceConfig) (datasource.IDataSource, error) {
	if dataSourceConfig.DataSourceType == "SQL" {
		return datasource.NewPostgresDatasource(dsf.logger, dsf.cfg, dataSourceConfig), nil
	}
	return nil, fmt.Errorf("invalid datasource type: %s", dataSourceConfig.DataSourceType)
}

func NewDatasourceFactory(logger *logrus.Logger, cfg *config.ApplicationConfig) *DatasourceFactory {
	return &DatasourceFactory{logger, cfg}
}
