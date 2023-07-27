package main

import (
	"log"
	"server/db"
)

func main() {
	db, err := db.NewDatabase()

	if err != nil {
		log.Println(err)
	}
	db.GetDB()
}
