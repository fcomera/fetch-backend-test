package db

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	//"os"
)

var DBConn *gorm.DB


// Initialize the SQLite database
func InitializeDatabase() {
	log.Println("Initializing database")
	var err error

	DBConn, err = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if DBConn == nil {
		panic("Database is not prepared")
	}

	if err != nil {
		log.Println("Error %s", err.Error())
		panic("Something went wrong")
	}
}
