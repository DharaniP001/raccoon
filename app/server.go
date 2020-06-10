package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"raccoon/config"
	"raccoon/logger"
	"raccoon/publisher"
	ws "raccoon/websocket"
	"syscall"
	"time"
)

// StartServer starts the server
func StartServer(ctx context.Context, cancel context.CancelFunc) {

	//@TODO - create publisher with ctx

	//@TODO - create events-channels, workers (go routines) with ctx

	//start server @TODOD - the wss handler should be passed with the events-channel
	wssServer := ws.CreateServer()
	logger.Info("Start Server -->")
	wssServer.StartHTTPServer(ctx, cancel)

	//deliveryChan := make(chan kafka.Event)
	//
	kafkaConfig := config.NewKafkaConfig()
	//topic := kafkaConfig.Topic()

	//kafkaMessage := &kafka.Message{
	//	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//	Value:          []byte("Test"),
	//}

	kafkaProducer, err := publisher.NewKafkaProducer(kafkaConfig)

	if err != nil {
		logger.Error("Error creating kafka producer", err)
	}

	_ = publisher.NewProducer(kafkaProducer, config.NewKafkaConfig())

	//newProducer.Produce(kafkaMessage,deliveryChan)

	go shutDownServer(ctx, cancel)
}

func shutDownServer(ctx context.Context, cancel context.CancelFunc) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			logger.Info(fmt.Sprintf("[App.Server] Received a signal %s", sig))
			time.Sleep(3 * time.Second)
			logger.Info("Exiting server")
			os.Exit(0)
		default:
			logger.Info(fmt.Sprintf("[App.Server] Received a unexpected signal %s", sig))
		}
	}
}
