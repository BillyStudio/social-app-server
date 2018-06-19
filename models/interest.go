package models

import (
	"database/sql"
	"fmt"
)


// 储存在数据库的 UserBasic interest
type UserInterest struct {
	UserId		string
	InterestTag	string
}

// 储存在数据库的 Interest area
type Area struct {
	InterestTag string
	MomentId	int64
}

// 返回给前端的 Interest area
type InterestArea struct {
	AreaName	string
	MomentStrId	string
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
	if err != nil {
		fmt.Println(err.Error());
		return false
	}
	defer statementInsert.Close() // Close the statement when we leave main()/the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(Tags, MomnentId)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func GetAllInterests() (InterestAreas []*InterestArea) {
	InterestAreas = make([]*InterestArea, 50)

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	Rows, err := db.Query("select interest_tag, fk_moment_id from AREA")
	iRow := 0
	for Rows.Next() {
		var area InterestArea
		var ColTag, ColMoment []byte
		err = Rows.Scan(&ColTag, &ColMoment)
		area.AreaName = string(ColTag)
		area.MomentStrId = string(ColMoment)
		fmt.Println("The Moment ID = ", area.MomentStrId)
		InterestAreas[iRow] = &area
		iRow ++
	}

	InterestAreas = InterestAreas[0:iRow]
	return InterestAreas
}

func UploadInterest(userInterest UserInterest) bool {
	// LEAVE: Split the tags

	// Check the database
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare("INSERT INTO INTEREST VALUES(?, ?)")
	if err != nil {
		fmt.Println(err.Error());
		return false
	}
	defer statementInsert.Close() // Close the statement when we leave main()/the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(userInterest.UserId, userInterest.InterestTag)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
