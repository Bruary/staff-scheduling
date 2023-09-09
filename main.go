package main

import (
	"encoding/json"
	"log"

	"github.com/Bruary/staff-scheduling/auth"
	"github.com/Bruary/staff-scheduling/core"
	"github.com/Bruary/staff-scheduling/db"
	_ "github.com/Bruary/staff-scheduling/docs"
	"github.com/Bruary/staff-scheduling/models"
	"github.com/Bruary/staff-scheduling/users"
	userModels "github.com/Bruary/staff-scheduling/users/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load config file", err)
	}

	// Establish DB connection
	dbConn := db.EstablishDBConnection()

	// New db service
	database := db.New(dbConn)

	// services registry
	usersService := users.New(database)
	authService := auth.New(usersService)

	// controllers
	usersController := users.NewControllerService(usersService)
	authController := auth.NewControllerService(authService)

	app := fiber.New()

	// documentaion open api
	app.Get("/swagger/*", swagger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Use(func(c *fiber.Ctx) error {

		c.Context().SetContentType("application/json")

		endpointConfig, ok := core.Endpoints[core.Endpoint{Path: string(c.Request().URI().Path()), Method: string(c.Request().Header.Method())}]
		if !ok {

			errResp := models.BaseResponse{
				ErrorType: "Endpoint does not exist",
				ErrorMsg:  "Endpoint was not found",
			}
			resp, _ := json.Marshal(errResp)
			c.Context().SetBody(resp)

			return nil
		}

		if !endpointConfig.RequireJWT {

			c.Next()
			return nil
		}

		if c.Get("Authorization") == "" {
			errResp := models.BaseResponse{
				ErrorType: "Token Missing",
				ErrorMsg:  "Token is missing",
			}
			resp, _ := json.Marshal(errResp)
			c.Context().SetBody(resp)
			return nil
		}

		resp := authController.IsTokenValid(c.Context(), auth.IsTokenValidRequest{
			Token: &auth.Token{
				Token:     c.Get("Authorization"),
				TokenType: auth.JWT,
			},
		})
		if resp.BaseResponse != nil {
			resp, _ := json.Marshal(resp)
			c.Context().SetBody(resp)
			return nil
		}

		if resp.Valid {
			c.Next()
			return nil
		} else {
			c.Context().SetBody([]byte(fiber.ErrUnauthorized.Error()))
			return nil
		}

	})

	v1.Post("/signup", func(c *fiber.Ctx) error {
		req := userModels.CreateUserRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := usersController.CreateUser(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	app.Listen(":3000")
}
