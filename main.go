package main

import (
	"encoding/json"
	"log"

	"github.com/Bruary/staff-scheduling/auth"
	"github.com/Bruary/staff-scheduling/core"
	"github.com/Bruary/staff-scheduling/core/models"
	"github.com/Bruary/staff-scheduling/db"
	shiftsRepo "github.com/Bruary/staff-scheduling/db/shifts"
	usersRepo "github.com/Bruary/staff-scheduling/db/users"
	_ "github.com/Bruary/staff-scheduling/docs"
	"github.com/Bruary/staff-scheduling/shifts"
	shiftsModels "github.com/Bruary/staff-scheduling/shifts/models"
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

	// New db services
	usersRepo := usersRepo.New(dbConn)
	shiftsRepo := shiftsRepo.New(dbConn)

	// services registry
	usersService := users.New(usersRepo)
	authService := auth.New(usersService)
	shiftsService := shifts.New(usersRepo, shiftsRepo)

	// controllers
	usersController := users.NewControllerService(usersService)
	authController := auth.NewControllerService(authService)
	shiftsController := shifts.NewControllerService(&shiftsService)

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
			errResp := models.JWTMissingError
			resp, _ := json.Marshal(errResp)
			c.Context().SetBody(resp)
			return nil
		}

		resp := authController.IsTokenValid(c.Context(), auth.IsTokenValidRequest{
			Token: c.Get("Authorization"),
		})
		if resp.BaseResponse != nil {
			resp, _ := json.Marshal(resp)
			c.Context().SetBody(resp)
			return nil
		}

		if resp.Valid {
			// fetch user
			user, err := usersRepo.GetUserByUid(c.Context(), resp.Claims.UserUid)
			if err != nil {
				errResp := models.UnauthorizedError
				resp, _ := json.Marshal(errResp)
				c.Context().SetBody(resp)
				return nil
			}

			// check if user have right access level
			if user.Type != string(endpointConfig.AccessLevel) {
				errResp := models.UserPermissionError
				resp, _ := json.Marshal(errResp)
				c.Context().SetBody(resp)
				return nil
			}

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

	v1.Post("/login", func(c *fiber.Ctx) error {
		req := auth.LoginRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := authController.Login(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	v1.Put("/user/permission", func(c *fiber.Ctx) error {
		req := userModels.UpdateUserPermissionLevelRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := usersController.UpdateUserPermissionLevel(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	v1.Delete("/user", func(c *fiber.Ctx) error {
		req := userModels.DeleteUserRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := usersController.DeleteUser(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	v1.Post("/shift", func(c *fiber.Ctx) error {
		req := shiftsModels.CreateShiftRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := shiftsController.CreateShift(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	v1.Delete("/shift", func(c *fiber.Ctx) error {
		req := shiftsModels.DeleteShiftRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := shiftsController.DeleteShift(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	v1.Patch("/shift", func(c *fiber.Ctx) error {
		req := shiftsModels.UpdateShiftRequest{}

		err := json.Unmarshal(c.Body(), &req)
		if err != nil {
			return err
		}

		response := shiftsController.UpdateShift(c.Context(), req)
		respBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}

		c.Context().SetBody(respBytes)
		return nil
	})

	app.Listen(":3000")
}
