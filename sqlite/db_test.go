package sqlite

import (
	"fmt"
	"strconv"
	"testing"
)

var fields = []Field{
	{"id", "integer not null primary key autoincrement"},
	{"name", "text"},
	{"value", "text"},
}

type TestTableObject struct {
	ID    int64
	Name  string
	Value string
}

func Test_InitDB(t *testing.T) {
	db := DB{":memory:", nil}
	db.Init()

	if db.Namespace != ":memory:" {
		t.Error("no db namspace")
	} else {
		t.Log("correctly assigned namespace memory")
	}
}
func Test_CreateTable(t *testing.T) {
	db := DB{":memory:", nil}
	db.Init()

	tbl := Table{
		"TestTable",
		fields,
		db,
	}

	if tbl.DB.Namespace != ":memory:" {
		t.Error("table namespace is incorrect")
		fmt.Println("table namespace error")
	} else {
		t.Log("table namespce is correct")
		fmt.Println("table namespace correct")
	}

	err := db.CreateTable(tbl)

	if err != nil {
		t.Error("error creating table... " + err.Error())
	} else {
		t.Log("table was created")
	}

}

func Test_TableCreateGet(t *testing.T) {
	db := DB{":memory:", nil}
	db.Init()

	tbl := Table{
		"TestTable",
		fields,
		db,
	}

	_ = db.CreateTable(tbl)

	obj := TestTableObject{1, "Object Name", "Object Title"}

	values := []string{"null", obj.Name, obj.Value}

	defer db.Close()

	id, err := tbl.Create(values)

	if err != nil {
		t.Error("error creating table object... " + err.Error())
	} else {
		t.Log("object created with id: " + strconv.FormatInt(id, 10))
		fmt.Println("object created with id: " + strconv.FormatInt(id, 10))
	}

	obj.ID = id

	var obj2 = new(TestTableObject)

	err = tbl.Get(obj.ID, &obj2.ID, &obj2.Name, &obj2.Value)

	if err != nil {
		t.Error("error getting table object... " + err.Error())
	} else {
		t.Log("object got with name: " + obj2.Name)
		fmt.Println("object got with name: " + obj2.Name)
	}

}
