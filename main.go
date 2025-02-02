//Handles receipt processing and points calculation based rules listed

package main

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Data Models

// Receipt represents the incoming JSON structure from API reqs

type Receipt struct {
	Retailer     string  `json:"retailer" binding:"required"`
	PurchaseDate string  `json:"purchaseDate" binding:"required"`
	PurchaseTime string  `json:"purchaseTime" binding:"required"`
	Items        []Item  `json:"items" binding:"required,min=1"`
	Total        string  `json:"total" binding:"required"`
}

// Item represents single product on receipt

type Item struct {
	ShortDescription string `json:"shortDescription" binding:"required"`
	Price            string `json:"price" binding:"required"`
}

// ReceiptDocument represents HOW storage happens in MongoDB

type ReceiptDocument struct {
	ID           string  `bson:"_id"`
	Retailer     string  `bson:"retailer"`
	PurchaseDate string  `bson:"purchaseDate"`
	PurchaseTime string  `bson:"purchaseTime"`
	Items        []Item  `bson:"items"`
	Total        string  `bson:"total"`
}

//Global Variables

var (
	mongoClient *mongo.Client
	collection  *mongo.Collection
)

//Main App Entry Pt

func main() {
	// Initialize MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	defer mongoClient.Disconnect(ctx)

	collection = mongoClient.Database("receipts").Collection("receipts")

// Initialize Gin router

	r := gin.Default()

	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)

	r.Run(":8080")
}

//API Handlers

//POST

func processReceipt(c *gin.Context) {
	var receipt Receipt
	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt. Please verify input."})
		return
	}

	// Validate receipt fields
	if !validateReceipt(receipt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receipt. Please verify input."})
		return
	}

	// Generate ID
	id := uuid.New().String()

	// Save to MongoDB
	doc := ReceiptDocument{
		ID:           id,
		Retailer:     receipt.Retailer,
		PurchaseDate: receipt.PurchaseDate,
		PurchaseTime: receipt.PurchaseTime,
		Items:        receipt.Items,
		Total:        receipt.Total,
	}

	//Insert doc into MongoDB

	_, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save receipt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// GET

func getPoints(c *gin.Context) {
	id := c.Param("id")

// Look up receipt in Mongo

	var doc ReceiptDocument
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "No receipt found for that ID"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve receipt"})
		return
	}

	points := calculatePoints(doc)
	c.JSON(http.StatusOK, gin.H{"points": points})
}
//Validation and Business Logic

//validateReceipt checks

func validateReceipt(receipt Receipt) bool {
	retailerRegex := regexp.MustCompile(`^[\w\s\-&]+$`)
	itemDescRegex := regexp.MustCompile(`^[\w\s\-]+$`)
	priceRegex := regexp.MustCompile(`^\d+\.\d{2}$`)

	if !retailerRegex.MatchString(receipt.Retailer) {
		return false
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return false
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return false
	}

	if !priceRegex.MatchString(receipt.Total) {
		return false
	}

	for _, item := range receipt.Items {
		if !itemDescRegex.MatchString(item.ShortDescription) {
			return false
		}
		if !priceRegex.MatchString(item.Price) {
			return false
		}
	}

	return true
}

//calculatePoints computes reward points

func calculatePoints(doc ReceiptDocument) int {
	points := 0

// Rule 1: Alphanumeric characters in retailer name

	alphanumCount := 0
	for _, c := range doc.Retailer {
		if unicode.IsLetter(c) || unicode.IsDigit(c) {
			alphanumCount++
		}
	}
	points += alphanumCount

// Rule 2: Round dollar

	total, _ := strconv.ParseFloat(doc.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

// Rule 3: Multiple of 0.25

	if int(total*100)%25 == 0 {
		points += 25
	}

// Rule 4: 5 points per two items

	points += (len(doc.Items) / 2) * 5

// Rule 5: Item description length multiple of 3

	for _, item := range doc.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(price * 0.2 + 0.5) // Round up
		}
	}

// Rule 6: Odd purchase day

	purchaseDate, _ := time.Parse("2006-01-02", doc.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

// Rule 7: Purchase time between 2pm and 4pm

	purchaseTime, _ := time.Parse("15:04", doc.PurchaseTime)
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}