//go:build e2e
// +build e2e

package e2e_test

import (
	"fmt"
	"testing"
	"time"
	"context"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/initiator"

	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"

	"github.com/sirupsen/logrus"
)


func TestEnd2EndApplicationKpisTest(t *testing.T){
	// assumption : postgres and redis are already running
	// load config from the new file
	rootLogger := logrus.New()
	done := make(chan bool, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // Ensure resources are released

	// set the config file path
	utils.SetEnvironmentVariable("CONFIG_FILE_PATH", "config/config.test.yaml")
	cfg := config.GetConfig("test", rootLogger)
	fmt.Println(cfg)

	

	go appStartUp(rootLogger, done)
	
	go func() {
		// Listens for intended termination and terminate the memory addresses
		rootLogger.Info("triggered executing")
		// sig := <-sigs
		<- done
		// rootLogger.Debug(sig)
		// close scheduler
		rootLogger.Info("Closing")
		// close(sigs)
		close(done)
	}()
	rootLogger.Info("Waiting for end")
	testExporter(&ctx, &cfg, done)
	
	<- done
	
	rootLogger.Info("Terminating program...")

}

func appStartUp(logger *logrus.Logger, done chan bool){

	app := initiator.Application{Done: done}
	if err := app.Init(); err != nil {
		logger.Panic("Failed to initialize application", err)
		return
	}
	if app.IsCollectorEnabled() {
		go app.StartCollector()
	}
	if app.IsApiEnabled() {
		go app.StartApi()
	}
	<- done
	fmt.Println("Clean up Scheduler")
	
}


func testExporter(ctx *context.Context,cfg *config.ApplicationConfig ,done chan bool) error{
	time.Sleep(15 * time.Second)
	baseUrl := "http://localhost:2112"
	response, err := utils.SimpleGetRequest(baseUrl + "/app-metrics")
	if err != nil {
		return err
	}
	fmt.Println(response)
	done <- true
	return nil
}


func fetchAllMetrics(data string) {
	
}