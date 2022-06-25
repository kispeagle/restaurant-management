package controller

import (
	"context"
	"log"
	"net/http"
	database "restaurant-management/databases"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type invoiceView struct {
	Invoice_id      string
	Payment_method  string
	OrderId         string
	Payment_status  string
	Payment_due     interface{}
	Table_number    interface{}
	Paymen_due_date time.Time
	Order_details   interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "invoice")

func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		result, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred when listing all invoices"})
			return
		}

		var allInvoices []bson.M
		err = result.All(ctx, &allInvoices)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoiceById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceId := c.Param("invoice_id")
		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_dd": invoiceId}).Decode(&invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var _invoiceView invoiceView

		allOrderItem, err := ItemsByOrder(invoice.OrderId)
		_invoiceView.Invoice_id = invoice.InvoiceId
		_invoiceView.Paymen_due_date = invoice.Payment_due_date
		_invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			_invoiceView.Payment_method = *invoice.Payment_method
		}
		_invoiceView.Invoice_id = invoice.InvoiceId
		_invoiceView.Payment_status = *invoice.Payment_status
		_invoiceView.Payment_due = allOrderItem[0]["payment_due"]
		_invoiceView.Table_number = allOrderItem[0]["table_number"]
		_invoiceView.Order_details = allOrderItem[0]["order_items"]

		c.JSON(http.StatusOK, _invoiceView)
	}
}

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var invoice models.Invoice

		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var order models.Order
		err = orderCollection.FindOne(ctx, bson.M{"orderId": invoice.OrderId}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order id was not found"})
			return
		}

		validationErr := validate.Struct(invoice)
		if validationErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceId = invoice.ID.Hex()

		rs, insertError := invoiceCollection.InsertOne(ctx, invoice)
		if insertError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cann't insert order"})
			return
		}
		c.JSON(http.StatusOK, rs)
	}
}

func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		invoiceId := c.Param("invoice_id")

		var invoice models.Invoice
		err := c.BindJSON(&invoice)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"invoice_id": invoiceId}
		var updateObj primitive.D
		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{"payment_method", invoice.Payment_method})
		}

		if invoice.Payment_status != nil {
			updateObj = append(updateObj, bson.E{"payment_status", invoice.Payment_status})
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", invoice.Updated_at})
		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}

		result, err := invoiceCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update invoice unsuccessfully"})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}
