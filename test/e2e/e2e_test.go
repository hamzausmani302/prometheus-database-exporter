//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hamzausmani302/prometheus-database-exporter/config"
	"github.com/hamzausmani302/prometheus-database-exporter/internal/initiator"

	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"

	"github.com/sirupsen/logrus"
)

var EXPECTED_MAP map[string]string = map[string]string{
	"taxi_rides123_total_wells1": "10",
	"taxi_rides12_total_wells1": "10",
	"taxi_rides_total_wells": "10",
}

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
	err := testExporter(&ctx, &cfg, done)
	
	done <- true
	if err != nil {
		t.Error(err)
	}
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
	
}

// Test the exporter running as a service
func testExporter(ctx *context.Context,cfg *config.ApplicationConfig ,done chan bool) error{
	time.Sleep(30 * time.Second)		// Wait for 30 sec then test 
	baseUrl := "http://localhost:2112"
	response, err := utils.SimpleGetRequest(baseUrl + "/metrics")
	if err != nil {
		return err
	}
	metrics := fetchAllMetrics(response)
	for _, metric := range metrics{
		m , ok:= EXPECTED_MAP[metric.MetricName]
		if ok && m != metric.Value  {
			return errors.New(fmt.Sprintf("expected = %s, got = %s", m , metric.Value))
		} 
	}
	fmt.Println(metrics)
	return nil
}


func fetchAllMetrics(data string) []MetricResult{
	lines := strings.Split(data, "\n")
	metrics := []MetricResult{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			m := getMetricFromString(line)
			if m.MetricName != "" {
				metrics = append(metrics,  m)
			}
		} 
	} 
	return metrics
}

type MetricResult struct{
	MetricName string
	Labels string
	Value string
}

func getMetricFromString(metricString string) MetricResult {
	var m MetricResult;
	// fmt.Println("metricString", metricString)
	pattern := "^([a-zA-Z_:][a-zA-Z0-9_:]*)\\s*(\\{[^}]*\\})?\\s+([\\d.e+-]+)$"
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(metricString)
	if matches != nil {
		m = MetricResult{
			MetricName: matches[1],
			Labels: matches[2],
			Value: matches[3],
		}
	}
	return m

}