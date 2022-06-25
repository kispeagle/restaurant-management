package main

import (
	"fmt"
	"os"

	"restaurant-management/logger"
	"restaurant-management/middleware"
	"restaurant-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.GetCustomProductionLogger()
	logger.Logger.Info("===== START SERVER =====")

	PORT := os.Getenv("PORT")
	fmt.Println(PORT)

	if PORT == "" {
		PORT = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	routes.FoodRoutes(router)
	routes.MenuRoutes(router)
	routes.InvoiceRoutes(router)
	routes.OrderItemRoutes(router)
	routes.OrderRoutes(router)
	routes.TableRoutes(router)

	router.Run("localhost:" + PORT)

}
