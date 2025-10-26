package main

import (
	"context"
	"strconv"
	"time"

	db "go-backend-task/db/sqlc"
	intdb "go-backend-task/internal/db"
	"go-backend-task/internal/logger"
	"go-backend-task/internal/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

var validate = validator.New()

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func sendErrorResponse(c *fiber.Ctx, status int, errType, msg string) error {
	resp := ErrorResponse{
		Error:   errType,
		Message: msg,
	}
	return c.Status(status).JSON(resp)
}

func main() {
	logger.Init()
	defer logger.Sync()

	intdb.Connect()
	app := fiber.New()
	app.Use(middleware.RequestIDMiddleware())
	app.Use(middleware.LoggingMiddleware())

	app.Get("/", func(c *fiber.Ctx) error {
		logger.Logger.Info("GET / called")
		return c.SendString("Hello from GoFiber!")
	})

	app.Get("/users", func(c *fiber.Ctx) error {
		pageStr := c.Query("page", "1")
		limitStr := c.Query("limit", "10")

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit

		queries := db.New(intdb.DB)
		users, err := queries.ListUsersWithPagination(
			context.Background(),
			db.ListUsersWithPaginationParams{
				Limit:  int32(limit),
				Offset: int32(offset),
			},
		)
		if err != nil {
			logger.Logger.Error("Failed to list users with pagination", zap.Error(err))
			return sendErrorResponse(c, 500, "ServerError", err.Error())
		}

		return c.JSON(users)
	})

	app.Post("/users", func(c *fiber.Ctx) error {
		logger.Logger.Info("POST /users called")

		type CreateUserInput struct {
			Name string `json:"name" validate:"required,min=2,max=100"`
			DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
		}

		var input CreateUserInput
		if err := c.BodyParser(&input); err != nil {
			logger.Logger.Error("Invalid input to create user", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Invalid input: "+err.Error())
		}

		if err := validate.Struct(input); err != nil {
			logger.Logger.Error("Validation failed on create user", zap.Error(err))
			return sendErrorResponse(c, 400, "ValidationError", err.Error())
		}

		t, _ := time.Parse("2006-01-02", input.DOB)

		var dob pgtype.Date
		if err := dob.Scan(t); err != nil {
			logger.Logger.Error("Date scan failed", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Date scan failed: "+err.Error())
		}

		queries := db.New(intdb.DB)
		user, err := queries.CreateUser(context.Background(), db.CreateUserParams{
			Name: input.Name,
			Dob:  dob,
		})
		if err != nil {
			logger.Logger.Error("Failed to create user", zap.Error(err))
			return sendErrorResponse(c, 500, "ServerError", "Failed to create user: "+err.Error())
		}

		logger.Logger.Info("User created", zap.Int32("id", user.ID))
		return c.Status(201).JSON(user)
	})

	app.Get("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		logger.Logger.Info("GET /users/:id called", zap.String("id", idParam))

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Logger.Error("Invalid user ID", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Invalid user ID")
		}

		queries := db.New(intdb.DB)
		user, err := queries.GetUserById(context.Background(), int32(id))
		if err != nil {
			logger.Logger.Error("User not found", zap.Error(err))
			return sendErrorResponse(c, 404, "NotFound", "User not found")
		}

		if !user.Dob.Valid {
			logger.Logger.Error("User has null DOB")
			return sendErrorResponse(c, 500, "ServerError", "User has null DOB")
		}
		dobTime := user.Dob.Time

		now := time.Now()
		age := now.Year() - dobTime.Year()
		if now.YearDay() < dobTime.YearDay() {
			age--
		}

		type UserResponse struct {
			ID   int32  `json:"id"`
			Name string `json:"name"`
			DOB  string `json:"dob"`
			Age  int    `json:"age"`
		}
		resp := UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  dobTime.Format("2006-01-02"),
			Age:  age,
		}

		logger.Logger.Info("User fetched", zap.Int32("id", user.ID))
		return c.JSON(resp)
	})

	app.Put("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		logger.Logger.Info("PUT /users/:id called", zap.String("id", idParam))

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Logger.Error("Invalid user ID", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Invalid user ID")
		}

		type UpdateUserInput struct {
			Name string `json:"name" validate:"required,min=2,max=100"`
			DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
		}

		var input UpdateUserInput
		if err := c.BodyParser(&input); err != nil {
			logger.Logger.Error("Invalid input to update user", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Invalid input: "+err.Error())
		}

		if err := validate.Struct(input); err != nil {
			logger.Logger.Error("Validation failed on update user", zap.Error(err))
			return sendErrorResponse(c, 400, "ValidationError", err.Error())
		}

		t, _ := time.Parse("2006-01-02", input.DOB)

		var dob pgtype.Date
		if err := dob.Scan(t); err != nil {
			logger.Logger.Error("Date scan failed", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Date scan failed: "+err.Error())
		}

		queries := db.New(intdb.DB)
		user, err := queries.UpdateUser(context.Background(), db.UpdateUserParams{
			ID:   int32(id),
			Name: input.Name,
			Dob:  dob,
		})
		if err != nil {
			logger.Logger.Error("Failed to update user", zap.Error(err))
			return sendErrorResponse(c, 500, "ServerError", "Failed to update user: "+err.Error())
		}

		logger.Logger.Info("User updated", zap.Int32("id", user.ID))
		return c.JSON(user)
	})

	app.Delete("/users/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		logger.Logger.Info("DELETE /users/:id called", zap.String("id", idParam))

		id, err := strconv.Atoi(idParam)
		if err != nil {
			logger.Logger.Error("Invalid user ID", zap.Error(err))
			return sendErrorResponse(c, 400, "BadRequest", "Invalid user ID")
		}

		queries := db.New(intdb.DB)
		err = queries.DeleteUser(context.Background(), int32(id))
		if err != nil {
			logger.Logger.Error("Failed to delete user", zap.Error(err))
			return sendErrorResponse(c, 500, "ServerError", "Failed to delete user: "+err.Error())
		}

		logger.Logger.Info("User deleted", zap.Int("id", id))
		return c.SendStatus(204)
	})

	app.Listen(":3000")
}
