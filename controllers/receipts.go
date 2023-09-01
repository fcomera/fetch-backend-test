package controllers

import (
	"log"
	"github.com/google/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/fcomera/fetch-backend-test/models"
	"github.com/fcomera/fetch-backend-test/utils"
)


// Attends the request to create a new receipt and calculate the points
// associated  to it. Also generates a UUID to be used as
// a unique identifier of the receipt
func NewReceipt(c *fiber.Ctx) error {
	re := new(models.ReceiptRequest)

	if err := c.BodyParser(re); err != nil {
		return err
	}

	re.PurchaseDateAndTime = utils.ParseDateAndTime(re.PurchaseDateStr, re.PurchaseTimeStr)

	models.CreateNewUUIDForReceipt(re)
	models.CreateReceipt(re)

	response := map[string]uuid.UUID{
		"id": re.ID,
	}

	return c.JSON(response)
}


// Attends the request that retrieves the points of a receipt.
func ReceiptPoints(c *fiber.Ctx) error {
	idToSearch := c.Params("id")
	receiptId, err := uuid.Parse(idToSearch)

	if err != nil {
		log.Print("Error %s", err.Error())
	}

	receipt := models.FindReceiptById(receiptId)

	response := map[string]int{
		"points": receipt.Points,
	}

	return c.JSON(response)
}