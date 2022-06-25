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

type orderItemPack struct {
	TableId    *string
	OrderItems []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "orderItem")

func GetOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		result, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cann't list order items"})
			return
		}

		var allOrderItems bson.M
		err = result.All(ctx, &allOrderItems)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allOrderItems)

	}
}
func GetOrderItemById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")

		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"orderItemId": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item was not found"})
			return
		}
		c.JSON(http.StatusOK, orderItem)

	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{"$match", bson.D{{"orderId", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "foodId"}, {"foreignField", "foodId"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"food", "order"}, {"localField", "orderId"}, {"foreignField", "orderId"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.TableId"}, {"foreignField", "tableId"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmtpyArrays", true}}}}

	project := bson.D{
		{
			"$project", bson.D{
				{"id", 0},
				{"amount", "$food.price"},
				{"totalCount", 1},
				{"foodName", "$foodName"},
				{"foodImage", "$food.foodImage"},
				{"tableNumber", "$table.tableNumber"},
				{"tableId", "$table.tableId"},
				{"orderId", "$order.orderId"},
				{"price", "$food.price"},
				{"quantity", 1},
			},
		},
	}

	groupStage := bson.D{{"$group", bson.D{{"id", bson.D{{"_id", bson.D{{"orderId", "$orderId"}, {"tableId", "$tableId"}, {"tableNumber", "$tableNumber"}}}, {"paymentDue", bson.D{{"$sum", "$amount"}}}, {"totalCount", bson.D{{"$sum", 1}}}, {"orderItems", bson.D{{"$push", "$$ROOT"}}}}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"paymenDue", 1},
			{"totalCount", 1},
			{"tableNumber", "$_id.tableNumber"},
			{"orderItems", 1},
		}},
	}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		project,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}
	result.All(ctx, &OrderItems)
	return OrderItems, err
}

func GetOrderItemByOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderId := c.Param("order_id")

		allOrderItem, err := ItemsByOrder(orderId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error coccurred when list all order item by order id"})
			return
		}
		c.JSON(http.StatusOK, allOrderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var _orderItemPack orderItemPack
		var order models.Order

		err := c.BindJSON(&_orderItemPack)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		orderItemToBeInserted := []interface{}{}
		order.Table_id = _orderItemPack.TableId
		order_id := OrderItemOrderCreator(order)

		for _, orderItem := range _orderItemPack.OrderItems {
			orderItem.Order_id = order_id
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": validationErr.Error()})
				return
			}
			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			var num = toFixed(*orderItem.Uint_price, 2)
			orderItem.Uint_price = &num
			orderItemToBeInserted = append(orderItemToBeInserted, orderItem)
		}

		result, err := orderItemCollection.InsertMany(ctx, orderItemToBeInserted)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, result)

	}
}

func UpdateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orderItemId := c.Param("orderItem_id")
		filter := bson.M{"order_item_id": orderItemId}

		var updateObj primitive.D

		if orderItem.Uint_price != nil {
			num := toFixed(*orderItem.Uint_price, 2)
			orderItem.Uint_price = &num
			updateObj = append(updateObj, bson.E{"unit_price", orderItem.Uint_price})
		}

		if orderItem.Quantity != "" {
			updateObj = append(updateObj, bson.E{"quantity", orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updateObj = append(updateObj, bson.E{"food_id", orderItem.Food_id})
		}

		if orderItem.Order_id != "" {
			updateObj = append(updateObj, bson.E{"order_id", orderItem.Order_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, bson.E{"updated_at", orderItem.Updated_at})

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := orderItemCollection.UpdateOne(
			ctx,
			filter,
			bson.D{{"$set", updateObj}},
			&opt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "order item update unsuccessfully"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
