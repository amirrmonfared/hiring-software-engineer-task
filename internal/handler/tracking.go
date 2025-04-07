package handler

import (
	"sweng-task/internal/model"
	"sweng-task/internal/service"

	"github.com/go-playground/validator/v10"
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type TrackingHandler struct {
	trackingService *service.TrackingService
	lineItemService *service.LineItemService
	log             *zap.SugaredLogger
	validate        *validator.Validate
}

func NewTrackingHandler(ts *service.TrackingService, ls *service.LineItemService, log *zap.SugaredLogger) *TrackingHandler {
	return &TrackingHandler{
		trackingService: ts,
		lineItemService: ls,
		log:             log,
		validate:        validator.New(),
	}
}

func (h *TrackingHandler) TrackEvent(c *fiber.Ctx) error {
	var event model.TrackingEvent
	if err := c.BodyParser(&event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := h.validate.Struct(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Validation failed",
			"details": err.Error(),
		})
	}

	if _, err := h.lineItemService.GetByID(event.LineItemID); err == service.ErrLineItemNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Line item not found",
		})
	}

	if err := h.trackingService.RecordEvent(event); err != nil {
		h.log.Errorw("Failed to record tracking event", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to record tracking event",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
	})
}
