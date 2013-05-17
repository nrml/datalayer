package sqlite

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Table struct {
	Name   string
	Fields []Field
	DB     *DB
}

// type Stringer interface {
// 	String() string
// }

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

	return err
}

func (tbl *Table) Fill(id int64, obj interface{}) error {
	statement := "select " + strings.Join(tbl.fieldNames(), ",") + " from " + tbl.Name + " where id = ?"

	stmt, err := tbl.DB.db.Prepare(statement)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	//err = stmt.QueryRow(&id).Scan(values...)
	rows, err := stmt.Query(&id)
	//obj.Fill(row, obj)
	tbl.fill(rows, obj)
	return err
}

func (tbl *Table) Create(obj interface{}) (int64, error) {
	elem := reflect.ValueOf(obj)
	length := elem.NumField() - 2

	fmt.Printf("fields count: %v\n", length)

	keys := make([]string, length)
	vals := make([]string, length)

	for i := 1; i < length+1; i++ {
		f := tbl.Fields[i]
		keys[i-1] = f.Name
		if f.Type == "text" {
			vals[i-1] = strings.Join([]string{"'", elem.Field(i).String(), "'"}, "")
		} else {
			it, ok := elem.Field(i).Interface().(int64)
			if ok {
				vals[i-1] = strconv.FormatInt(it, 10)
			}
		}
	}

	statement := "insert into " + tbl.Name + "(" + strings.Join(keys, ",") + ") values(" + strings.Join(vals, ",") + ")"
	fmt.Println(statement)
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
func (tbl *Table) Update(id int64, values ...interface{}) error {

	//TODO concurrency

	matches := make([]string, len(values))
	i := 0

	for _, _ = range values {
		//skip id field
		f := tbl.Fields[i+1]
		sval := "?"

		if f.Type == "string" {
			sval = "'?'"
		}
		matches[i] = fmt.Sprintf("%s=%s", f.Name, sval)
		i++
	}

	statement := fmt.Sprintf("update %s set %s where id=%s", tbl.Name, strings.Join(matches, ","), strconv.FormatInt(id, 10))

	fmt.Println(statement)

	tx, err := tbl.DB.db.Begin()

	stmt, err := tx.Prepare(statement)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)

	tx.Commit()

	return err

}

func (tbl *Table) Delete(id int64) error {
	tx, err := tbl.DB.db.Begin()
	if err != nil {
		fmt.Println(err)
		return err
	}
	stmt, err := tx.Prepare("delete from " + tbl.Name + " where id=?")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	tx.Commit()

	return err
}
