package commands

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/essem/learngo/addressbook/pb"
)

// Delete a person from address book
func Delete(c pb.AddressBookServiceClient, reader *bufio.Reader) {
	fmt.Print("ID: ")
	idStr, _ := reader.ReadString('\n')

	id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
	if err != nil {
		fmt.Printf("Could not convert id: %v\n", err)
		return
	}

	r, err := c.Delete(context.Background(), &pb.DeleteRequest{Id: id})
	if err != nil {
		fmt.Printf("Could not delete: %v\n", err)
		return
	}

	if !r.Success {
		fmt.Println("Failed to delete")
	}
}
