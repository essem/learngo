package commands

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Delete a person from address book
func Delete(book *pb.AddressBook, reader *bufio.Reader) error {
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

	copy(book.People[found:], book.People[found+1:])
	book.People = book.People[:len(book.People)-1]

	return nil
}
