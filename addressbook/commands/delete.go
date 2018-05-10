package commands

import (
	"bufio"
	"context"
	"fmt"
	"strconv"
	"strings"

	pb "github.com/essem/learngo/addressbookpb"
)

// Delete a person from address book
func Delete(c pb.AddressBookServiceClient, reader *bufio.Reader) {
	fmt.Print("ID: ")
	IDStr, _ := reader.ReadString('\n')

	ID, err := strconv.Atoi(strings.TrimSpace(IDStr))
	if err != nil {
		fmt.Printf("Could not convert id: %v\n", err)
		return
	}

	r, err := c.Delete(context.Background(), &pb.DeleteRequest{Id: int32(ID)})
	if err != nil {
		fmt.Printf("Could not delete: %v\n", err)
		return
	}

	if !r.Success {
		fmt.Println("Failed to delete")
	}
}
