package models

import (
	"time"
	"log"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"github.com/fcomera/fetch-backend-test/db"
	"unicode"
	"math"
	"strings"
)

// Receipt Model
type ReceiptRequest struct {
	ID uuid.UUID `json:"id, omitempty" gorm:"primaryKey"`
	Retailer string `json:"retailer"`
	PurchaseDateStr string `json:"purchaseDate"`
	PurchaseTimeStr string `json:"purchaseTime"`
	PurchaseDateAndTime time.Time `json:"-"`
	Total float32 `json:"total,string"`
	Items []Item `json:"items" gorm:"foreignKey:ReceiptRequestRefer"`
	Points int `json:"-" gorm:"default:-1"`
}


// Item Model
type Item struct {
	gorm.Model
	ReceiptRequestRefer uuid.UUID `json:"-"`
	ShortDescription string `json:"shortDescription"`
	Price float32 `json:"price,string"`
}


// Generates a new UUID for an unidentified receipt
func CreateNewUUIDForReceipt(r *ReceiptRequest) {
	newUuid := uuid.New()
	r.ID = newUuid
}

// Queries the database to find a receipt by an UUID
func FindReceiptById(id uuid.UUID) ReceiptRequest {
	var receipt ReceiptRequest = ReceiptRequest{ID: id}

	result := db.DBConn.Model(&ReceiptRequest{}).Preload("Items").First(&receipt)

	if result.Error != nil {
		log.Println(result.Error)
	}

	return receipt
}


// Calculates the points of a receipt and
// stores the receipt and its items on the
// database
func CreateReceipt(r *ReceiptRequest){

	dbo := db.DBConn

	if dbo == nil {
		log.Println(dbo)
	}

	r.Points = CalculatePoints(*r)

	tx := dbo.Begin()
	tx.Create(r)
	tx.Commit()

}


// Updates the points column of a stored receipt
func UpdateReceiptPoints(r ReceiptRequest, points int) {
	var receipt *ReceiptRequest = &r
	receipt.Points = points
	db.DBConn.Model(receipt).Update("points", points)
}

// Obtains the alphanumeric value from the retailer name of a receipt
func checkName(retailerName string) int {
	points := 0

	for _, c := range retailerName {
		if unicode.IsDigit(c) || unicode.IsLetter(c) {
			points += 1
		}
	}

	return points
}


// Verifies the total of the receipt is rounded in dollars(with no cents)
// Returns 50 points if that is the case
func checkRoundDollar(price float32) int {
	points := 0
	intVal := int(price)
	isRound := price - float32(intVal)

	if !(isRound > 0.00) {
		points += 50
	}

	return points
}


// Verifies the total of a receipt is a multiple of 0.25 cents.
// Returns 25 points if that is the case
func checkPriceIsMultipleOfQuarter(price float32) int {
	points := 0

	isMultiple := (math.Mod(float64(price), 0.25) == 0.00)

	if isMultiple {
		points += 25
	}

	return points
}


// Returns 5 points for each 2 items in the receipt
func checkItems(items []Item) int {
	points := len(items) / 2

	points = points * 5

	return points
}


// Checks and verifies if the description of an item
// is multiple of 3 and if that is the case it adds
// the ceil value of the item price times(*) 0.2
func checkItemsLength(items []Item) int {
	points := 0

	for _, item := range items {
		description := strings.TrimSpace(item.ShortDescription)
		stringLenIsMultipleOfThree := (len(description) % 3 == 0)
		if stringLenIsMultipleOfThree {
			points += int(math.Ceil(float64(item.Price) * 0.2))
		}
	}

	return points
}


// Verifies the day of the purchase is odd. Returns 6 points if
// that is the case
func checkPurchaseDayIsOdd(purchaseDateAndTime time.Time) int {
	points := 0
	day := purchaseDateAndTime.Day()
	isOdd := day % 2 != 0

	if isOdd {
		points += 6
	}

	return points
}


// Checks if the purchase time is between 2:00 pm and 4:00 pm. If
// that is the case it returns 10 points
func checkPurchaseTime(purchaseDateAndTime time.Time) int {
	points := 0

	startTime, _ := time.Parse(time.Kitchen, "2:00PM")
	endTime, _ := time.Parse(time.Kitchen, "4:00PM")
	purchaseTime, _ := time.Parse(time.Kitchen, purchaseDateAndTime.Format(time.Kitchen))

	if purchaseTime.After(startTime) && purchaseTime.Before(endTime) {
		points += 10
	}

	return points
}


// Pass through all the checks to determine the points of a receipt
func CalculatePoints(r ReceiptRequest) int{
	points := 0

	points += checkName(r.Retailer)
	points += checkRoundDollar(r.Total)
	points += checkPriceIsMultipleOfQuarter(r.Total)
	points += checkItems(r.Items)
	points += checkItemsLength(r.Items)
	points += checkPurchaseDayIsOdd(r.PurchaseDateAndTime)
	points += checkPurchaseTime(r.PurchaseDateAndTime)

	return points
}