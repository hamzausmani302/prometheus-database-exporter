package queryscheduler

import (
	"time"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/datasource"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/schema"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/cache"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/go-scheduler"
	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"

	"github.com/sirupsen/logrus"
)

/* Interface to be implemneted for the task scheduler, will be helpful for mocking query scheduler
Underneath we are using the scheduler libraray but can be implemneted for custom implementation

Only scheduler can also be mocked by assigning a different implementation with the same interface the go-scheduler package follows
interface {
	Start() error
	Stop()
	RunEvery(duration time.Duration, func task.Function, ...taks.Params) task.ID, error
	Note: Although I dont think that is a good idea
}
*/
type IQueryScheduler interface{
	// Initialization of objects & scheduler
	Init() error
	// Start the scheduler
	Start() error
	// Define what to do with each query object
	ExecuteTask(query *schema.Query) error
	// Stop the schduler
	Stop() error
}

type QueryScheduler struct{
	Queries []*schema.Query
	cfg *config.ApplicationConfig
	logger *logrus.Logger
	scheduler *scheduler.Scheduler
	programChannel *chan bool
	cacheStore *cache.ICache
}

func (q *QueryScheduler) Init() error{
	q.logger.Infof("total number of Queries : %d", len(q.Queries))
	for _, query := range q.Queries {
		// assiging the schduled task id hash
		var taskId string
		if id, err := q.scheduler.RunEvery(time.Duration(query.QueryRefreshTime) * time.Second, q.ExecuteTask, query  ); err != nil {
			q.logger.Errorf("Error while running task with id = %s", id )
			q.logger.Debugf("Error while running task with id = %s | query = %s | %d", query, query.QueryRefreshTime)
			return err
		}else{
			taskId = string(id)
		}
		query.SetHash(taskId)
	}
	return nil
}
func (q *QueryScheduler) Start() error{
	if err := q.scheduler.Start(); err != nil {
		q.logger.Fatal("Error running scheduler", err)
	}
	return nil
}

func (q *QueryScheduler) Stop() error {
	q.scheduler.Stop()
	return nil
}
// The actual task workflow will be written here
func (q *QueryScheduler) ExecuteTask(query *schema.Query) error {
	now := time.DateTime
	q.logger.Infof("Executing task for %s %s %s | %d", query.Name, query.Query, now, query.QueryRefreshTime)
	ds := *query.GetDataSource()
	// get data from database
	if err := ds.Connect(); err != nil {
		q.logger.Errorf("Error connecting to data source %s", query.DataSource)
		return err
	}
	df := ds.GetData(datasource.SQLQuery{
		Query: query.Query,
	})
	q.logger.Debug(df)
	
	// put the data with the key in cacheStore 
	if bytesDf,err := utils.DataFrameToCSVBytes(df); err == nil{
		(*q.cacheStore).Set(query.Query, bytesDf, int64(query.QueryRefreshTime * 2))
	} else{
		q.logger.Errorf("error converting dataframe to bytes ", df)
	}
	return nil
}

func NewQuerySchduler(logger *logrus.Logger, cfg *config.ApplicationConfig, baseScheduler *scheduler.Scheduler, queries []*schema.Query, store *cache.ICache, channel *chan bool) *QueryScheduler {
	return &QueryScheduler{logger: logger, cfg: cfg, scheduler: baseScheduler, Queries: queries, programChannel: channel, cacheStore: store}
}