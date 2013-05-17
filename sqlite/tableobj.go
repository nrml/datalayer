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

func (to *Table) fill(rows *sql.Rows) []interface{} {

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
		elem := reflect.New(to.Type).Elem()

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
				fmt.Printf("trying to set: %v with %v\n", to.Type.Field(i).Name, col)
				to.setVal(&f, col)
			}

		}
		fmt.Println("checking grow")
		if idx > len(results)-1 {
			//fmt.Printf("growing results with: %v\n", elem.Interface())
			fmt.Println("growing")
			results = to.growResult(results, elem.Interface())
		}
		idx++
	}
	return results
}
func (to *Table) growResult(results []interface{}, obj interface{}) []interface{} {
	length := len(results)

	more := make([]interface{}, length+1)

	for i := range results {
		more[i] = results[i]
	}
	more[length] = obj

	return more

}
func (to *Table) setVal(field *reflect.Value, val interface{}) {
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
