package processor

import (
	"context"
	"log/slog"
	"messaging-system/app/dto"
	"messaging-system/app/mocks"
	"messaging-system/app/model"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestMessageProcessor_Start(t *testing.T) {
	mockService, mockClient, mockCache, logger := createMockServices(t)
	processor := NewMessageProcessor(mockService, mockClient, mockCache, logger)

	processor.Start(context.Background())

	processor.ticker.Stop()

	assert.True(t, processor.isRunning)
	assert.NotNil(t, processor.ticker)
}

func TestMessageProcessor_Stop(t *testing.T) {
	mockService, mockClient, mockCache, logger := createMockServices(t)
	processor := NewMessageProcessor(mockService, mockClient, mockCache, logger)

	ctx := context.Background()
	processor.Start(ctx)

	go func() {
		time.Sleep(10 * time.Millisecond)
		processor.Stop(ctx)
	}()

	time.Sleep(20 * time.Millisecond)
	assert.False(t, processor.isRunning)
}

func TestMessageProcessor_GetSentMessages(t *testing.T) {
	mockService, mockClient, mockCache, logger := createMockServices(t)
	processor := NewMessageProcessor(mockService, mockClient, mockCache, logger)

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusSent, 10).
			Return([]model.Message{
				{Content: "Message 1", Status: StatusSent},
			}, nil)

		messages, err := processor.GetSentMessages(ctx, 10)
		assert.NoError(t, err)
		assert.NotNil(t, messages)
	})

	t.Run("error", func(t *testing.T) {
		ctx := context.Background()
		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusSent, 10).
			Return(nil, assert.AnError)

		messages, err := processor.GetSentMessages(ctx, 10)
		assert.Nil(t, messages)
		assert.Error(t, err)
	})
}

func TestMessageProcessor_processMessages(t *testing.T) {
	mockService, mockClient, mockCache, logger := createMockServices(t)
	processor := NewMessageProcessor(mockService, mockClient, mockCache, logger)

	ctx := context.Background()
	t.Run("fetch messages error", func(t *testing.T) {
		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return(nil, assert.AnError)

		processor.processMessages(ctx)
	})

	t.Run("no messages to process", func(t *testing.T) {
		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return([]model.Message{}, nil)

		processor.processMessages(ctx)
	})

	t.Run("process messages successfully", func(t *testing.T) {
		message := model.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+90555",
			Content:     "Hello",
			Status:      StatusUnsent,
		}

		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return([]model.Message{message}, nil)

		mockClient.
			EXPECT().
			SendMessage(gomock.Any()).
			Return(&dto.MessageResponse{Message: "accepted", MessageID: "webhook123"}, nil)

		mockService.
			EXPECT().
			MarkMessageAsSent(gomock.Any(), message.ID, "webhook123").
			Return(nil)

		mockCache.
			EXPECT().
			Set(gomock.Any(), message.ID.Hex(), gomock.Any()).
			Return(nil)

		processor.processMessages(ctx)
	})

	t.Run("send message error", func(t *testing.T) {
		message := model.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+90555",
			Content:     "Hello",
			Status:      StatusUnsent,
		}

		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return([]model.Message{message}, nil)

		mockClient.
			EXPECT().
			SendMessage(gomock.Any()).
			Return(nil, assert.AnError)

		processor.processMessages(ctx)
	})

	t.Run("mark message as sent error", func(t *testing.T) {
		message := model.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+90555",
			Content:     "Hello",
			Status:      StatusUnsent,
		}

		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return([]model.Message{message}, nil)

		mockClient.
			EXPECT().
			SendMessage(gomock.Any()).
			Return(&dto.MessageResponse{Message: "accepted", MessageID: "webhook123"}, nil)

		mockService.
			EXPECT().
			MarkMessageAsSent(gomock.Any(), message.ID, "webhook123").
			Return(assert.AnError)

		processor.processMessages(ctx)
	})

	t.Run("cache set error", func(t *testing.T) {
		message := model.Message{
			ID:          primitive.NewObjectID(),
			PhoneNumber: "+90555",
			Content:     "Hello",
			Status:      StatusUnsent,
		}

		mockService.
			EXPECT().
			GetMessages(gomock.Any(), StatusUnsent, MessageLimit).
			Return([]model.Message{message}, nil)

		mockClient.
			EXPECT().
			SendMessage(gomock.Any()).
			Return(&dto.MessageResponse{Message: "accepted", MessageID: "webhook123"}, nil)

		mockService.
			EXPECT().
			MarkMessageAsSent(gomock.Any(), message.ID, "webhook123").
			Return(nil)

		mockCache.
			EXPECT().
			Set(gomock.Any(), message.ID.Hex(), gomock.Any()).
			Return(assert.AnError)

		processor.processMessages(ctx)
	})
}

func createMockServices(t *testing.T) (*mocks.MockIMessageService, *mocks.MockIClient,
	*mocks.MockICacheService, *slog.Logger) {
	t.Helper()
	mockController := gomock.NewController(t)

	mockService := mocks.NewMockIMessageService(mockController)
	mockClient := mocks.NewMockIClient(mockController)
	mockCache := mocks.NewMockICacheService(mockController)
	logger := slog.Default()

	return mockService, mockClient, mockCache, logger
}
