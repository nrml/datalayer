package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Table struct {
	Name   string
	Type   reflect.Type
	Fields []Field
	DB     *DB
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
func (tbl *Table) Search(search string) ([]interface{}, error) {
	statement := "select " + strings.Join(tbl.fieldNames(), ",") + " from " + tbl.Name + " where " + search

	stmt, err := tbl.DB.db.Prepare(statement)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	filled := tbl.fill(rows)

	return filled, err
}
func (tbl *Table) Get(id int64) (interface{}, error) {
	statement := "select " + strings.Join(tbl.fieldNames(), ",") + " from " + tbl.Name + " where id = ?"

	stmt, err := tbl.DB.db.Prepare(statement)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(&id)

	filled := tbl.fill(rows)

	if len(filled) > 0 {
		obj := filled[0]
		return obj, err
	}

	return nil, err
}

func (tbl *Table) List() ([]interface{}, error) {
	statement := "select " + strings.Join(tbl.fieldNames(), ",") + " from " + tbl.Name

	stmt, err := tbl.DB.db.Prepare(statement)

	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()

	filled := tbl.fill(rows)

	return filled, err
}

func (tbl *Table) Create(obj interface{}) (int64, error) {
	elem := reflect.ValueOf(obj)
	length := elem.NumField() - 1

	keys := make([]string, length)
	vals := make([]string, length)
	values := make([]interface{}, length)

	for i := 1; i < length+1; i++ {
		f := tbl.Fields[i]
		keys[i-1] = f.Name
		vals[i-1] = "?"
		face := elem.Field(i).Interface()

		values[i-1] = face
	}

	statement := "insert into " + tbl.Name + "(" + strings.Join(keys, ",") + ") values(" + strings.Join(vals, ",") + ")"

	tx, err := tbl.DB.db.Begin()

	if err != nil {
		return -1, err
	}
	stmt, err := tx.Prepare(statement)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(values...)

	id, _ := res.LastInsertId()

	tx.Commit()

	return id, err

}
func (tbl *Table) Update(id int64, obj interface{}) error {

	elem := reflect.ValueOf(obj)
	length := len(tbl.Fields)
	// //TODO concurrency
	matches := make([]string, length-1)
	values := make([]interface{}, length-1)

	//skip id field
	for i := 1; i < length; i++ {
		f := tbl.Fields[i]
		values[i-1] = elem.Field(i).Interface()
		matches[i-1] = fmt.Sprintf("%s=?", f.Name)
	}

	statement := fmt.Sprintf("update %s set %s where id=%s", tbl.Name, strings.Join(matches, ","), strconv.FormatInt(id, 10))

	tx, err := tbl.DB.db.Begin()

	stmt, err := tx.Prepare(statement)
	if err != nil {
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
		log.Println(err)
		return err
	}
	stmt, err := tx.Prepare("delete from " + tbl.Name + " where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	tx.Commit()

	return err
}

func (tbl *Table) fill(rows *sql.Rows) []interface{} {

	columns, err := rows.Columns()

	values := make([]interface{}, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	results := make([]interface{}, 0)
	idx := 0
	for rows.Next() {
		elem := reflect.New(tbl.Type).Elem()

		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				//TODO () handle nulls
			} else {
				f := elem.Field(i)
				log.Printf("trying to set: %v with %v\n", tbl.Type.Field(i).Name, col)
				tbl.setVal(&f, col)
			}

		}
		if idx > len(results)-1 {
			results = tbl.growResult(results, elem.Interface())
		}
		idx++
	}
	return results
}
func (tbl *Table) growResult(results []interface{}, obj interface{}) []interface{} {
	length := len(results)

	more := make([]interface{}, length+1)

	for i := range results {
		more[i] = results[i]
	}
	more[length] = obj

	return more

}
func (tbl *Table) setVal(field *reflect.Value, val interface{}) {
	str, ok := val.(string)
	if ok {
		field.SetString(str)
		return
	}
	i, ok := val.(int64)
	if ok {
		field.SetInt(i)
		return
	}
	ui, ok := val.(uint64)
	if ok {
		field.SetUint(ui)
		return
	}
	b, ok := val.(bool)
	if ok {
		field.SetBool(b)
		return
	}
	f, ok := val.(float64)
	if ok {
		field.SetFloat(f)
		return
	}
	t, ok := val.(time.Time)
	if ok {
		field.Set(reflect.ValueOf(t))
		return
	}

}
