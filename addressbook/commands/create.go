package commands

import (
	"bufio"
	"context"
	"fmt"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Create new person and add to address book
func Create(c pb.AddressBookServiceClient, reader *bufio.Reader) {
	fmt.Print("Name: ")
	Name, _ := reader.ReadString('\n')

	fmt.Print("E-mail: ")
	Email, _ := reader.ReadString('\n')

	person := &pb.Person{
		Id:    0,
		Name:  strings.TrimSpace(Name),
		Email: strings.TrimSpace(Email),
	}

	r, err := c.Create(context.Background(), &pb.CreateRequest{Person: person})
	if err != nil {
		fmt.Printf("Could not create: %v\n", err)
		return
	}

	fmt.Printf("New ID: %d\n", r.Id)
}
