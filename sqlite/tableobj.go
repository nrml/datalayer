package sqlite

import (
	"database/sql"
	"fmt"
	"reflect"
)

type Fillable interface {
	Fill(*sql.Row, Fillable)
}

type TableObj struct {
}

func (to *Table) fill(rows *sql.Rows, obj interface{}) {
	elem := reflect.ValueOf(obj).Elem()
	//length := elem.NumField() - 1

	columns, err := rows.Columns()

	// Make a slice for the values
	values := make([]interface{}, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
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
				to.setVal(elem.Field(i), col)
			}
		}
		fmt.Println("-----------------------------------")
	}
}
func (to *Table) setVal(field reflect.Value, val interface{}) {
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

}
