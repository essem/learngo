package main

import (
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql"

	"github.com/essem/learngo/addressbook/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcService struct {
	db         *appDB
	grpcServer *grpc.Server
}

func (s *grpcService) init() {
	grpcServer := grpc.NewServer()
	pb.RegisterAddressBookServiceServer(grpcServer, s)
	reflection.Register(grpcServer)

	s.grpcServer = grpcServer
}

func (s *grpcService) serve(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *grpcService) shutdown() {
	s.grpcServer.GracefulStop()
}

func (s *grpcService) List(ctx context.Context, in *pb.Empty) (*pb.ListReply, error) {
	log.Println("ListRequest", in)

	people, err := s.db.list()
	if err != nil {
		log.Fatal(err)
	}

	return &pb.ListReply{People: people}, nil
}

func (s *grpcService) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	log.Println("CreateRequest", in)

	id, err := s.db.create(in.Person)
	if err != nil {
		log.Fatal(err)
	}

	return &pb.CreateReply{Id: id}, nil
}

func (s *grpcService) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadReply, error) {
	log.Println("ReadRequest", in)

	person, err := s.db.read(in.Id)
	if err != nil {
		log.Println(err)
		return &pb.ReadReply{Person: nil}, nil
	}

	return &pb.ReadReply{Person: person}, nil
}

func (s *grpcService) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	log.Println("UpdateRequest", in)

	err := s.db.update(in.Person)
	if err != nil {
		log.Println(err)
		return &pb.UpdateReply{Success: false}, nil
	}

	return &pb.UpdateReply{Success: true}, nil
}

func (s *grpcService) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	log.Println("DeleteRequest", in)

	err := s.db.delete(in.Id)
	if err != nil {
		log.Println(err)
		return &pb.DeleteReply{Success: false}, nil
	}

	return &pb.DeleteReply{Success: true}, nil
}
