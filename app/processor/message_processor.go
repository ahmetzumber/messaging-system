package processor

import (
	"context"
	"encoding/json"
	"log/slog"
	"messaging-system/app/dto"
	"messaging-system/app/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	StatusSent   = "sent"
	StatusUnsent = "unsent"

	MessageLimit = 2
)

type IMessageService interface {
	GetMessages(ctx context.Context, status string, limit int) ([]model.Message, error)
	MarkMessageAsSent(ctx context.Context, messageID primitive.ObjectID, webhookMessageID string) error
}

type IClient interface {
	SendMessage(request *dto.MessageRequest) (*dto.MessageResponse, error)
}

type ICacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type MessageProcessor struct {
	service   IMessageService
	client    IClient
	cache     ICacheService
	logger    *slog.Logger
	ticker    *time.Ticker
	isRunning bool
	stopChan  chan bool
}

func NewMessageProcessor(service IMessageService, client IClient,
	cache ICacheService, logger *slog.Logger) *MessageProcessor {
	return &MessageProcessor{
		service:  service,
		client:   client,
		cache:    cache,
		logger:   logger,
		stopChan: make(chan bool),
	}
}

func (p *MessageProcessor) Start(ctx context.Context) {
	if p.isRunning {
		p.logger.Warn("Message processor already running")
		return
	}

	p.logger.Info("Starting message processor")
	p.isRunning = true
	p.ticker = time.NewTicker(2 * time.Minute)

	go func() {
		for {
			select {
			case <-p.ticker.C:
				p.processMessages(ctx)
			case <-p.stopChan:
				p.logger.Warn("Stopping processor...")
				p.ticker.Stop()
				p.isRunning = false
				return
			}
		}
	}()
}

func (p *MessageProcessor) Stop(_ context.Context) {
	if !p.isRunning {
		p.logger.Warn("Message processor not running")
		return
	}

	p.stopChan <- true
}

func (p *MessageProcessor) GetSentMessages(ctx context.Context, limit int) ([]model.Message, error) {
	return p.service.GetMessages(ctx, StatusSent, limit)
}

func (p *MessageProcessor) processMessages(ctx context.Context) {
	messages, err := p.service.GetMessages(ctx, StatusUnsent, MessageLimit)
	if err != nil {
		p.logger.Error("failed to fetch messages", slog.Any("error", err))
		return
	}

	if len(messages) == 0 {
		p.logger.Info("no unsent messages found")
		return
	}

	for _, message := range messages {
		request := message.ConvertToRequest()
		resp, err := p.client.SendMessage(request)
		if err != nil {
			p.logger.Error("failed to send message",
				slog.String("messageId", message.ID.Hex()),
				slog.Any("error", err),
			)
			continue
		}

		p.logger.Info("message sent",
			"messageId", message.ID,
			"webhookMessageId", resp.MessageID,
		)

		if markErr := p.service.MarkMessageAsSent(ctx, message.ID, resp.MessageID); markErr != nil {
			p.logger.Error("failed to mark message as sent",
				"messageId", message.ID,
				"error", markErr,
			)
			continue
		}

		cacheValue := &model.CacheMessage{
			MessageID: resp.MessageID,
			SentAt:    time.Now().UTC().Format(time.RFC3339),
		}

		jsonValue, err := json.Marshal(cacheValue)
		if err != nil {
			p.logger.Error("failed to marshal cache value", "error", err)
			continue
		}

		if err := p.cache.Set(ctx, message.ID.Hex(), jsonValue, 24*time.Hour); err != nil {
			p.logger.Error("failed to cache message", "messageId", message.ID, "error", err)
			continue
		}
	}
}
