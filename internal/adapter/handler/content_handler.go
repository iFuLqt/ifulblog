package handler

import (
	"fmt"
	"ifulblog/internal/adapter/handler/request"
	"ifulblog/internal/adapter/handler/response"
	"ifulblog/internal/core/domain/entity"
	"ifulblog/internal/core/service"
	"ifulblog/lib/conv"
	validatorLib "ifulblog/lib/validator"
	"os"
	"strings"
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

	//FE
	GetContentWithQuery(c *fiber.Ctx) error
	ContentDetail(c *fiber.Ctx) error
}

type contentHandler struct {
	contentService service.ContentService
}

// ContentDetail implements ContentHandler.
func (ch *contentHandler) ContentDetail(c *fiber.Ctx) error {
	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] ContentDetail - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), contentID)
	if err != nil {
		code = "[HANDLER] ContentDetail - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Success"

	respContent := response.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccesResponse.Data = respContent
	return c.JSON(defaultSuccesResponse)
}

// GetContentWithQuery implements ContentHandler.
func (ch *contentHandler) GetContentWithQuery(c *fiber.Ctx) error {
	page := 1
	if c.Query("page") != "" {
		page, err = conv.StringToInt(c.Query("page"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 1"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid page number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	limit := 6
	if c.Query("limit") != "" {
		limit, err = conv.StringToInt(c.Query("limit"))
		if err != nil {
			code = "[HANDLER] GetContentWithQuery - 2"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = "Invalid list number"

			return c.Status(fiber.StatusBadRequest).JSON(errorResp)
		}
	}

	orderBy := "created_at"
	if c.Query("orderBy") != "" {
		orderBy = c.Query("orderBy")
	}

	orderType := "desc"
	if c.Query("orderType") != "" {
		orderType = c.Query("orderType")
	}

	search := ""
	if c.Query("search") != "" {
		search = c.Query("search")
	}

	reqEntity := entity.QueryString{
		Limit:     int64(limit),
		Page:      int64(page),
		OrderBy:   orderBy,
		OrderType: orderType,
		Search:    search,
		Status: "PUBLISH",
	}

	results,totalData, totalPages, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContentWithQuery - 3"
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
			CreatedAt:    content.CreatedAt.Local().Format("02 January 2006"),
			CategoryName: content.Category.Title,
			Author:       content.User.Name,
		}

		responseContents = append(responseContents, respContent)
	}

	defaultSuccesResponse.Data = responseContents
	defaultSuccesResponse.Pagination = &response.PaginationResponse{
		TotalRecords: totalData,
		Page: int64(page),
		PerPage: int64(limit),
		TotalPages: totalPages,
	}

	return c.JSON(defaultSuccesResponse)
}

// CreateContent implements ContentHandler.
func (ch *contentHandler) CreateContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] CreateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	userID := claims.UserID
	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] CreateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[HANDLER] CreateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err := ch.contentService.CreateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Content created successfully"
	defaultSuccesResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccesResponse)
}

// DeleteContent implements ContentHandler.
func (ch *contentHandler) DeleteContent(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] DeleteContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] DeleteContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	err = ch.contentService.DeleteContent(c.Context(), contentID)
	if err != nil {
		code = "[HANDLER] DeleteContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Delete Content Success"
	defaultSuccesResponse.Data = nil

	return c.JSON(defaultSuccesResponse)
}

// GetContentByID implements ContentHandler.
func (ch *contentHandler) GetContentByID(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] GetContentByID - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] GetContentByID - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	result, err := ch.contentService.GetContentByID(c.Context(), contentID)
	if err != nil {
		code = "[HANDLER] GetContentByID - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Success"

	respContent := response.ContentResponse{
		ID:           result.ID,
		Title:        result.Title,
		Excerpt:      result.Excerpt,
		Description:  result.Description,
		Image:        result.Image,
		Tags:         result.Tags,
		Status:       result.Status,
		CategoryID:   result.CategoryID,
		CreatedByID:  result.CreatedByID,
		CreatedAt:    result.CreatedAt.Format(time.RFC3339),
		CategoryName: result.Category.Title,
		Author:       result.User.Name,
	}

	defaultSuccesResponse.Data = respContent
	return c.JSON(defaultSuccesResponse)

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

	reqEntity := entity.QueryString{
		Limit:      10,
		Page:       1,
		OrderBy:    "",
		OrderType:  "",
		Search:     "",
		CategoryID: 0,
	}

	results, err := ch.contentService.GetContents(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] GetContents - 2"
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
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] UpdateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	userID := claims.UserID
	var req request.ContentRequest
	if err = c.BodyParser(&req); err != nil {
		code = "[HANDLER] UpdateContent - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Invalid request body"

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = validatorLib.ValidateStruct(&req); err != nil {
		code = "[HANDLER] UpdateContent - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	idParam := c.Params("contentID")
	contentID, err := conv.StringToInt64(idParam)
	if err != nil {
		code = "[HANDLER] UpdateContent 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	tags := strings.Split(req.Tags, ",")
	reqEntity := entity.ContentEntity{
		ID:          contentID,
		Title:       req.Title,
		Excerpt:     req.Excerpt,
		Description: req.Description,
		Image:       req.Image,
		Tags:        tags,
		Status:      req.Status,
		CategoryID:  req.CategoryID,
		CreatedByID: int64(userID),
	}

	err = ch.contentService.UpdateContent(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] CreateContent - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Content updated successfully"
	defaultSuccesResponse.Data = nil

	return c.Status(fiber.StatusCreated).JSON(defaultSuccesResponse)
}

// UploadImageR2 implements ContentHandler.
func (ch *contentHandler) UploadImageR2(c *fiber.Ctx) error {
	claims := c.Locals("user").(*entity.JwtData)
	if claims.UserID == 0 {
		code = "[HANDLER] UpdateContent - 1"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = "Unauthrized access"

		return c.Status(fiber.StatusUnauthorized).JSON(errorResp)
	}

	var req request.FileUploadRequest
	file, err := c.FormFile("image")
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 2"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusBadRequest).JSON(errorResp)
	}

	if err = c.SaveFile(file, fmt.Sprintf("./temp/content/%s", file.Filename)); err != nil {
		code = "[HANDLER] UploadImageR2 - 3"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	req.Image = fmt.Sprintf("./temp/content/%s", file.Filename)
	reqEntity := entity.FileUploadEntity{
		Name: fmt.Sprintf("%d-%d", int64(claims.UserID), time.Now().UnixNano()),
		Path: req.Image,
	}

	imageUrl, err := ch.contentService.UploadImageR2(c.Context(), reqEntity)
	if err != nil {
		code = "[HANDLER] UploadImageR2 - 4"
		log.Errorw(code, err)
		errorResp.Meta.Status = false
		errorResp.Meta.Message = err.Error()

		return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
	}

	if imageUrl != "" {
		err = os.Remove(req.Image)
		if err != nil {
			code = "[HANDLER] UploadImageR2 - 5"
			log.Errorw(code, err)
			errorResp.Meta.Status = false
			errorResp.Meta.Message = err.Error()

			return c.Status(fiber.StatusInternalServerError).JSON(errorResp)
		}
	}

	urlImageResp := map[string]interface{}{
		"urlImage": imageUrl,
	}

	defaultSuccesResponse.Meta.Status = true
	defaultSuccesResponse.Meta.Message = "Content updated successfully"
	defaultSuccesResponse.Data = urlImageResp

	return c.Status(fiber.StatusCreated).JSON(defaultSuccesResponse)
}

func NewContentHandler(contentService service.ContentService) ContentHandler {
	return &contentHandler{
		contentService: contentService,
	}
}
