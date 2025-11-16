package handler

import (
	"messaging-system/app/mocks"
	"messaging-system/app/model"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

const limit = 10

func TestHandler_GetSentMessages(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockProcessor := mocks.NewMockIMessageProcessor(mockController)
	mockHandler := NewMessageHandler(mockProcessor)

	app := fiber.New()
	app.Get("/processor/sent-messages", mockHandler.GetSentMessages)

	t.Run("invalid limit query param", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/processor/sent-messages?limit=abc", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("processor returns error", func(t *testing.T) {
		mockProcessor.
			EXPECT().
			GetSentMessages(gomock.Any(), limit).
			Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/processor/sent-messages?limit=10", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("no sent messages found", func(t *testing.T) {
		mockProcessor.
			EXPECT().
			GetSentMessages(gomock.Any(), limit).
			Return([]model.Message{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/processor/sent-messages?limit=10", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("successful retrieval of sent messages", func(t *testing.T) {
		mockProcessor.
			EXPECT().
			GetSentMessages(gomock.Any(), limit).
			Return([]model.Message{
				{Content: "Hello"},
				{Content: "World"},
			}, nil)

		req := httptest.NewRequest(http.MethodGet, "/processor/sent-messages?limit=10", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestHandler_StartStopJob(t *testing.T) {
	mockController := gomock.NewController(t)
	defer mockController.Finish()

	mockProcessor := mocks.NewMockIMessageProcessor(mockController)
	mockHandler := NewMessageHandler(mockProcessor)

	app := fiber.New()
	app.Post("/processor/:action", mockHandler.StartStopJob)

	t.Run("successfully start action", func(t *testing.T) {
		mockProcessor.
			EXPECT().
			Start(gomock.Any())

		req := httptest.NewRequest(http.MethodPost, "/processor/start", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successfully stop action", func(t *testing.T) {
		mockProcessor.
			EXPECT().
			Stop(gomock.Any())

		req := httptest.NewRequest(http.MethodPost, "/processor/stop", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("invalid action name - 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/processor/invalidaction", nil)
		resp, err := app.Test(req)

		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
