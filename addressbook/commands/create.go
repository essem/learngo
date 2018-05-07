package commands

import (
	"bufio"
	"fmt"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Create new person and add to address book
func Create(book *pb.AddressBook, nextID int32, reader *bufio.Reader) error {
	fmt.Print("Name: ")
	Name, _ := reader.ReadString('\n')

	fmt.Print("E-mail: ")
	Email, _ := reader.ReadString('\n')

	person := &pb.Person{
		Id:    nextID,
		Name:  strings.TrimSpace(Name),
		Email: strings.TrimSpace(Email),
	}
	book.People = append(book.People, person)

	return nil
}
