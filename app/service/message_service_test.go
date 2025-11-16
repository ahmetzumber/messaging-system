package service

import (
	"context"
	"messaging-system/app/mocks"
	"messaging-system/app/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMessageService_GetMessages(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockIRepository(mockController)
	messageService := NewMessageService(mockRepo)
	t.Run("retrieve messages successfully", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetMessages(gomock.Any(), "unsent", 10).
			Return([]model.Message{}, nil)

		messages, err := messageService.GetMessages(ctx, "unsent", 10)
		assert.Nil(t, err)
		assert.NotNil(t, messages)
	})

	t.Run("error retrieving messages", func(t *testing.T) {
		mockRepo.
			EXPECT().
			GetMessages(gomock.Any(), "unsent", 10).
			Return(nil, assert.AnError)

		messages, err := messageService.GetMessages(ctx, "unsent", 10)
		assert.NotNil(t, err)
		assert.Nil(t, messages)
	})
}

func TestMessageService_MarkMessageAsSent(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	ctx := context.Background()
	mockRepo := mocks.NewMockIRepository(mockController)
	messageService := NewMessageService(mockRepo)
	t.Run("mark message as sent successfully", func(t *testing.T) {
		mockRepo.
			EXPECT().
			MarkMessageAsSent(gomock.Any(), gomock.Any(), "webhook123").
			Return(nil)

		err := messageService.MarkMessageAsSent(ctx, primitive.NewObjectID(), "webhook123")
		assert.Nil(t, err)
	})

	t.Run("error marking message as sent", func(t *testing.T) {
		mockRepo.
			EXPECT().
			MarkMessageAsSent(gomock.Any(), gomock.Any(), "webhook123").
			Return(assert.AnError)

		err := messageService.MarkMessageAsSent(ctx, primitive.NewObjectID(), "webhook123")
		assert.NotNil(t, err)
	})
}
