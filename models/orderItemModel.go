package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `bson:"_id"`
	Quantity    string             `json:"quantity" validate:"required, eq=S|eq=M|eq=L"`
	OrderItemId string             `json:"orderItemId" validate:"required"`
	UintPrice   *float64           `json:"unitPrice" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
	FoodId      *string            `json:"foodId" validate:"required"`
	OrderId     string             `json:"orderId" validate:"required"`
}
