package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id"`
	InvoiceId        string             `json:"invoice_id"`
	OrderId          string             `json:"order_id"`
	Payment_method   *string            `json:Payment_method" validate:"eq=CARD|eq=CASH|eq="`
	Payment_status   *string            `json:"payment_status" validate:"required,eq=PENDING|eq=DONE"`
	Payment_due_date time.Time          `json:"payment_due_date`
	Created_at       time.Time          `json:Created_at`
	Updated_at       time.Time          `json:Updated_at`
}
