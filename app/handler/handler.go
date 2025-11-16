package handler

import (
	"context"
	"messaging-system/app/dto"
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

// GetSentMessages godoc
// @Summary Get sent messages
// @Description Retrieves a list of sent messages with an optional limit
// @Tags processor
// @Accept json
// @Produce json
// @Param limit query int false "Number of messages to retrieve" default(10)
// @Success 200 {array} model.Message "List of sent messages"
// @Failure 400 {object} dto.ErrorResponse "Invalid limit parameter"
// @Failure 404 {object} dto.ErrorResponse "No unsent messages found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /processor/sent-messages [get]
func (h *Handler) GetSentMessages(c *fiber.Ctx) error {
	ctx := c.Context()
	limit := c.Query("limit", "10")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid limit parameter",
		})
	}

	messages, err := h.processor.GetSentMessages(ctx, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.ErrorResponse{
			Error: err.Error(),
		})
	}

	if len(messages) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(dto.ErrorResponse{
			Error: "no unsent messages found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}

// StartStopJob godoc
// @Summary Start or stop message processor
// @Description Starts or stops the message processor job based on the action parameter
// @Tags processor
// @Accept json
// @Produce json
// @Param action path string true "Action to perform" Enums(start, stop)
// @Success 200 {object} dto.SuccessResponse "Message processor started or stopped successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid action parameter"
// @Router /processor/{action} [post]
func (h *Handler) StartStopJob(c *fiber.Ctx) error {
	ctx := c.Context()
	action := c.Params("action")
	switch action {
	case ActionStart:
		h.processor.Start(ctx)
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Message: "message processor started",
		})
	case ActionStop:
		h.processor.Stop(ctx)
		return c.Status(fiber.StatusOK).JSON(dto.SuccessResponse{
			Message: "message processor stopped",
		})
	default:
		return c.Status(fiber.StatusBadRequest).JSON(dto.ErrorResponse{
			Error: "invalid action, use 'start' or 'stop'",
		})
	}
}
