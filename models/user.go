package models

import (
	"errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var (
	UserList map[string]*User
)

func init() {

	/* 连接数据库测试 */
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		panic(err.Error())	// proper error handling instead of panic in your app
	}

	// Use the DB normally, execute the querys etc
}

type User struct {
	PhoneId  string
	Username string
	Password string
}

type Profile struct {
	Gender  string
	Age     int
	Address string
	Email   string
}

func AddUser(u User) string {

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare("INSERT INTO user VALUES( ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer statementInsert.Close() // Close the statement when we leave main() / the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(u.PhoneId, u.Username, u.Password)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return u.PhoneId
}

func GetUser(uid string) (u *User, err error) {
	db, err := sql.Open("mysql", "app_root:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statement for reading data
	statementOut, err := db.Prepare("SELECT user_name FROM user WHERE user_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer statementOut.Close()

	// Query the username
	var Username string;
	err = statementOut.QueryRow(uid).Scan(&Username) // WHERE number = uid
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The username of 123 is: %s", Username)
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	if u, ok := UserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func Login(username, password string) bool {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}
