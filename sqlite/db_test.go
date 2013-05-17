package sqlite

import (
	"fmt"
	"strconv"
	"testing"
)

type TestObject struct {
	ID    int64  `json:"id" sql:"id integer not null primary key autoincrement"`
	Name  string `json:"name"  sql:"name text"`
	Value string `json:"value" sql:"value text"`
	//_     TableObj
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

	_, err := db.CreateTable("TestTable", TestObject{})

	if err != nil {
		t.Error("error creating table... " + err.Error())
	} else {
		t.Log("table was created")
	}

}

func Test_TableCreateGetUpdate(t *testing.T) {
	db := DB{":memory:", nil}
	db.Init()

	tbl, err := db.CreateTable("TestTable", TestObject{})

	obj := TestObject{Name: "Object Name", Value: "Object Title"}

	defer db.Close()
	fmt.Println("about to create")
	id, err := tbl.Create(obj)

	if err != nil {
		t.Error("error creating table object... " + err.Error())
		return
	} else {
		t.Log("object created with id: " + strconv.FormatInt(id, 10))
		fmt.Println("object created with id: " + strconv.FormatInt(id, 10))
	}

	obj.ID = id

	fmt.Println("attempting to fill " + strconv.FormatInt(obj.ID, 10))

	res, err := tbl.Get(id)
	obj2 := res.(TestObject)

	fmt.Printf("made it past fill with object: %d-%s-%s\n", obj2.ID, obj2.Name, obj2.Value)
	//err = tbl.Get(obj.ID, &obj2.ID, &obj2.Name, &obj2.Value)

	if err != nil {
		t.Error("error getting table object... " + err.Error())
		return
	} else {
		t.Log("object got with name: " + obj2.Name)
		fmt.Println("object got with name: " + obj2.Name)
	}

	obj2.Name = "Object Name Updated"

	err = tbl.Update(obj2.ID, obj2)

	if err != nil {
		t.Error("error updating table object... " + err.Error())
		return
	} else {
		t.Log("object updated: " + obj2.Name)
		fmt.Println("object updated wth name : " + obj2.Name)
	}

	res, err = tbl.Get(obj2.ID)
	obj3 := res.(TestObject)

	if err != nil {
		t.Error("error getting table object after update... " + err.Error())
		return
	} else {
		t.Log("updated object got with name: " + obj3.Name)
		fmt.Println("updated object got with name: " + obj3.Name)
	}

	list, err := tbl.List()

	if err != nil {
		t.Error("error getting listing objects after update... " + err.Error())
		return
	} else {
		fmt.Printf("found %v objects in list\n", len(list))
		tobj1 := list[0].(TestObject)
		fmt.Printf("found object1 in list:  %v\n", tobj1)

	}

	err = tbl.Delete(obj3.ID)

	if err != nil {
		t.Error("error deleting... " + err.Error())
		return
	} else {
		t.Log("obj deleted with id: " + strconv.FormatInt(obj3.ID, 10))
		fmt.Println("object deleted with id: " + strconv.FormatInt(obj3.ID, 10))
	}

	res, err = tbl.Get(obj3.ID)

	if res == nil {
		t.Log("confirmed delete")
		fmt.Println("confirmed delete")
	} else {
		obj4 := res.(TestObject)
		t.Error("object did not delete with name " + obj4.Name)
		fmt.Println("object did not delete with name: " + obj4.Name)
	}

}
