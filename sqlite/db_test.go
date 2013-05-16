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
	_     TableObj
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

func Test_TableCreateGetUpdate(t *testing.T) {
	db := DB{":memory:", nil}
	db.Init()

	tbl := Table{
		"TestTable",
		fields,
		db,
	}

	_ = db.CreateTable(tbl)

	obj := TestTableObject{1, "Object Name", "Object Title", TableObj{}}

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

	obj2 := new(TestTableObject)

	fmt.Println("attempting to fill " + strconv.FormatInt(obj.ID, 10))
	err = tbl.Fill(id, obj2)
	fmt.Printf("made it past fill with object: %d-%s-%s\n", obj2.ID, obj2.Name, obj2.Value)
	//err = tbl.Get(obj.ID, &obj2.ID, &obj2.Name, &obj2.Value)

	if err != nil {
		t.Error("error getting table object... " + err.Error())
	} else {
		t.Log("object got with name: " + obj2.Name)
		fmt.Println("object got with name: " + obj2.Name)
	}

	obj2.Name = "Object Name Updated"

	err = tbl.Update(obj2.ID, obj2.Name, obj2.Value)

	if err != nil {
		t.Error("error updating table object... " + err.Error())
		return
	} else {
		t.Log("object updated: " + obj2.Name)
		fmt.Println("object updated wth name : " + obj2.Name)
	}

	var obj3 = new(TestTableObject)

	err = tbl.Get(obj2.ID, &obj3.ID, &obj3.Name, &obj3.Value)

	if err != nil {
		t.Error("error getting table object after update... " + err.Error())
	} else {
		t.Log("updated object got with name: " + obj3.Name)
		fmt.Println("updated object got with name: " + obj3.Name)
	}

	err = tbl.Delete(obj3.ID)

	if err != nil {
		t.Error("error deleting... " + err.Error())
	} else {
		t.Log("obj deleted with id: " + strconv.FormatInt(obj3.ID, 10))
		fmt.Println("object deleted with id: " + strconv.FormatInt(obj3.ID, 10))
	}

	var obj4 = new(TestTableObject)

	err = tbl.Get(obj3.ID, &obj4.ID, &obj4.Name, &obj4.Value)

	if err != nil {
		t.Log("confirmed delete")
		fmt.Println("confirmed delete")
	} else {
		t.Error("object did not delete with name " + obj4.Name)
		fmt.Println("object did not delete with name: " + obj4.Name)
	}

}
