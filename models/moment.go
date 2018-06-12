package models

import (
	"time"
	"database/sql"
	"os"
	"log"
	"strconv"
	"fmt"
	"io/ioutil"
	"social-app-server/utilities"
)

// 由客户端上传的Moment
type MomentContent struct {
	Token  string
	Text   string
	Image  string
	Tag    string
}

// 返回给客户端的Moment
type MomentReturn struct {
	StrId	string
	Time 	string
	Tags    string
	Text    string
	Image   string
	User    string
	Likes   string
}

// 与数据库交互的Moment，注意数据库中存储的id是长整型
type Moment struct {
	id			   int64
	PublishTime    string
	IfTag	       bool
	IfText         bool
	IfImage        bool
	ForeignKeyUser string
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

func AddOne(content MomentContent) int64 {
	// 将发送时间作为id
	var m Moment
	m.id = time.Now().UTC().UnixNano()
	fmt.Printf("m.id=%v\n", m.id)

	m.PublishTime = time.Now().Format("2006-01-02 15:04:05")	// 2006-01-02 15:04:05据说是Go的诞生时间

	/* 首先检查token是否存在，即将token与user_id对应 */
	// 连接数据库
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statement for reading data
	statementQuery, err := db.Prepare("SELECT fk_user_id FROM TOKEN WHERE token_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer statementQuery.Close()
	// Query the foreign key of user id
	var ColumnUserId []byte
	err = statementQuery.QueryRow(content.Token).Scan(&ColumnUserId) // WHERE token_id = Token
	utilities.CheckError(err)
	m.ForeignKeyUser = string(ColumnUserId)

	/* 将标签、文本和图片均作为文件，存储在res文件夹下 */

	// 存储标签为tag，并添加到AREA数据表中
	if content.Tag != "" {
		TagLocation := "res/" + strconv.FormatInt(m.id, 10) + ".tag"

		AddInterest(content.Tag, m.id)

		f, err := os.OpenFile(TagLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if utilities.CheckError(err) {
			if _, err := f.Write([]byte(content.Tag)); err != nil {
				log.Fatal(err.Error())
			}
		}
		if err := f.Close(); err != nil {
			log.Fatal(err.Error())
		}
		m.IfTag = true
	} else {
		m.IfTag = false
	}

	// 存储文本为txt
	if content.Text != "" {
		TextLocation := "res/" + strconv.FormatInt(m.id, 10) + ".txt"
		f, err := os.OpenFile(TextLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(content.Text)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		m.IfText = true
	} else {
		m.IfText = false
	}

	// 存储图片为img
	if content.Image != "" {
		ImageLocation := "res/" + strconv.FormatInt(m.id, 10) + ".img"
		f, err := os.OpenFile(ImageLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(content.Image)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
		m.IfImage = true
	} else {
		m.IfImage = false
	}

	/* 储存 m 到数据库中 */

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare(
		"INSERT INTO MOMENT VALUES(?, ?, ?, ?, ?, ?, 0)")
	if err != nil {
		panic(err.Error())
	}
	defer statementInsert.Close() // Close the statement when we leave main()/the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(m.id, m.PublishTime, m.IfTag, m.IfText, m.IfImage, m.ForeignKeyUser)
	if err != nil {
		fmt.Println(err)
	}

	return m.id
}

func GetOne(Token string, MomentId int64) (content MomentContent, err error) {

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	content.Token = Token
	// LEAVE: Check if the token really exixts in the table

	// Prepare statement for reading data
	statement, err := db.Prepare("SELECT if_tag, if_text, if_image FROM MOMENT WHERE moment_id = ?")
	utilities.CheckError(err)
	defer statement.Close()

	// Query
	var IfTag, IfText, IfImage bool
	err = statement.QueryRow(MomentId).Scan(&IfTag, &IfText, &IfImage) // WHERE moment_id = MomentId
	fmt.Println("IfTag:-->", IfTag)

	// Get the file content
	if IfTag {
		TagLocation := "res/" + strconv.FormatInt(MomentId, 10) + ".tag"
		BytesTag, err := ioutil.ReadFile(TagLocation)
		utilities.CheckError(err)
		Tags := string(BytesTag)
		fmt.Printf("Tags:%v\n", Tags)
		content.Tag = Tags
	}
	if IfText {
		TextLocation := "res/" + strconv.FormatInt(MomentId, 10) + ".txt"
		BytesText, err := ioutil.ReadFile(TextLocation) // just pass the file name
		if err != nil {
			log.Fatal(err)
		}
		text := string(BytesText)
		fmt.Printf("text: %v\n", text)
		content.Text = text
	}
	if IfImage {
		ImageLocation := "res/" + strconv.FormatInt(MomentId, 10) + ".img"
		BytesImage, err := ioutil.ReadFile(ImageLocation)
		if err != nil {
			log.Fatal(err)
		}
		image := string(BytesImage)
		fmt.Println("base64 string of image: ", image)
		content.Image = image
	}
	return content, err
}

func GetAll() []*MomentReturn {
	var Moments []*MomentReturn
	Moments = make([]*MomentReturn, 50) // pre-allocate memory
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Execute the query
	rows, err := db.Query("SELECT * FROM MOMENT")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	// 建立interface到slice的索引，values中存储每一行的数据
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	// 按行读取
	iRow := 0
	for rows.Next() {
		var moment MomentReturn
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		utilities.CheckError(err)
		// Now do something with the data
		var value string
		for i, col := range values {
			fmt.Println(i)
			value = string(col)
			switch i {
			case 0: moment.StrId = value
				fmt.Println("The Moment Id is ", moment.StrId)
			break
			case 1: moment.Time = value
			break
			case 2:
				IfTag, err := strconv.ParseBool(value)
				utilities.CheckError(err)
				if IfTag {
					TagLocation := "res/" + moment.StrId + ".tag"
					BytesTag, err := ioutil.ReadFile(TagLocation)
					utilities.CheckError(err)
					Tags := string(BytesTag)
					fmt.Printf("Tags:%v\n", Tags)
					moment.Tags = Tags
				} else {
					moment.Tags = ""
				}
			break
			case 3: moment.Text = value
			break
			case 4: moment.Image = value
			break
			case 5: moment.User = value
			break
			case 6: moment.Likes = value
			default:
				break
			}
		}
		Moments[iRow] = &moment
		iRow ++
	}
	Moments = Moments[0:iRow]
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return Moments
}

func Delete(MomentId int64) bool {
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statement for reading data
	statement, err := db.Prepare("delete from MOMENT where moment_id = ?")
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
	}
	defer statement.Close()

	// Executing deletion
	res, err := statement.Exec(MomentId) // WHERE moment_id = MomentId
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
		return false
	}
	num, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
		return false
	}
	fmt.Println(num)

	return true
}

