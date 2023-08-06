package main

import (
	"api-fiber-gorm/database"
	"api-fiber-gorm/router"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	validate := validator.New()

	database.ConnectDB()

	app.Get("/api/test", func(ctx *fiber.Ctx) error {
		type User struct {
			ID        uint   `validate:"required,omitempty"`
			Firstname string `validate:"required,min=8,max=32"`
			Password  string `validate:"gte=10"` // gte = Greater than or equal
		}

		user := User{
			ID:        1,
			Firstname: "Fiber1",
			/*
				if you delete Firstname field
				you'll get response like this: Error:Field validation for 'Firstname' failed on the 'required' tag"
			*/
			Password: "FiberPassword123",
			/*
				if you enter "Fiber" in Password
				you'll get response like this: Error:Field validation for 'Password' failed on the 'gte' tag"+

			*/
		}

		if err := validate.Struct(user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
		}

		return ctx.Status(fiber.StatusOK).JSON("success time")
	})

	router.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
