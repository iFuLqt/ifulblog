package handler

import (
	"ifulblog/internal/adapter/handler/response"
	"ifulblog/internal/core/domain/entity"
	"ifulblog/internal/core/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type ContentHandler interface {
	GetContents(c *fiber.Ctx) error
	GetContentByID(c *fiber.Ctx) error
	CreateContent(c *fiber.Ctx) error
	UpdateContent(c *fiber.Ctx) error
	DeleteContent(c *fiber.Ctx) error
	UploadImageR2(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// CreateContent implements ContentHandler.
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContentByID implements ContentHandler.
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetContents implements ContentHandler.
func (ch *contentHandler) GetContents(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] GetContents - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	results, err := ch.contentService.GetContents(c.Context())
	if err != nil {
		code = "[HANDLER] GetContens - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Success"

	responseContents := []response.ContentResponse{}
	for _, content := range results {
		respContent := response.ContentResponse{
			ID:           content.ID,
			Title:        content.Title,
			Excerpt:      content.Excerpt,
			Description:  content.Description,
			Image:        content.Image,
			Tags:         content.Tags,
			Status:       content.Status,
			CategoryID:   content.CategoryID,
			CreatedByID:  content.CreatedByID,
			CreatedAt:    content.CreatedAt.Format(time.RFC3339),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}
		
		responseContents = append(responseContents, respContent)
	}

	defaultSuccesResponse.Data = responseContents
	return c.JSON(defaultSuccesResponse)
}

// UpdateContent implements ContentHandler.
func (ch *contentHandler) UpdateContent(c *fiber.Ctx) error {
	panic("unimplemented")
}

// UploadImageR2 implements ContentHandler.
func (ch *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
