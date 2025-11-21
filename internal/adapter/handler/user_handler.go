package handler

import (
	"ifulblog/internal/adapter/handler/request"
	"ifulblog/internal/adapter/handler/response"
	"ifulblog/internal/core/domain/entity"
	"ifulblog/internal/core/service"
	validatorLib "ifulblog/lib/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type UserHandler interface {
	GetUserByID(c *fiber.Ctx) error
	UpdatePassword(c *fiber.Ctx) error
}

type userHandler struct {
	userService service.UserService
}

// GetUserByID implements UserHandler.
func (u *userHandler) GetUserByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetUserByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	user, err := u.userService.GetUserByID(c.Context(), int64(userID))
	if err != nil {
		code = "[HANDLER] GetUserByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Success"

	resp := response.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	defaultSuccesResponse.Data = resp

	return c.JSON(defaultSuccesResponse)
}

// UpdatePassword implements UserHandler.
func (u *userHandler) UpdatePassword(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] UpdatePassword - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.UpdatePasswordRequest
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] UpdatePassword - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[HANDLER] UpdatePassword - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	err := u.userService.UpdatePassword(c.Context(), req.NewPassword, int64(userID))
	if err != nil {
		code = "[HANDLER] UpdatePassword - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}
	
	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Success"
	defaultSuccesResponse.Data = nil

	return c.JSON(defaultSuccesResponse)

}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{
		userService: userService,
	}
}
