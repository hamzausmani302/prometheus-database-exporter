package schema

import (
	"errors"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// struct representing single label object
type Label struct {
	Name        string `yaml:"name"`
	StaticValue string `yaml:"staticValue"`
	ColumnName  string `yaml:"columnName"`
}

// UnmarshalYAML accepts both "columnName" (camelCase) and "column_name" (snake_case)
// so that configs written with either convention work correctly.
func (l *Label) UnmarshalYAML(value *yaml.Node) error {
	type labelAlias struct {
		Name            string `yaml:"name"`
		StaticValue     string `yaml:"staticValue"`
		ColumnName      string `yaml:"columnName"`
		ColumnNameSnake string `yaml:"column_name"`
	}
	var alias labelAlias
	if err := value.Decode(&alias); err != nil {
		return err
	}
	l.Name = alias.Name
	l.StaticValue = alias.StaticValue
	l.ColumnName = alias.ColumnName
	if l.ColumnName == "" {
		l.ColumnName = alias.ColumnNameSnake
	}
	return nil
}

func (l Label) IsStaticValue() bool {
	return l.StaticValue != ""
}

// struct representing single Metric object
type Metric struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Help   string `yaml:"help"`
	Column string `yaml:"column"`
}

// represents the query object for queries defined in config
type Query struct {
	Name             string `yaml:"name"`
	hash             string
	DataSource       string `yaml:"dataSource"`
	dataSource       *datasource.IDataSource
	Query            string   `yaml:"query"`
	Pipeline		 []map[string]interface{} `yaml:"pipeline"`
	QueryTimeout     int      `yaml:"queryTimeout"`
	QueryRefreshTime int      `yaml:"queryRefreshTime"`
	Labels           []Label  `yaml:"labels"`
	Metrics          []Metric `yaml:"metrics"`
}
// query has pipeline or only query
func (query *Query) IsPipeline() bool {
	return len(query.Pipeline) > 0
}

// Set the value of hash from outside
func (query *Query) SetHash(hash string) {
	query.hash = hash
}

// Get the value of Query Hash
func (query *Query) GetHash() string {
	return query.hash
}

/*// Generate hash by the following way
 		MD5(query Name + SQL Query + label names + metrics labels )
// */
func (query *Query) GenerateHash() {

	payload := ""
	for _, label := range query.Labels {
		payload += label.Name
		payload += label.ColumnName
		payload += label.StaticValue
	}
	for _, metric := range query.Metrics {
		payload += metric.Name
		payload += metric.Type
	}
	payload += query.Query
	payload += query.Name
	query.hash = utils.Hash(payload)
}

func (query *Query) GetDataSource() *datasource.IDataSource {
	return query.dataSource
}

// Convert Yaml data to query object
func (query *Query) Load(logger *logrus.Logger, queryData map[string]interface{}, dataSources map[string]datasource.IDataSource) error {
	//parse content into bytes first
	content, err := yaml.Marshal(queryData)
	if err != nil {
		logger.Error("Error marshalling query into bytes")
		return err
	}

	err = yaml.Unmarshal(content, query)
	if err != nil {
		logger.Error("Error Unmshalling for ", string(content), err)
		return err
	}
	// assign datasource
	logger.Debugf("loaded query %q targeting datasource %q", query.Name, query.DataSource)
	ds, ok := dataSources[query.DataSource]
	if !ok {
		logger.Errorf("data source %s not found", query.DataSource)
		return errors.New("data source not found")
	}
	query.dataSource = &ds
	return nil
}

func LoadMany(logger *logrus.Logger, queries []map[string]interface{}, dataSources map[string]datasource.IDataSource) []*Query {
	var result []*Query
	for i, queryMap := range queries {
		result = append(result, &Query{})
		if err := result[i].Load(logger, queryMap, dataSources); err != nil {
			logger.Errorf("Unable to parse queryMapping from queries in config  = %s ", queryMap)
		}

	}
	return result
}
