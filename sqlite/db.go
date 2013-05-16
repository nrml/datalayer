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
		err := errors.New("Db needs a namespace")
		return err
	}

	tst, err := sql.Open("sqlite3", db.Namespace)
	db.db = tst

	if db == nil {
		log.Fatal("database cannot be initialized")
	}

	return err
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) CreateTable(name string, tbl interface{}) (Table, error) {
	elem := reflect.TypeOf(tbl)

	//HACK: don't count last prop reflected on as it is tableObj
	length := elem.NumField() - 1

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

	return Table{name, fields, db}, err

}
