package models

import (
	"time"
	"database/sql"
	"os"
	"log"
	"strconv"
	"fmt"
	"io/ioutil"
)

// 由客户端上传的Moment
type MomentContent struct {
	Text   string
	Image  string
	Tag    string
}

// 储存在数据库的Moment
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

	/* 将标签、文本和图片均作为文件，存储在res文件夹下 */

	// 存储标签为tag
	if content.Tag != "" {
		TagLocation := "res/" + strconv.FormatInt(m.id, 10) + ".tag"

		f, err := os.OpenFile(TagLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if CheckError(err) {
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

	// 连接数据库
	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare(
		"INSERT INTO MOMENT VALUES(?, ?, ?, ?, ?, ?)")
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

func GetOne(MomentId int64) (content MomentContent, err error) {

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statement for reading data
	statement, err := db.Prepare("SELECT moment_tag,text_location,image_location FROM MOMENT WHERE moment_id = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer statement.Close()
	// Query the username
	var ColumnTag, ColumnText, ColumnImage []byte
	err = statement.QueryRow(MomentId).Scan(&ColumnTag, &ColumnText, &ColumnImage) // WHERE moment_id = MomentId
	tag := string(ColumnTag)
	TextLocation := string(ColumnText)
	ImageLocation := string(ColumnImage)
	content.Tag = tag

	// Get the file content
	if TextLocation != "" {
		BytesText, err := ioutil.ReadFile(TextLocation) // just pass the file name
		if err != nil {
			log.Fatal(err)
		}
		text := string(BytesText)
		fmt.Printf("text: %v\n", text)
		content.Text = text
	}
	if ImageLocation != "" {
		BytesImage, err := ioutil.ReadFile(ImageLocation)
		if err != nil {
			log.Fatal(err)
		}
		image := string(BytesImage)
		fmt.Println("base64 codes of image: %v", image)
		content.Image = image
	}

	return content, err
}

func GetAll() map[int64]*Moment {
	var Moments map[int64]*Moment
	Moments = make(map[int64]*Moment)	// allocate memory
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
	for rows.Next() {
		var moment Moment
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		CheckError(err)
		// Now do something with the data
		var value string
		for i, col := range values {
			fmt.Println(i)
			value = string(col)
			switch i {
			case 0: moment.id, err = strconv.ParseInt(value, 10, 64)
				if err != nil {
					panic(err.Error())
				}
				break
			case 1: moment.PublishTime = value
				break
			case 2: moment.Tag = value
				break
			case 3: moment.TextLocation = value
				break
			case 4: moment.ImageLocation = value
				break
			default:
				break
			}
		}
		Moments[moment.id] = &moment
	}

	// Fetch rows

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

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
		return false
	}
	return true
}
