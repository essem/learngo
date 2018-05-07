package commands

import (
	"fmt"

	pb "github.com/essem/learngo/addressbookpb"
)

// List every people in address book
func List(book *pb.AddressBook) {
	for _, person := range book.People {
		fmt.Printf("ID: %d, Name: %s, Email: %s\n", person.Id, person.Name, person.Email)
	}
}
