package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"restaurant-management/logger"
	"restaurant-management/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		result, err := menuCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred when listing menu"})
			return
		}
		var allMenus []primitive.M
		if err = result.All(ctx, &allMenus); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allMenus)
	}
}

func GetMenuById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		menuId := c.Param("menu_id")
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": menuId}).Decode(&menu)
		if err != nil {
			msg := fmt.Sprintf("Error occured while fetching menu id: %d", menuId)
			c.JSON(http.StatusInternalServerError, msg)
			return
		}
		c.JSON(http.StatusOK, menu)
	}
}

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		var menu models.Menu
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		errVal := validate.Struct(menu)
		if errVal != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menu.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		menu.ID = primitive.NewObjectID()
		menu.Menu_id = menu.ID.Hex()

		result, insertError := menuCollection.InsertOne(ctx, menu)
		if insertError != nil {
			c.JSON(http.StatusInternalServerError, "Menu cann't created")
			return
		}
		c.JSON(http.StatusOK, result)

	}
}

func inTimeSpan(start, end, now time.Time) bool {
	return start.Before(now) && end.After(now)
}

func UpdateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var menu models.Menu
		logger.Logger.Infof("%#v", menu)
		err := c.BindJSON(&menu)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		menu_id := c.Param("menu_id")
		filter := bson.M{"menu_id": menu_id}
		var updateObj primitive.D

		// logger.Logger.Infof("%s - %#v - %#v", menu.Start_date.Format(time.RFC3339), menu.End_date, time.Now().Format(time.RFC3339))
		if menu.Start_date != nil && menu.End_date != nil {
			if !inTimeSpan(*menu.Start_date, *menu.End_date, time.Now()) {
				msg := "kindly retype the time"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				defer cancel()
				return
			}
			updateObj = append(updateObj, bson.E{"start_date", menu.Start_date})
			updateObj = append(updateObj, bson.E{"end_date", menu.End_date})
			if menu.Name != "" {
				updateObj = append(updateObj, bson.E{"name", menu.Name})
			}
			if menu.Category != "" {
				updateObj = append(updateObj, bson.E{"category", menu.Category})
			}

			menu.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			updateObj = append(updateObj, bson.E{"updated_at", menu.Updated_at})
			upsert := true

			opt := options.UpdateOptions{
				Upsert: &upsert,
			}

			result, err := menuCollection.UpdateOne(ctx, filter, bson.D{{"$set", updateObj}}, &opt)
			if err != nil {
				msg := "Menu update failed"
				c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
				return
			}
			defer cancel()

			c.JSON(http.StatusOK, result)
		}

	}
}
