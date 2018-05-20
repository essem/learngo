package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

func main() {
	initConfig()
	initLog()

	db := &appDB{}
	dbConnStr := viper.GetString("database.connStr")
	log.WithField("dbConnStr", dbConnStr).Info("Connect database...")
	numPeople, err := db.open(dbConnStr)
	if err != nil {
		log.WithError(err).Fatal("Failed to open database")
	}
	log.WithField("num", numPeople).Info("Loaded people from database")
	defer db.close()

	log.Info("Init grpc service...")
	grpcBindAddress := viper.GetString("grpc.bindAddress")
	grpcService := &grpcService{db: db}
	grpcService.init()

	log.Info("Init web service...")
	webBindAddress := viper.GetString("web.bindAddress")
	webService := &webService{db: db}
	webService.init(webBindAddress)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.WithField("signal", sig).Info("Caught signal")

		log.Info("Stop grpc service...")
		grpcService.shutdown()

		log.Info("Stop web service...")
		webService.shutdown()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.WithField("address", grpcBindAddress).Info("Start grpc service")
		if err := grpcService.serve(grpcBindAddress); err != nil {
			log.WithError(err).Fatal("Failed to grpc service")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.WithField("address", webBindAddress).Info("Start web service")
		if err := webService.serve(); err != nil {
			log.WithError(err).Fatal("Failed to web service")
		}
		wg.Done()
	}()

	wg.Wait()
}
