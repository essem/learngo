package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/essem/learngo/addressbook/pb"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	dbFileName = "addressbook.db"
	port       = ":50051"
)

type server struct {
	nextID int32
	book   pb.AddressBook
}

func (s *server) List(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	log.Println("ListRequest", in)

	return &pb.ListReply{People: s.book.People}, nil
}

func (s *server) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	log.Println("CreateRequest", in)

	person := &pb.Person{
		Id:    s.nextID,
		Name:  strings.TrimSpace(in.Person.Name),
		Email: strings.TrimSpace(in.Person.Email),
	}
	s.book.People = append(s.book.People, person)

	s.nextID++

	return &pb.CreateReply{Id: person.Id}, nil
}

func (s *server) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadReply, error) {
	log.Println("ReadRequest", in)

	for _, person := range s.book.People {
		if person.Id == in.Id {
			return &pb.ReadReply{Person: person}, nil
		}
	}

	log.Printf("Not found ID: %d", in.Id)

	return &pb.ReadReply{Person: nil}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	log.Println("UpdateRequest", in)

	for _, person := range s.book.People {
		if person.Id == in.Person.Id {
			person.Name = strings.TrimSpace(in.Person.Name)
			person.Email = strings.TrimSpace(in.Person.Email)
			return &pb.UpdateReply{Success: true}, nil
		}
	}

	return &pb.UpdateReply{Success: false}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	log.Println("DeleteRequest", in)

	found := -1
	for i, person := range s.book.People {
		if person.Id == in.Id {
			found = i
			break
		}
	}

	if found == -1 {
		return &pb.DeleteReply{Success: false}, nil
	}

	copy(s.book.People[found:], s.book.People[found+1:])
	s.book.People = s.book.People[:len(s.book.People)-1]
	return &pb.DeleteReply{Success: true}, nil
}

func (s *server) load(dbFileName string) error {
	s.book = pb.AddressBook{}
	s.nextID = 1

	if _, err := os.Stat(dbFileName); err != nil {
		log.Fatalln("No database file exists.")
		return err
	}

	in, err := ioutil.ReadFile(dbFileName)
	if err != nil {
		log.Fatalln("Failed to read:", err)
		return err
	}

	if err := proto.Unmarshal(in, &s.book); err != nil {
		log.Fatalln("Failed to unmarshal:", err)
		return err
	}

	for i := 1; i < len(s.book.People); i++ {
		id := s.book.People[i].Id
		if id >= s.nextID {
			s.nextID = id + 1
		}
	}
	log.Printf("%d people loaded\n", len(s.book.People))

	return nil
}

func (s *server) write(dbFileName string) error {
	out, err := proto.Marshal(&s.book)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dbFileName, out, 0644); err != nil {
		return err
	}

	return nil
}

func main() {
	svr := &server{}
	svr.load(dbFileName)

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

	log.Println("Write changed to file")
	svr.write(dbFileName)
}
