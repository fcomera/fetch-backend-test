package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/fcomera/fetch-backend-test/controllers"
	"github.com/fcomera/fetch-backend-test/db"
	"github.com/fcomera/fetch-backend-test/migrations"
)



func main() {
	app := fiber.New()

	db.InitializeDatabase()
	migrations.MigrateModels()

	api := app.Group("/api")

	api.Post("/receipts/process", controllers.NewReceipt)
	api.Get("/receipts/:id/points", controllers.ReceiptPoints)
	app.Listen(":3000")
}