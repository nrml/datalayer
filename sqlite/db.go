package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
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
	} else {
		fmt.Println("database initialized")
	}

	return err
}

func (db *DB) Close() {
	db.db.Close()
}

func (db *DB) CreateTable(tbl Table) error {
	flds := make([]string, len(tbl.Fields))

	i := 0
	for _ = range tbl.Fields {
		flds[i] = fmt.Sprintf("%s %s", tbl.Fields[i].Name, tbl.Fields[i].Type)
		i++
	}

	stmt := fmt.Sprintf("create table %s (%s)", tbl.Name, strings.Join(flds, ","))

	fmt.Println(stmt)
	_, err := db.db.Exec(stmt)

	return err

}
