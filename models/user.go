package models

import (
	"errors"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
	"crypto/md5"
	"io"
	"strconv"
	"time"
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

	// Use the DB normally, execute the queries etc
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
	statementInsert, err := db.Prepare("INSERT INTO USER VALUES( ?, ?, ?)")
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

func GetUser(PhoneId string) (u User, err error) {
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var UserItem User
	UserItem.PhoneId = PhoneId;
	fmt.Printf("PhoneId: %v\n", UserItem.PhoneId)
	// Prepare statement for reading data
	RowUserName, err := db.Prepare("SELECT user_name FROM USER WHERE user_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer RowUserName.Close()
	// Query the username
	err = RowUserName.QueryRow(PhoneId).Scan(&UserItem.Username) // WHERE number = uid
	fmt.Printf("Username:%v\n", UserItem.Username)

	RowPassword, err := db.Prepare("SELECT password FROM USER WHERE user_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer RowPassword.Close()
	err = RowPassword.QueryRow(PhoneId).Scan(&UserItem.Password)
	fmt.Printf("Password:%v\n", UserItem.Password)

	return UserItem, err
}

func GetAllUsers() []User {
	var UserList []User
	UserList = make([]User, 20)	// allocate memory

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT * FROM USER")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	iRow := 0
	for rows.Next() {
		NewUser := User{};

		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)

			if strings.ToUpper(strings.TrimSpace(columns[i])) == "USER_ID" {
				NewUser.PhoneId = value;
			} else if strings.ToUpper(strings.TrimSpace(columns[i])) == "USER_NAME" {
				NewUser.Username = value;
			} else if strings.ToUpper(strings.TrimSpace(columns[i])) == "PASSWORD"{
				NewUser.Password = value;
			}
			fmt.Printf("new user --> %#v\n ", NewUser)
		}
		fmt.Println("-----------------------------------")
		UserList[iRow] = NewUser
		iRow = iRow + 1
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	return nil, errors.New("User Not Exist")
}

func Login(PhoneId, password string) (token string, err error) {
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	StatementFind, err := db.Prepare("select user_name from USER where user_id= ? and password= ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer StatementFind.Close()

	// Query the username
	row := StatementFind.QueryRow(PhoneId, password)
	var username string;
	err = row.Scan(&username);
	fmt.Printf("Match username: %v\n", username)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
		return "", err;
	}

	// generate tokens
	h := md5.New()
	fmt.Println("h-->%v", h)
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	fmt.Println("h-->%v", h)

	token = fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println("token-->%v", token)
	return token, nil
}

