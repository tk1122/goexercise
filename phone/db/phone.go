package db

import (
	"database/sql"
	"fmt"
)

// Phone represent schema of phone_number table
type Phone struct {
	Id     int
	Number string
}

// DB is a convinient wrapper of *sql.DB
type DB struct {
	db *sql.DB
}

func (db *DB) AllPhones() ([]Phone, error) {
	rows, err := db.db.Query("SELECT id, value FROM phone_number")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.Id, &p.Number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (db *DB) FindPhoneExcept(number string, id int) (*Phone, error) {
	var p Phone
	err := db.db.QueryRow("SELECT * FROM phone_number WHERE value=$1 AND id!=$2", number, id).Scan(&p.Id, &p.Number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil

}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	stm := `INSERT INTO phone_number(value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(stm, phone).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("cannot insert phone %s: %s", phone, err)
	}
	return id, nil
}

func (db *DB) Seed() error {
	data := []string{
		"1234567890",
		"123 456 7891",
		"(123) 456 7892",
		"(123) 456-7893",
		"123-456-7894",
		"123-456-7890",
		"1234567892",
		"(123)456-7892",
	}
	for _, number := range data {
		if _, err := insertPhone(db.db, number); err != nil {
			return err
		}
	}
	return nil
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetDB(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = createPhoneNumberTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func createPhoneNumberTable(db *sql.DB) error {
	stm := `
		CREATE TABLE IF NOT EXISTS phone_number (
			id SERIAL,
			value VARCHAR(255))
		`
	_, err := db.Exec(stm)
	if err != nil {
		return fmt.Errorf("cannot create table phone_number: %s", err)
	}
	return err
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return fmt.Errorf("cannot create database %s: %s", name, err)
	}
	return nil
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return fmt.Errorf("cannot drop database %s: %s", name, err)
	}
	return createDB(db, name)
}

func (db *DB) getPhone(id int) (string, error) {
	var number string
	err := db.db.QueryRow("SELECT * FROM phone_number WHERE id=$1", id).Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func (db *DB) UpdatePhone(p Phone) error {
	stm := `UPDATE phone_number SET value=$2 WHERE id=$1`
	_, err := db.db.Exec(stm, p.Id, p.Number)
	return err
}

func (db *DB) DeletePhone(id int) error {
	stm := `DELETE FROM phone_number WHERE id=$1`
	_, err := db.db.Exec(stm, id)
	return err
}
