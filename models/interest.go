package models

import (
	"database/sql"
	"social-app-server/utilities"
	"fmt"
)


// 储存在数据库的 User interest
type INTEREST struct {
	UserId		string
	InterestTag	string
}
// 储存在数据库的 Interest area
type AREA struct {
	InterestTag string
	MomentId	int64
}

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
}


func AddInterest(Tags string, MomnentId int64) bool {
	// LEAVE: Split the tags

	// Check the database
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare("INSERT INTO AREA VALUES(?, ?)")
	utilities.CheckError(err)
	defer statementInsert.Close() // Close the statement when we leave main()/the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(Tags, MomnentId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}
