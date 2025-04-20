package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/yandiriswandi/pos-sanbercode-final-project/controllers"
	"github.com/yandiriswandi/pos-sanbercode-final-project/middlewares"
	"github.com/yandiriswandi/pos-sanbercode-final-project/utils/swagger"
)

func StartSever() *gin.Engine {
	router := gin.Default()

	swagger.Initiator(router)

	//crud category
	router.POST("/category", middlewares.Authorize(1), controllers.CreateCategory)
	router.PUT("/category/:id", middlewares.Authorize(1), controllers.UpdateCategory)
	router.DELETE("/category/:id", middlewares.Authorize(1), controllers.DeleteCategory)
	router.GET("/category", middlewares.Authorize(), controllers.GetCategoryList)
	router.GET("/category/:id", middlewares.Authorize(), controllers.GetCategoryByID)

	//auth
	router.POST("/login", controllers.Login)

	router.POST("/upload", controllers.UploadFile)

	//crud product
	router.POST("/product", middlewares.Authorize(1), controllers.CreateProduct)
	router.PUT("/product/:id", middlewares.Authorize(1), controllers.UpdateProduct)
	router.DELETE("/product/:id", middlewares.Authorize(1), controllers.DeleteProduct)
	router.GET("/product", controllers.GetProducts)
	router.GET("/product/:id", controllers.GetProduct)

	//crud user
	router.POST("/user", middlewares.Authorize(1), controllers.CreateUser)
	router.PUT("/user/:id", middlewares.Authorize(1), controllers.UpdateUser)
	router.DELETE("/user/:id", middlewares.Authorize(1), controllers.DeleteUser)
	router.GET("/user", middlewares.Authorize(1), controllers.GetUsers)
	router.GET("/user/:id", middlewares.Authorize(), controllers.GetUser)

	//crud cart
	router.POST("/cart", middlewares.Authorize(), controllers.CreateCart)
	router.PUT("/cart/:id", middlewares.Authorize(), controllers.UpdateCart)
	router.DELETE("/cart/:id", middlewares.Authorize(), controllers.DeleteCart)
	router.GET("/cart", middlewares.Authorize(), controllers.GetCarts)
	router.GET("/cart/:id", middlewares.Authorize(), controllers.GetCart)

	return router
}
