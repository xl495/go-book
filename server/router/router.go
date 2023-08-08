package router

import (
	"api-fiber-gorm/handler"
	"api-fiber-gorm/middleware"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// Middleware
	api := app.Group("/api", logger.New())
	api.Get("/", handler.Hello)

	// Auth
	auth := api.Group("/auth")
	auth.Post("/login", handler.Login)

	// User
	user := api.Group("/user")
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Patch("/:id", middleware.Protected(), handler.UpdateUser)
	user.Delete("/:id", middleware.Protected(), handler.DeleteUser)

	// Product
	product := api.Group("/product")
	product.Get("/", handler.GetAllProducts)
	product.Get("/:id", handler.GetProduct)
	product.Post("/", middleware.Protected(), handler.CreateProduct)
	product.Delete("/:id", middleware.Protected(), handler.DeleteProduct)

	//	 Category
	category := api.Group("/category", middleware.Protected())

	category.Get("/", handler.GetCategoryList)
	category.Get("/:id", handler.GetCategory)
	category.Post("/", handler.CreateCategory)
	category.Patch("/:id", handler.UpdateCategory)
	category.Delete("/:id", handler.DeleteCategory)

	//	Book
	book := api.Group("/book", func(ctx *fiber.Ctx) error {
		fmt.Println("middleware")
		return ctx.Next()
	}, middleware.Protected())
	book.Get("/", handler.GetBookList)
	book.Get("/:id", handler.GetBook)
	book.Post("/", handler.CreateBook)
	book.Patch("/:id", handler.UpdateBook)
	book.Delete("/:id", handler.DeleteBook)

	// 404
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"data":    nil,
			"message": "没有找到该路由",
		})
	})
}
