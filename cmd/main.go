package main

import (
	"github.com/bhupeshpandey/task-manager-ashland/internal/config"
	"github.com/bhupeshpandey/task-manager-ashland/internal/logger"
	"github.com/bhupeshpandey/task-manager-ashland/internal/models"
	"github.com/bhupeshpandey/task-manager-ashland/internal/msgqueue"
	. "github.com/bhupeshpandey/task-manager-ashland/internal/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
)

func main() {
	// Load configuration
	conf, err := config.LoadConfig("./config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	taskLogger := logger.NewLogger(conf.Logging)

	// Create a new Prometheus registry
	reg := prometheus.NewRegistry()

	conf.MessageQueue.Registry = reg

	msgQueue := msgqueue.NewMessageQueue(conf.MessageQueue, taskLogger)

	go StartServer(conf.PrometheusConfig, taskLogger)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := msgQueue.ReceiveMessages(&wg)
		if err != nil {
			taskLogger.Log(models.ErrorLevel, "Unable to receive messages from message queue", err.Error())
		}
	}()
	wg.Wait()
}
