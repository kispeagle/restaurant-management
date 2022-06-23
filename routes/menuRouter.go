package routes

import (
	controller "restaurant-management/controllers"

	"github.com/gin-gonic/gin"
)

func MenuRoutes(incomingRoute *gin.Engine) {
	incomingRoute.GET("/menus", controller.GetMenus())
	incomingRoute.GET("/menus/:menu_id", controller.GetMenuById())
	incomingRoute.POST("/menus/", controller.CreateMenu())
	incomingRoute.PATCH("/menus/:menu_id", controller.UpdateMenu())
}
