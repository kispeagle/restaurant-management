package controller

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	database "restaurant-management/databases"
	"restaurant-management/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "food")
var menuCollection *mongo.Collection = database.OpenCollection(database.MongoDB, "menu")
var validate = validator.New()

func GetFoods() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil && recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err := strconv.Atoi(c.Query("Page"))
		if err != nil && page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", "1"}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"_id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
				},
			},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage,
			groupStage,
			projectStage,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred when listing food items"})
			return
		}

		var allFoods []primitive.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allFoods)
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
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
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
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	out := math.Pow(10, float64(precision))
	return float64(round(num*out)) / out
}

func UpdateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		idFood := c.Param("food_id")
		var food models.Food
		var menu models.Menu
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateObj primitive.D
		if food.Name != nil {
			updateObj = append(updateObj, bson.E{"name", food.Name})
		}

		if food.Price != nil {
			updateObj = append(updateObj, bson.E{"price", food.Price})
		}

		if food.FoodImage != nil {
			updateObj = append(updateObj, bson.E{"food_image", food.FoodImage})
		}

		if food.MenuId != nil {
			err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.MenuId}).Decode(&menu)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Menu was not found"})
				return
			}
			updateObj = append(updateObj, bson.E{"menu_id", food.MenuId})
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		updateObj = append(updateObj, primitive.E{"updated_at", food.Updated_at})
		upsert := true
		filter := bson.M{"food_id": idFood}
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}
		result, err := foodCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "food item update unsuccessfully"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
