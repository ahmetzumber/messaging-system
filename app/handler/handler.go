package handler

import (
	"context"
	"messaging-system/app/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	ActionStart = "start"
	ActionStop  = "stop"
)

type IMessageProcessor interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
	GetSentMessages(ctx context.Context, limit int) ([]model.Message, error)
}

type Handler struct {
	processor IMessageProcessor
}

func NewMessageHandler(processor IMessageProcessor) *Handler {
	return &Handler{processor: processor}
}

func (h *Handler) RegisterRoutes(server *fiber.App) {
	processor := server.Group("/processor")
	processor.Get("/sent-messages", h.GetSentMessages)
	processor.Post("/:action", h.StartStopJob)
}

func (h *Handler) GetSentMessages(c *fiber.Ctx) error {
	ctx := c.Context()
	limit := c.Query("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid limit parameter",
		})
	}

	messages, err := h.processor.GetSentMessages(ctx, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if len(messages) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no unsent messages found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

func (h *Handler) StartStopJob(c *fiber.Ctx) error {
	ctx := c.Context()
	action := c.Params("action")
	switch action {
	case ActionStart:
		h.processor.Start(ctx)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "message processor started",
		})
	case ActionStop:
		h.processor.Stop(ctx)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "message processor stopped",
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid action, use 'start' or 'stop'",
		})
	}
}
