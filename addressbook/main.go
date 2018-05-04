package main

import (
	"fmt"
	"io/ioutil"
	"log"

	pb "github.com/essem/learngo/addressbookpb"
	"github.com/golang/protobuf/proto"
)

func write() error {
	p := pb.Person{
		Id:    1234,
		Name:  "John Doe",
		Email: "jdoe@example.com",
		Phones: []*pb.Person_PhoneNumber{
			{Number: "555-4321", Type: pb.Person_HOME},
		},
	}

	fmt.Println("Write", p)

	out, err := proto.Marshal(&p)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile("person.db", out, 0644); err != nil {
		return err
	}

	return nil
}

func read() error {
	in, err := ioutil.ReadFile("person.db")
	if err != nil {
		return err
	}

	var p pb.Person
	if err := proto.Unmarshal(in, &p); err != nil {
		return err
	}

	fmt.Println("Read", p)
	return nil
}

func main() {
	err := write()
	if err != nil {
		log.Fatalln("Failed to write:", err)
	}

	err = read()
	if err != nil {
		log.Fatalln("Failed to read:", err)
	}
}
