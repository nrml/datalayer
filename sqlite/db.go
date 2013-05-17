package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"reflect"
	"strings"
)

type DB struct {
	Namespace string
	db        *sql.DB
}

func (db *DB) Init() error {
	if db.Namespace == "" {
		err := errors.New("db needs a namespace")
		return err
	}
	tst, err := sql.Open("sqlite3", db.Namespace)
	db.db = tst

	if db.db == nil {
		log.Fatal("database cannot be initialized")
	}

	return err
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) Table(name string, tblType interface{}) Table {
	elem := reflect.TypeOf(tblType)
	length := elem.NumField()
	fields := make([]Field, length)
	for i := 0; i < length; i++ {
		f := elem.Field(i)
		tag := f.Tag.Get("sql")
		if tag != "" {
			parts := strings.Split(tag, " ")
			fields[i] = Field{parts[0], strings.Join(parts[1:], " ")}
		}
	}

	return Table{name, reflect.TypeOf(tblType), fields, db}
}

func (db *DB) CreateTable(name string, tblType interface{}) (Table, error) {
	elem := reflect.TypeOf(tblType)

	length := elem.NumField()

	defs := make([]string, length)
	fields := make([]Field, length)

	for i := 0; i < length; i++ {
		f := elem.Field(i)
		tag := f.Tag.Get("sql")
		if tag != "" {
			defs[i] = tag
			parts := strings.Split(tag, " ")
			fields[i] = Field{parts[0], strings.Join(parts[1:], " ")}
		}
	}

	stmt := fmt.Sprintf("create table %s (%s)", name, strings.Join(defs, ","))

	_, err := db.db.Exec(stmt)

	if err != nil {
		fmt.Printf("ERROR creating table: %v   - %v\n", err.Error(), "setting error to nil")
		err = nil
	}

	return Table{name, reflect.TypeOf(tblType), fields, db}, err

}
