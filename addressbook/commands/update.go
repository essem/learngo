package commands

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/essem/learngo/addressbook/pb"
)

// Update a people information
func Update(c pb.AddressBookServiceClient, reader *bufio.Reader) {
	fmt.Print("ID: ")
	idStr, _ := reader.ReadString('\n')

	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil {
		fmt.Printf("Could not convert id: %v\n", err)
		return
	}

	fmt.Print("Name: ")
	Name, _ := reader.ReadString('\n')

	fmt.Print("E-mail: ")
	Email, _ := reader.ReadString('\n')

	person := pb.Person{}
	person.Id = id
	person.Name = strings.TrimSpace(Name)
	person.Email = strings.TrimSpace(Email)

	r, err := c.Update(context.Background(), &pb.UpdateRequest{Person: &person})
	if err != nil {
		fmt.Printf("Could not update: %v\n", err)
		return
	}

	if !r.Success {
		fmt.Println("Failed to update")
	}
}
