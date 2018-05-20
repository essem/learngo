package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("default")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("Failed to read default config", err)
	}

	viper.SetConfigName("local")
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalln("Failed to read local config", err)
		}
	}
}

func main() {
	initConfig()

	db := &appDB{}
	dbConnStr := viper.GetString("database.connStr")
	log.Println("Connect database...")
	numPeople, err := db.open(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Printf("There are %d people in database", numPeople)
	defer db.close()

	log.Println("Init grpc service...")
	grpcBindAddress := viper.GetString("grpc.bindAddress")
	grpcService := &grpcService{db: db}
	grpcService.init()

	log.Println("Init web service...")
	webBindAddress := viper.GetString("web.bindAddress")
	webService := &webService{db: db}
	webService.init(webBindAddress)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("Caught signal: %+v", sig)

		log.Println("Stop grpc service...")
		grpcService.shutdown()

		log.Println("Stop web service...")
		webService.shutdown()
	}()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		log.Printf("Start grpc service on %s", grpcBindAddress)
		if err := grpcService.serve(grpcBindAddress); err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Printf("Start web service on %s", webBindAddress)
		if err := webService.serve(); err != nil {
			log.Fatalln(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
