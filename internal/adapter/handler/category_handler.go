package handler

import (
	"ifulblog/internal/adapter/handler/request"
	"ifulblog/internal/adapter/handler/response"
	"ifulblog/internal/core/domain/entity"
	"ifulblog/internal/core/service"
	"ifulblog/lib/conv"
	validatorLib "ifulblog/lib/validator"

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
	var req request.CategoryRequest
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] CreateCategory - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	if err := c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateCategory - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err := validatorLib.ValidateStruct(req); err != nil {
		code = "[HANDLER] CreateCategory - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	reqEntity := entity.CategoryEntity{
		Title: req.Title,
		User: entity.UserEntity{
			ID: int64(userID),
		},
	}

	err := ch.categoryService.CreateCategory(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateCategory - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Data = nil
	defaultSuccesResponse.Pagination = nil
	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Category created successfully"

	return c.JSON(defaultSuccesResponse)
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
		code = "[HANDLER] GetCategories - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponses := []response.SuccessCategoryResponse{}
	for _, result := range results {
		categoryResponse := response.SuccessCategoryResponse{
			ID:            result.ID,
			Title:         result.Title,
			Slug:          result.Slug,
			CreatedByName: result.User.Name,
		}
		categoryResponses = append(categoryResponses, categoryResponse)
	}
	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Pagination = nil
	defaultSuccesResponse.Meta.Message = "Categories fetched successfully"
	defaultSuccesResponse.Data = categoryResponses

	return c.JSON(defaultSuccesResponse)
}

// GetCategoryByID implements CategoryHandler.
func (ch *categoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	userID := claims.UserID
	if userID == 0 {
		code = "[HANDLER] GetCategoryByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthorized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("categoryID")
	id, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		code = "[HANDLER] GetCategoryByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	categoryResponse := response.SuccessCategoryResponse{
		ID: id,
		Title: result.Title,
		Slug: result.Slug,
		CreatedByName: result.User.Name,
	}
	
	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Pagination = nil
	defaultSuccesResponse.Meta.Message = "Category fetched detail successfully"
	defaultSuccesResponse.Data = categoryResponse

	return c.JSON(defaultSuccesResponse)

}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandler{
		categoryService: categoryService,
	}
}
