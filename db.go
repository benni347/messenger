package main

import (
	"database/sql"

	"github.com/benni347/encryption"
	utils "github.com/benni347/messengerutils"
	_ "github.com/go-sql-driver/mysql"
)

func database(msg string, chatId uint64) {
	m := &utils.MessengerUtils{
		Verbose: true,
	}
	db, err := sql.Open("mysql", "root:password")
	if err != nil {
		utils.PrintError("During the opening of the sql databse an error ocurerd:", err)
	}
	defer db.Close()

	// Conect and print the version.
	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		utils.PrintError("During the query of the version an error occured:", err)
	}
	m.PrintInfo("Version:", version)

	hash := encryption.CalculateHash([]byte(msg))
	m.PrintInfo("Hash:", hash)
}
