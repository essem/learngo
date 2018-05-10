package commands

import (
	"context"
	"fmt"

	"github.com/essem/learngo/addressbook/pb"
)

// List every people in address book
func List(c pb.AddressBookServiceClient) {
	r, err := c.List(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Printf("Could not list: %v\n", err)
	}

	for _, person := range r.People {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", person.Id, person.Name, person.Email)
	}
}
