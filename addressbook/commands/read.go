package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Read a person information in address book
func Read(book *pb.AddressBook, reader *bufio.Reader) error {
	fmt.Print("ID: ")
	IDStr, _ := reader.ReadString('\n')

	ID, err := strconv.Atoi(strings.TrimSpace(IDStr))
	if err != nil {
		return err
	}

	found := -1
	for i, person := range book.People {
		if person.Id == int32(ID) {
			found = i
			break
		}
	}

	if found == -1 {
		return errors.New("Not found")
	}

	person := book.People[found]
	fmt.Printf("ID: %d\n", person.Id)
	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("E-mail: %s\n", person.Email)

	return nil
}
