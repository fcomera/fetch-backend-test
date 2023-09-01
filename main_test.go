package main

import (
	"fmt"
	"bytes"
	"encoding/json"
	"github.com/fcomera/fetch-backend-test/controllers"
	"github.com/fcomera/fetch-backend-test/db"
	"github.com/fcomera/fetch-backend-test/migrations"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"testing"

)

type TestCaseReceipt struct {
	Retailer       string         `json:"retailer"`
	PurchaseDate   string         `json:"purchaseDate"`
	PurchaseTime   string         `json:"purchaseTime"`
	Total          string         `json:"total"`
	Items          []TestCaseItem `json:"items"`
	ExpectedPoints int            `json:"-"`
}

type TestCaseItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func GetResponse(app *fiber.App, method string, url string, requestBody any) map[string]any {
	var responseBody map[string]any
	var err error

	if requestBody == nil {
		requestBody = []byte("")
	}

	bodyJson, err := json.Marshal(requestBody)

	if err != nil {
		log.Println(err)
	}

	req := httptest.NewRequest(method, url, bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, 50)

	bodyData := make([]byte, resp.ContentLength)
	_, _ = resp.Body.Read(bodyData)
	err = json.Unmarshal(bodyData, &responseBody)

	if err != nil {
		log.Println(err)
	}

	return responseBody
}

func TestFiberApp(t *testing.T) {
	app := fiber.New()

	db.InitializeDatabase()
	migrations.MigrateModels()

	api := app.Group("/api")

	api.Post("/receipts/process", controllers.NewReceipt)
	api.Get("/receipts/:id/points", controllers.ReceiptPoints)

	// Test Cases
	receiptsTestCases := []TestCaseReceipt{
		TestCaseReceipt{
			Retailer:     "Walgreens",
			PurchaseDate: "2022-01-02",
			PurchaseTime: "08:13",
			Total:        "9.0",
			Items: []TestCaseItem{
				TestCaseItem{
					ShortDescription: "Pepsi - 12-oz",
					Price:            "1.25",
				},
				TestCaseItem{
					ShortDescription: "Dasani",
					Price:            "1.40",
				},
			},
			ExpectedPoints: 90,
		},
		TestCaseReceipt{
			Retailer:     "Target",
			PurchaseDate: "2022-01-02",
			PurchaseTime: "13:13",
			Total:        "1.25",
			Items: []TestCaseItem{
				TestCaseItem{
					ShortDescription: "Pepsi - 12-oz",
					Price:            "1.25",
				},
			},
			ExpectedPoints: 31,
		},
		TestCaseReceipt{
			Retailer:     "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Total:        "35.35",
			Items: []TestCaseItem{
				TestCaseItem{
					ShortDescription: "Mountain Dew 12PK",
					Price:            "6.49",
				}, TestCaseItem{
					ShortDescription: "Emils Cheese Pizza",
					Price:            "12.25",
				}, TestCaseItem{
					ShortDescription: "Knorr Creamy Chicken",
					Price:            "1.26",
				}, TestCaseItem{
					ShortDescription: "Doritos Nacho Cheese",
					Price:            "3.35",
				}, TestCaseItem{
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            "12.00",
				},
			},
			ExpectedPoints: 28,
		},
		TestCaseReceipt{
			Retailer:     "M&M Corner Market",
			PurchaseDate: "2022-03-20",
			PurchaseTime: "14:33",
			Total:        "9.00",
			Items: []TestCaseItem{
				TestCaseItem{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				}, TestCaseItem{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				}, TestCaseItem{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				}, TestCaseItem{
					ShortDescription: "Gatorade",
					Price:            "2.25",
				},
			},
			ExpectedPoints: 109,
		},
	}

	for _, element := range receiptsTestCases {
		body := GetResponse(app, "POST", "/api/receipts/process", element)
		assert.NotNil(t, body)
		assert.Contains(t, body, "id")
		elementToCalculate := body["id"]
		pointsRoute := fmt.Sprintf("/api/receipts/%s/points", elementToCalculate)
		body = GetResponse(app, "GET", pointsRoute, nil)
		assert.NotNil(t, body)
		assert.Contains(t, body, "points")
		var bodyPoints interface{} = int(body["points"].(float64))
		assert.Equal(t, element.ExpectedPoints, bodyPoints.(int))
	}

}
