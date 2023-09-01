package migrations

import (
	"log"
	"github.com/fcomera/fetch-backend-test/models"
	"github.com/fcomera/fetch-backend-test/db"
)


// Generate the database schema with the models
func MigrateModels() {
	if db.DBConn == nil {
		log.Println("Whoopsie")
	}
	db.DBConn.AutoMigrate(&models.ReceiptRequest{}, &models.Item{})
}