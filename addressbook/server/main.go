package main

import (
	"database/sql"
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
	db *sql.DB
}

func (s *server) List(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	log.Println("ListRequest", in)

	rows, err := s.db.Query("SELECT id, name, email FROM people")
	if err != nil {
		log.Fatal(err)
	}

	people := make([]*pb.Person, 0)
	for rows.Next() {
		var id int32
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			log.Fatal(err)
		}
		people = append(people, &pb.Person{Id: id, Name: name, Email: email})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return &pb.ListReply{People: people}, nil
}

func (s *server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	log.Println("CreateRequest", in)

	r, err := s.db.Exec("INSERT INTO people (name, email) VALUES (?, ?)", in.Person.Name, in.Person.Email)
	if err != nil {
		log.Fatal(err)
	}

	id, err := r.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	return &pb.CreateReply{Id: int32(id)}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadReply, error) {
	log.Println("ReadRequest", in)

	rows := s.db.QueryRow("SELECT id, name, email FROM people WHERE id = ?", in.Id)
	var id int32
	var name, email string
	if err := rows.Scan(&id, &name, &email); err != nil {
		log.Printf("Query failed: %v", err)
		return &pb.ReadReply{Person: nil}, nil
	}

	person := &pb.Person{Id: id, Name: name, Email: email}
	return &pb.ReadReply{Person: person}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	log.Println("UpdateRequest", in)

	r, err := s.db.Exec("UPDATE people SET name = ?, email = ? WHERE id = ?",
		in.Person.Name, in.Person.Email, in.Person.Id)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return &pb.UpdateReply{Success: false}, nil
	}

	numAffected, err := r.RowsAffected()
	if err != nil || numAffected != 1 {
		return &pb.UpdateReply{Success: false}, nil
	}

	return &pb.UpdateReply{Success: true}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	log.Println("DeleteRequest", in)

	r, err := s.db.Exec("DELETE FROM people WHERE id = ?", in.Id)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return &pb.DeleteReply{Success: false}, nil
	}

	numAffected, err := r.RowsAffected()
	if err != nil || numAffected != 1 {
		return &pb.DeleteReply{Success: false}, nil
	}

	return &pb.DeleteReply{Success: true}, nil
}

func (s *server) init() {
	log.Println("Init server")

	db, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	s.db = db

	var numPeople int
	err = db.QueryRow("SELECT COUNT(id) FROM people").Scan(&numPeople)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("There are %d people in database", numPeople)
}

func (s *server) cleanup() {
	log.Println("Cleanup server")

	s.db.Close()
}

func main() {
	svr := &server{}
	svr.init()

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

	svr.cleanup()
}
