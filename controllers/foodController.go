package controller

import (
	"context"
	"fmt"
	"net/http"
	database "restaurant-management/databases"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func GetFoodById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("food_id")

		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var food models.Food
		err := foodCollection.FindOne(context, bson.M{"food_id": id}).Decode(&food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occur while fetching food item"})
		}
		c.JSON(http.StatusOK, food)
	}
}

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		context, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var menu models.Menu
		var food models.Food

		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := validate.Struct(food)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = menuCollection.FindOne(context, bson.M{"menu_id": food.MenuId}).Decode(&menu)
		if err != nil {
			msg := fmt.Sprintf("menu was not found")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		food.CreatedAt = time.Now()
		food.CreatedAt.Format(time.RFC3339)
		food.UpdatedAt = time.Now()
		food.UpdatedAt.Format(time.RFC3339)
		food.ID = primitive.NewObjectID()
		food.FoodId = food.ID.Hex()
		var fixedPrice = toFixed(*food.Price, 2)
		food.Price = &fixedPrice

		result, insertErr := foodCollection.InsertOne(context, food)
		if insertErr != nil {
			msg := fmt.Sprintf("Food item was not created")
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func round(num float64) int {
	return 0
}

func toFixed(num float64, precision int) float64 {
	return 0.0
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
