package sqlite

import (
	"fmt"
	"strings"
)

type Table struct {
	Name   string
	Fields []Field
	DB     DB
}

func (t *Table) fieldNames() []string {
	keys := make([]string, len(t.Fields))
	i := 0
	for _, _ = range t.Fields {
		keys[i] = t.Fields[i].Name
		i++
	}
	return keys
}

func (tbl *Table) Get(id int64, values ...interface{}) error {
	statement := "select " + strings.Join(tbl.fieldNames(), ",") + " from " + tbl.Name + " where id = ?"

	stmt, err := tbl.DB.db.Prepare(statement)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stmt.QueryRow(&id).Scan(values...)

	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func (tbl *Table) Create(values []string) (int64, error) {
	cnt := len(tbl.Fields)

	keys := make([]string, cnt)
	vals := make([]string, cnt)

	i := 0
	for _, _ = range tbl.Fields {
		f := tbl.Fields[i]
		keys[i] = f.Name
		if f.Type == "text" {
			vals[i] = strings.Join([]string{"'", values[i], "'"}, "")
		} else {
			vals[i] = values[i]
		}

		i++
	}

	statement := "insert into " + tbl.Name + "(" + strings.Join(keys, ",") + ") values(" + strings.Join(vals, ",") + ")"

	tx, err := tbl.DB.db.Begin()

	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	stmt, err := tx.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec()

	id, _ := res.LastInsertId()

	if err != nil {
		return -1, err
	}

	tx.Commit()

	return id, err

}
