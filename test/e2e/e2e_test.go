//go:build e2e
// +build e2e

package e2e_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/initiator"

	"github.com/hamzausmani302/prometheus-database-exporter/pkg/utils"

	"github.com/sirupsen/logrus"
)

var cacheStore = flag.String("cachestore", "local", "the store storage backend")
var schedulerStorage = flag.String("schedulerstore", "memory", "the storage backend to use for schduler")

var EXPECTED_MAP map[string]string = map[string]string{
	"product_inventory_product_count": "1",
	"taxi_rides_total_wells":          "10",
}

func TestEnd2EndApplicationKpisTest(t *testing.T) {
	// assumption : postgres and redis are already running
	rootLogger := logrus.New()
	done := make(chan bool, 1)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel() // Ensure resources are released

	// set the config file path
	utils.SetEnvironmentVariable("CONFIG_FILE_PATH", "config/config.test.yaml")
	go appStartUp(rootLogger, done)

	go func() {
		rootLogger.Info("triggered executing")
		<-done
		rootLogger.Info("Closing")
		close(done)
	}()
	time.Sleep(20 * time.Second)
	rootLogger.Info("Waiting for end")
	err := testExporter(rootLogger, ctx, done)
	done <- true
	if err != nil {
		t.Error(err)
	}
	<-done

	rootLogger.Info("Terminating program...")
}

func appStartUp(logger *logrus.Logger, done chan bool) {
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
	<-done
}

// waitForReady polls the /metrics endpoint until it returns HTTP 200 or the context expires.
func waitForReady(ctx context.Context, url string) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for %s to be ready", url)
		default:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				if body, err := io.ReadAll(resp.Body); err == nil && len(body) != 0 {
					resp.Body.Close()
					return nil
				}

			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// testExporter waits for the exporter to be ready, then validates metrics.
func testExporter(logger *logrus.Logger, ctx context.Context, done chan bool) error {
	logger.Info("Starting testing function")
	baseURL := "http://localhost:2112"
	if err := waitForReady(ctx, baseURL+"/metrics"); err != nil {
		return err
	}
	logger.Info("Exporter is ready")

	response, err := utils.SimpleGetRequest(baseURL + "/metrics")
	if err != nil {
		return err
	}
	metrics := fetchAllMetrics(response)
	for _, metric := range metrics {
		m, ok := EXPECTED_MAP[metric.MetricName]
		if ok && m != metric.Value {
			return errors.New(fmt.Sprintf("expected = %s, got = %s", m, metric.Value))
		}
	}
	return nil
}

func fetchAllMetrics(data string) []MetricResult {
	lines := strings.Split(data, "\n")
	metrics := []MetricResult{}
	for _, line := range lines {
		if !strings.HasPrefix(line, "#") {
			m := getMetricFromString(line)
			if m.MetricName != "" {
				metrics = append(metrics, m)
			}
		}
	}
	return metrics
}

type MetricResult struct {
	MetricName string
	Labels     string
	Value      string
}

func getMetricFromString(metricString string) MetricResult {
	var m MetricResult
	pattern := "^([a-zA-Z_:][a-zA-Z0-9_:]*)\\s*(\\{[^}]*\\})?\\s+([\\d.e+-]+)$"
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(metricString)
	if matches != nil {
		m = MetricResult{
			MetricName: matches[1],
			Labels:     matches[2],
			Value:      matches[3],
		}
	}
	return m
}
