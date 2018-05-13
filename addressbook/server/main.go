package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"

	"github.com/essem/learngo/addressbook/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbConnStr = "dev:password@tcp(127.0.0.1:3306)/addressbook"
	port      = ":50051"
)

type server struct {
	db *appDB
}

func (s *server) List(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	log.Println("ListRequest", in)

	people, err := s.db.list()
	if err != nil {
		log.Fatal(err)
	}

	return &pb.ListReply{People: people}, nil
}

func (s *server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	log.Println("CreateRequest", in)

	id, err := s.db.create(in.Person)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.CreateReply{Id: id}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadReply, error) {
	log.Println("ReadRequest", in)

	person, err := s.db.read(in.Id)
	if err != nil {
		log.Println(err)
		return &pb.ReadReply{Person: nil}, nil
	}

	return &pb.ReadReply{Person: person}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	log.Println("UpdateRequest", in)

	err := s.db.update(in.Person)
	if err != nil {
		log.Println(err)
		return &pb.UpdateReply{Success: false}, nil
	}

	return &pb.UpdateReply{Success: true}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	log.Println("DeleteRequest", in)

	err := s.db.delete(in.Id)
	if err != nil {
		log.Println(err)
		return &pb.DeleteReply{Success: false}, nil
	}

	return &pb.DeleteReply{Success: true}, nil
}

func main() {
	db := &appDB{}
	numPeople, err := db.open(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	log.Printf("There are %d people in database", numPeople)
	defer db.close()

	svr := &server{db: db}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterAddressBookServiceServer(s, svr)
	reflection.Register(s)

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		log.Printf("Caught signal: %+v", sig)
		log.Println("Try to stop gracefully...")
		s.GracefulStop()
	}()

	log.Printf("Start service on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
