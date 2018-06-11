package utilities

import (
	"fmt"
	"log"
)

func CheckError(err error) bool {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err.Error())
		return false
	}
	return true
}

