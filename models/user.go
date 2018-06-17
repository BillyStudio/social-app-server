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
	"social-app-server/utilities"
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
	id string
	username string
	password string
}

type Profile struct {
	gender  string
	dob     string
	motto   string
}

func AddUser(u User) (string, error) {

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare("INSERT INTO USER VALUES( ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return "Failed when preparing INSERT statement.", err
	}
	defer statementInsert.Close() // Close the statement when we leave main() / the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(u.id, u.username, u.password)
	if err != nil {
		fmt.Println(err)
		return "Failed when executing INSERT statement.", err
	}
	return u.id, err
}

func AddProfile()  {

}

func GetUser(userId string) (u User, err error) {
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	utilities.CheckError(err)
	defer db.Close()

	var UserItem User
	UserItem.id = userId;

	// Prepare statement for reading data
	RowUserName, err := db.Prepare("SELECT user_name FROM USER WHERE user_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return UserItem, err
	}
	defer RowUserName.Close()

	// Query the username
	err = RowUserName.QueryRow(userId).Scan(&UserItem.username)
	fmt.Printf("Username:%v\n", UserItem.username)

	RowPassword, err := db.Prepare("SELECT password FROM USER WHERE user_id = ?")
	if err != nil {
		fmt.Println(err.Error())
		return UserItem, err
	}
	defer RowPassword.Close()
	err = RowPassword.QueryRow(userId).Scan(&UserItem.password)
	fmt.Printf("Password:%v\n", UserItem.password)

	return UserItem, err
}

func GetAllUsers() (UserList []*User, err error) {
	UserList = make([]*User, 50)	// allocate memory, query at most 50 users once a time

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT * FROM USER")
	if err != nil {
		fmt.Println(err.Error())
		return UserList[0:1], err
	}

	// Get column names
	columns, err := rows.Columns()
	utilities.CheckError(err)

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))
	// rows.Scan wants '[]interface{}' as an argument, so we must copy the references into such a slice
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
			fmt.Println(err)
			return UserList[0:1], err
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
				NewUser.id = value;
			} else if strings.ToUpper(strings.TrimSpace(columns[i])) == "USER_NAME" {
				NewUser.username = value;
			} else if strings.ToUpper(strings.TrimSpace(columns[i])) == "PASSWORD"{
				NewUser.password = value;
			}
			fmt.Printf("new user --> %#v\n ", NewUser)
		}
		fmt.Println("-----------------------------------")
		UserList[iRow] = &NewUser
		iRow = iRow + 1
	}
	UserList = UserList[0:iRow]
	if err = rows.Err(); err != nil {
		return UserList, err
	}

	return UserList, nil
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

	// insert into table TOKEN
	StatementInsert, err := db.Prepare("insert into TOKEN values( ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer StatementInsert.Close() // Close the statement when we leave main() / the program terminates

	// Executing inserting
	_, err = StatementInsert.Exec(token, PhoneId)
	utilities.CheckError(err)

	return token, nil
}

func Logout(Token string) (err error) {

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	StatementRemove, err := db.Prepare("delete from TOKEN where token_id= ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer StatementRemove.Close()

	// Executing deletion
	res, err := StatementRemove.Exec(Token) // WHERE token_id = Token
	if err != nil {
		fmt.Println(err)
		return err
	}
	num, err := res.RowsAffected()
	fmt.Printf("rows affected = %v\n", num)

	return err
}