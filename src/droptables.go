package main

import (
	"fmt"
	"gosecondhand/src/database"
	"log"
)

func main() {
	db := database.ConnectDB()
	// Suspend execution of this function until surround has ran
	defer db.Close()

	// ERROR HANDLING FOR DATABASE CONNECTION
	err := db.Ping()
	if err != nil {
		log.Fatal("Connection could not be verified with Ping(): ", err)
	}

	err = database.DropAllTables(db)
	if err != nil {
		log.Fatal("Something went wrong during table dropping call: ", err)
	} else {
		fmt.Println("All tables have been dropped!")
	}
}
