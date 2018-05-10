package commands

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Update a people information
func Update(c pb.AddressBookServiceClient, reader *bufio.Reader) {
	fmt.Print("ID: ")
	IDStr, _ := reader.ReadString('\n')

	ID, err := strconv.Atoi(strings.TrimSpace(IDStr))
	if err != nil {
		fmt.Printf("Could not convert id: %v\n", err)
		return
	}

	fmt.Print("Name: ")
	Name, _ := reader.ReadString('\n')

	fmt.Print("E-mail: ")
	Email, _ := reader.ReadString('\n')

	person := pb.Person{}
	person.Id = int32(ID)
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
