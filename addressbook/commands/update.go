package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Update a people information
func Update(book *pb.AddressBook, reader *bufio.Reader) error {
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

	fmt.Print("Name: ")
	Name, _ := reader.ReadString('\n')

	fmt.Print("E-mail: ")
	Email, _ := reader.ReadString('\n')

	person := book.People[found]
	person.Name = strings.TrimSpace(Name)
	person.Email = strings.TrimSpace(Email)

	return nil
}
