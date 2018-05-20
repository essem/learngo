package main

import (
	"net"

	log "github.com/sirupsen/logrus"

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
	log.WithField("pb", in).Info("ListRequest")

	people, err := s.db.list()
	if err != nil {
		log.WithError(err).Fatal("Failed to list")
	}

	return &pb.ListReply{People: people}, nil
}

func (s *grpcService) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateReply, error) {
	log.WithField("pb", in).Info("CreateRequest")

	id, err := s.db.create(in.Person)
	if err != nil {
		log.WithError(err).Fatal("Failed to create")
	}

	return &pb.CreateReply{Id: id}, nil
}

func (s *grpcService) Read(ctx context.Context, in *pb.ReadRequest) (*pb.ReadReply, error) {
	log.WithField("pb", in).Info("ReadRequest")

	person, err := s.db.read(in.Id)
	if err != nil {
		log.WithError(err).Debug("Failed to read")
		return &pb.ReadReply{Person: nil}, nil
	}

	return &pb.ReadReply{Person: person}, nil
}

func (s *grpcService) Update(ctx context.Context, in *pb.UpdateRequest) (*pb.UpdateReply, error) {
	log.WithField("pb", in).Info("UpdateRequest")

	err := s.db.update(in.Person)
	if err != nil {
		log.WithError(err).Debug("Failed to update")
		return &pb.UpdateReply{Success: false}, nil
	}

	return &pb.UpdateReply{Success: true}, nil
}

func (s *grpcService) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteReply, error) {
	log.WithField("pb", in).Info("DeleteRequest")

	err := s.db.delete(in.Id)
	if err != nil {
		log.WithError(err).Debug("Failed to delete")
		return &pb.DeleteReply{Success: false}, nil
	}

	return &pb.DeleteReply{Success: true}, nil
}
