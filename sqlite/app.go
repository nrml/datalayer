package sqlite

// import (
// 	"errors"
// 	"fmt"
// )

// var t interface{}

// type Repo struct {
// 	Type string
// }

// func (this *Repo) Get(id int64) (t, error) {
// 	var reg Registration
// 	stmt, err := this.db.Prepare("select id, email from registration where id = ?")
// 	if err != nil {
// 		fmt.Println("error getting user by id: " + err.Error())
// 		return reg, err
// 	}
// 	defer stmt.Close()
// 	var email string
// 	err = stmt.QueryRow(&id).Scan(&id, &email)
// 	if err != nil {
// 		fmt.Println("error in getbyid function for " + email + " --- " + err.Error())
// 		return reg, err
// 	}
// 	fmt.Println("found", id, email)

// 	reg = Registration{id, email, ""}

// 	return reg, err
// }
// func (this *Repo) GetByEmail(email string) (t, error) {
// 	var reg Registration
// 	stmt, err := this.db.Prepare("select id, email from registration where email = ?")
// 	if err != nil {
// 		fmt.Println("error getting " + email + " " + err.Error())
// 		return reg, err
// 	}
// 	defer stmt.Close()
// 	var id int64
// 	err = stmt.QueryRow(&id).Scan(&id, &email)
// 	if err != nil {
// 		fmt.Println("\nerror in getbyemail function for " + email + " --- " + err.Error())
// 		return reg, err
// 	}
// 	fmt.Println("found", id, email)

// 	reg = Registration{id, email, ""}

// 	return reg, err
// }
// func (this *Repo) List() error {
// 	rows, err := this.db.Query("select id, email, password from registration")
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var id int
// 		var email string
// 		var password string
// 		rows.Scan(&id, &email, &password)
// 		fmt.Println(id, email, password)
// 	}

// 	return errors.New("haven't implemented list of registration")
// 	//TODO: return set of registration

// }
// func (this *Repo) Create(r *Registration) (*t, error) {

// 	_, err := this.GetByEmail(r.Email)

// 	if err == nil {
// 		return r, errors.New("User already exists.")
// 	} else {
// 		fmt.Println("user does not exist, okay to create user -- " + r.Email + " -- error:" + err.Error())
// 	}

// 	tx, err := this.db.Begin()
// 	if err != nil {
// 		fmt.Println(err)
// 		return r, err
// 	}
// 	stmt, err := tx.Prepare("insert into registration(id, email, password) values(?, ?, ?)")
// 	if err != nil {
// 		fmt.Println(err)
// 		return r, err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(nil, r.Email, r.Password)

// 	if err != nil {
// 		fmt.Println("COULD NOT CREATE user -- " + r.Email + " -- error: " + err.Error())
// 	} else {
// 		fmt.Println("CREATED user -- " + r.Email)
// 	}

// 	tx.Commit()

// 	r.Password = ""
// 	return r, err

// }
// func (this *Repo) Update(r *Registration) (*t, error) {
// 	var reg Registration
// 	reg, err := this.GetByEmail(r.Email)

// 	if reg.Email == "" {
// 		return r, errors.New("User does not exist.")
// 	}

// 	tx, err := this.db.Begin()
// 	if err != nil {
// 		fmt.Println(err)
// 		return r, err
// 	}
// 	stmt, err := tx.Prepare("update registration set email=?, password=? where id=?")
// 	if err != nil {
// 		fmt.Println(err)
// 		return r, err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(r.Email, r.Password, r.ID)
// 	tx.Commit()

// 	r.Password = ""
// 	return r, err

// }
// func (this *Repo) Delete(id int64) error {

// 	tx, err := this.db.Begin()
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	stmt, err := tx.Prepare("delete from registration where id=?")
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(id)
// 	tx.Commit()

// 	return err
// }
