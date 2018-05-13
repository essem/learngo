package main

import (
	"reflect"
	"testing"

	"github.com/essem/learngo/addressbook/pb"
	_ "github.com/go-sql-driver/mysql"
)

const testDbConnStr = "dev:password@tcp(127.0.0.1:3306)/addressbook_test"

func Test_appDB(t *testing.T) {
	// Open and reset all data
	db := appDB{}
	_, err := db.open(testDbConnStr)
	if err != nil {
		t.Errorf("appDB.open() error = %v", err)
		return
	}
	db.reset()
	db.close()

	// Open again
	numPeople, err := db.open(testDbConnStr)
	if err != nil {
		t.Errorf("appDB.open() error = %v", err)
		return
	}
	defer db.close()

	if numPeople != 0 {
		t.Errorf("appDB.open() numPeople = %v, want %v", numPeople, 0)
		return
	}

	t.Run("Empty database", func(t *testing.T) {
		got, err := db.list()
		if err != nil {
			t.Errorf("appDB.list() error = %v", err)
			return
		}
		want := make([]*pb.Person, 0)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("appDB.list() = %v, want %v", got, want)
		}
	})

	t.Run("Create a person", func(t *testing.T) {
		person := pb.Person{
			Name:  "test",
			Email: "test@test.com",
		}
		got, err := db.create(&person)
		if err != nil {
			t.Errorf("appDB.list() error = %v", err)
			return
		}
		var want int64 = 1
		if !reflect.DeepEqual(got, want) {
			t.Errorf("appDB.list() = %v, want %v", got, want)
		}
	})

	t.Run("List database", func(t *testing.T) {
		got, err := db.list()
		if err != nil {
			t.Errorf("appDB.list() error = %v", err)
			return
		}
		want := []*pb.Person{
			{Id: 1, Name: "test", Email: "test@test.com"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("appDB.list() = %v, want %v", got, want)
		}
	})
}
