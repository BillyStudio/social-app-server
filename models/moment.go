package models

import (
	"errors"
	"time"
	"database/sql"
	"os"
	"log"
	"strconv"
	"fmt"
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
	Tag	           string
	TextLocation  string
	ImageLocation string
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

func AddOne(content MomentContent) (MomentId int64) {
	// 将发送时间作为id
	// LEAVE：将除上发送人的余数作为id可以避免1秒内的碰撞
	MomentId = time.Now().UTC().UnixNano()
	var m Moment
	m.id = MomentId
	m.PublishTime = time.Now().Format("2006-01-02 15:04:05")	// 2006-01-02 15:04:05据说是Go的诞生时间
	m.Tag = content.Tag

	/* 将文本与图片作为文件存储 */

	// 存储文本为txt
	if content.Text != "" {
		m.TextLocation = "res/" + strconv.FormatInt(MomentId, 10) + ".txt"
		f, err := os.OpenFile(m.TextLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(content.Text)); err != nil {
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	// 存储图片为img
	if content.Image != "" {
		m.ImageLocation = "res/" + strconv.FormatInt(MomentId, 10) + ".img"
		f, err := os.OpenFile(m.ImageLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Write([]byte(content.Image)); err != nil {	// LEAVE: 存储之后再读回来
			log.Fatal(err)
		}
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}

	/* 储存 m 到数据库中 */

	db, err := sql.Open("mysql", "ubuntu:IS1501@/social_app")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// Prepare statements for inserting data
	statementInsert, err := db.Prepare(
		"INSERT INTO MOMENT VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer statementInsert.Close() // Close the statement when we leave main()/the program terminates

	// Executing inserting
	_, err = statementInsert.Exec(m.id, m.PublishTime, m.Tag, m.TextLocation, m.ImageLocation)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	return MomentId
}

func GetOne(MomentId int) (moment *Moment, err error) {

	return nil, errors.New("ObjectId Not Exist")
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
	/*columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}*/

	// 建立interface到slice的索引，values中存储每一行的数据
	/*values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}*/
	// 按行读取
	for rows.Next() {
		var moment Moment
		// get RawBytes from data
		err = rows.Scan(&moment.id, &moment.PublishTime, &moment.Tag, &moment.TextLocation, &moment.ImageLocation)
		fmt.Println("moment:%v", moment);
		if err != nil {
			fmt.Println(err);
			panic(err.Error())
		}
		Moments[moment.id] = &moment
	}

	// Fetch rows

	return Moments
}

func Delete(MomentId int64) {
	var Moments map[int64]*Moment
	delete(Moments, MomentId)
}

