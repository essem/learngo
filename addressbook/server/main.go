package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	dbConnStr       = "dev:password@tcp(127.0.0.1:3306)/addressbook"
	grpcBindAddress = ":50051"
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

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("Caught signal: %+v", sig)

		log.Println("Stop grpc service...")
		grpcService.shutdown()
	}()

	log.Printf("Start grpc service on %s", grpcBindAddress)
	if err := grpcService.serve(grpcBindAddress); err != nil {
		log.Fatalln(err)
	}
}
