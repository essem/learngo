package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/essem/learngo/addressbook/commands"
	pb "github.com/essem/learngo/addressbookpb"
	"github.com/golang/protobuf/proto"
)

const dbFileName = "addressbook.db"

func write(book *pb.AddressBook) error {
	out, err := proto.Marshal(book)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(dbFileName, out, 0644); err != nil {
		return err
	}

	return nil
}

func load(book *pb.AddressBook) error {
	in, err := ioutil.ReadFile(dbFileName)
	if err != nil {
		return err
	}

	if err := proto.Unmarshal(in, book); err != nil {
		return err
	}

	return nil
}

func main() {
	book := pb.AddressBook{}
	var nextID int32 = 1
	if _, err := os.Stat("addressbook.db"); err == nil {
		if err := load(&book); err != nil {
			log.Fatalln("Failed to load:", err)
			return
		}
		for i := 1; i < len(book.People); i++ {
			id := book.People[i].Id
			if id >= nextID {
				nextID = id + 1
			}
		}
		fmt.Printf("%d people loaded\n", len(book.People))
	} else {
		fmt.Println("No database file exists.")
	}

CommandLoop:
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command(l/c/r/u/d/q): ")
		command, _ := reader.ReadString('\n')
		switch strings.TrimSpace(command) {
		case "l":
			commands.List(&book)
		case "c":
			commands.Create(&book, nextID, reader)
			nextID++
		case "r":
			if err := commands.Read(&book, reader); err != nil {
				fmt.Printf("Failed to read: %s\n", err)
			}
		case "u":
			if err := commands.Update(&book, reader); err != nil {
				fmt.Printf("Failed to update: %s\n", err)
			}
		case "d":
			if err := commands.Delete(&book, reader); err != nil {
				fmt.Printf("Failed to delete: %s\n", err)
			}
		case "q":
			break CommandLoop
		default:
			fmt.Printf("Invalid command: %s", command)
		}
	}

	err := write(&book)
	if err != nil {
		log.Fatalln("Failed to write:", err)
	}
}
