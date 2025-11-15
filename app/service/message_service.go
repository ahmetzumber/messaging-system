package service

import (
	"context"
	"messaging-system/app/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IRepository interface {
	GetMessages(ctx context.Context, status string, limit int) ([]model.Message, error)
	MarkMessageAsSent(ctx context.Context, messageID primitive.ObjectID, webhookMessageID string) error
}

type MessageService struct {
	repo IRepository
}

func NewMessageService(repo IRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) GetMessages(ctx context.Context, status string, limit int) ([]model.Message, error) {
	messages, err := s.repo.GetMessages(ctx, status, limit)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *MessageService) MarkMessageAsSent(ctx context.Context, messageID primitive.ObjectID, webhookMessageID string) error {
	return s.repo.MarkMessageAsSent(ctx, messageID, webhookMessageID)
}
