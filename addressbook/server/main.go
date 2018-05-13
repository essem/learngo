package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	dbConnStr       = "dev:password@tcp(127.0.0.1:3306)/addressbook"
	grpcBindAddress = ":50051"
	webBindAddress  = ":8090"
)

func main() {
	db := &appDB{}
	numPeople, err := db.open(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Printf("There are %d people in database", numPeople)
	defer db.close()

	log.Println("Init grpc service...")
	grpcService := &grpcService{db: db}
	grpcService.init()

	log.Println("Init web service...")
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
