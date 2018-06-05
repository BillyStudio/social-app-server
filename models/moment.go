package models

import (
	"errors"
	"time"
	"database/sql"
	"os"
	"log"
	"strconv"
)

var (
	Moments map[int]*Moment
)

// 由客户端上传的Moment
type MomentContent struct {
	Text   string
	Image  string
	Tag    string
}

// 储存在数据库的Moment
type Moment struct {
	id			   int
	PublishTime    time.Time
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

func AddOne(content MomentContent) (MomentId int) {
	// 将发送时间作为id
	// LEAVE：将除上发送人的余数作为id可以避免1秒内的碰撞
	MomentId = time.Now().Second()
	var m Moment
	m.id = MomentId
	m.PublishTime = time.Now()
	m.Tag = content.Tag

	/* 将文本与图片作为文件存储 */

	// 存储文本为txt
	if content.Text != "" {
		m.TextLocation = strconv.Itoa(MomentId) + ".txt"
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
		m.ImageLocation = strconv.Itoa(MomentId) + ".img"
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
		"INSERT INTO MOMENT VALUES(moment_id=?,moment_time=?,moment_tag=?,text_location=?,image_location=?")
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

func GetAll() map[int]*Moment {
	return Moments
}

func Delete(MomentId int) {
	delete(Moments, MomentId)
}

