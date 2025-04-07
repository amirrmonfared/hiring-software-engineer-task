package handler

import (
	"strconv"
	"sweng-task/internal/service"

	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AdHandler struct {
	adService *service.AdService
	log       *zap.SugaredLogger
}

func NewAdHandler(adService *service.AdService, log *zap.SugaredLogger) *AdHandler {
	return &AdHandler{
		adService: adService,
		log:       log,
	}
}

func (h *AdHandler) GetWinningAds(c *fiber.Ctx) error {
	placement := c.Query("placement")
	if placement == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "placement query param is required",
		})
	}

	category := c.Query("category")
	keyword := c.Query("keyword")

	limitStr := c.Query("limit", "1")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "limit must be an integer between 1 and 10",
		})
	}

	ads, err := h.adService.GetWinningAds(placement, category, keyword, limit)
	if err != nil {
		h.log.Errorw("Failed to get winning ads", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Failed to get winning ads",
		})
	}

	return c.JSON(ads)
}
