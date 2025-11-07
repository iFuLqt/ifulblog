package handler

import (
	"ifulblog/internal/adapter/handler/response"
	"ifulblog/internal/core/domain/entity"
	"ifulblog/internal/core/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)


var defaultSuccesResponse response.DefaultSuccessResponse


type CategoryHandler interface {
	GetCategories(c *fiber.Ctx) error
	GetCategoryByID(c *fiber.Ctx) error 
	CreateCategory(c *fiber.Ctx) error
	EditCategoryByID(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
}

type categoryHandler struct {
	categoryService service.CategoryService
}

// CreateCategory implements CategoryHandler.
func (ch *categoryHandler) CreateCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteCategory implements CategoryHandler.
func (ch *categoryHandler) DeleteCategory(c *fiber.Ctx) error {
	panic("unimplemented")
}

// EditCategoryByID implements CategoryHandler.
func (ch *categoryHandler) EditCategoryByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetCategories implements CategoryHandler.
func (ch *categoryHandler) GetCategories(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID

	if userID == 0 {
		code = "[HANDLER] GetCategories - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"
		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.categoryService.GetCategories(c.Context())
	if err != nil {
		code = "[HANDLER] GetCategories"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID: result.ID,
			Title: result.Title,
			Slug: result.Slug,
			CreatedByName: result.User.Name,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}
	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Categories fetched successfully"
	defaultSuccesResponse.Data = categoryResponses

	return c.JSON(defaultSuccesResponse)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}
