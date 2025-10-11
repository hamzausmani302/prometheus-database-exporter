package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hamzausmani302/prometheus-database-exporter/internal/initiator"
	"github.com/sirupsen/logrus"
)

func main() {
	rootLogger := logrus.New()
	sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    done := make(chan bool, 1)

	app := initiator.Application{Done: done}
	app.Init()
	go app.StartCollector()
	go app.StartApi()

	go func() {
		// Listens for intended termination and terminate the memory addresses
		rootLogger.Info("triggered executing")
        sig := <-sigs
        rootLogger.Debug(sig)
        done <- true
		// close scheduler
		if err := app.CleanUp(); err != nil {
			rootLogger.Error(err)
		}else{
			rootLogger.Info("closed successfully")
		}
		close(sigs)
		close(done)
	}()
	
	<- done
	rootLogger.Info("Exiting out of main thread")
}

