package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/seus31/todo-application/backend/dto/requests"
	"github.com/seus31/todo-application/backend/dto/responses"
	"github.com/seus31/todo-application/backend/models"
	"github.com/seus31/todo-application/backend/services"
	"github.com/seus31/todo-application/backend/utils"
	"time"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	var req requests.CreateUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Request parsing failed"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashPassword,
	}

	if err := uc.UserService.CreateUser(utils.GetContextFromFiber(ctx), user); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (uc *UserController) GetUsers(ctx *fiber.Ctx) error {
	var req requests.GetUsersRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid parameters"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	offset := (req.Page - 1) * req.Limit
	users, err := uc.UserService.GetUsers(utils.GetContextFromFiber(ctx), req.Limit, offset)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get users"})
	}

	return ctx.Status(fiber.StatusOK).JSON(users)
}

func (uc *UserController) GetUser(ctx *fiber.Ctx) error {
	var req requests.GetUserRequest
	if err := ctx.ParamsParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid parameters"})
	}

	if err := utils.ValidateStruct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := uc.UserService.GetUser(utils.GetContextFromFiber(ctx), req.ID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	response := responses.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
