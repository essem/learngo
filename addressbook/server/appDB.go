package main

import (
	"database/sql"
	"errors"

	"github.com/essem/learngo/addressbook/pb"
	_ "github.com/go-sql-driver/mysql"
)

type appDB struct {
	sqlDB *sql.DB
}

func (db *appDB) open(connStr string) (int, error) {
	sqlDB, err := sql.Open("mysql", connStr)
	if err != nil {
		return 0, err
	}

	var numPeople int
	err = sqlDB.QueryRow("SELECT COUNT(id) FROM people").Scan(&numPeople)
	if err != nil {
		sqlDB.Close()
		return 0, err
	}

	db.sqlDB = sqlDB
	return numPeople, nil
}

func (db *appDB) close() {
	db.sqlDB.Close()
}

func (db *appDB) list() ([]*pb.Person, error) {
	rows, err := db.sqlDB.Query("SELECT id, name, email FROM people")
	if err != nil {
		return nil, err
	}

	people := make([]*pb.Person, 0)
	for rows.Next() {
		var id int64
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			return nil, err
		}
		people = append(people, &pb.Person{Id: id, Name: name, Email: email})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return people, nil
}

func (db *appDB) create(person *pb.Person) (int64, error) {
	r, err := db.sqlDB.Exec("INSERT INTO people (name, email) VALUES (?, ?)", person.Name, person.Email)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (db *appDB) read(id int64) (*pb.Person, error) {
	rows := db.sqlDB.QueryRow("SELECT name, email FROM people WHERE id = ?", id)
	var name, email string
	if err := rows.Scan(&name, &email); err != nil {
		return nil, err
	}

	return &pb.Person{Id: id, Name: name, Email: email}, nil
}

func (db *appDB) update(person *pb.Person) error {
	r, err := db.sqlDB.Exec("UPDATE people SET name = ?, email = ? WHERE id = ?",
		person.Name, person.Email, person.Id)
	if err != nil {
		return err
	}

	numAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if numAffected != 1 {
		return errors.New("Number of affected rows is not one")
	}

	return nil
}

func (db *appDB) delete(id int64) error {
	r, err := db.sqlDB.Exec("DELETE FROM people WHERE id = ?", id)
	if err != nil {
		return err
	}

	numAffected, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if numAffected != 1 {
		return errors.New("Number of affected rows is not one")
	}

	return nil
}
