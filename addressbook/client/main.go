package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/essem/learngo/addressbook/commands"
	pb "github.com/essem/learngo/addressbookpb"
	"google.golang.org/grpc"
)

const address = "localhost:50051"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Can not connect: %v\n", err)
		return
	}
	defer conn.Close()
	c := pb.NewAddressBookServiceClient(conn)

CommandLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command(l/c/r/u/d/q): ")
		command, _ := reader.ReadString('\n')
		switch strings.TrimSpace(command) {
		case "l":
			commands.List(c)
		case "c":
			commands.Create(c, reader)
		case "r":
			commands.Read(c, reader)
		case "u":
			commands.Update(c, reader)
		case "d":
			commands.Delete(c, reader)
		case "q":
			break CommandLoop
		default:
			fmt.Printf("Invalid command: %s", command)
		}
	}
}
