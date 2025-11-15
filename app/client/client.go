package client

import (
	"log/slog"
	"messaging-system/app/dto"
	"messaging-system/config"

	"github.com/go-resty/resty/v2"
)

const (
	HeaderAuth = "x-ins-auth-key"
)

type Client struct {
	resty  *resty.Client
	conf   *config.Client
	logger *slog.Logger
}

func NewClient(conf *config.Client, logger *slog.Logger) *Client {
	restyClient := resty.
		New().
		SetBaseURL(conf.URL).
		SetHeader(HeaderAuth, conf.ApiKey)

	logger.Info("Client initialized", "URL", conf.URL)
	return &Client{
		resty:  restyClient,
		conf:   conf,
		logger: logger,
	}
}

func (c *Client) SendMessage(request *dto.MessageRequest) (*dto.MessageResponse, error) {
	if validErr := request.Validate(); validErr != nil {
		c.logger.Error("Invalid message request", "error", validErr)
		return nil, validErr
	}

	response := &dto.MessageResponse{}
	res, err := c.resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		SetResult(response).
		Post("/fa36ba41-1794-4295-8053-d198499d2fc9")

	if err != nil {
		c.logger.Error("Error sending message", "error", err)
		return nil, err
	}

	c.logger.Info("Message successfully sent", "StatusCode", res.StatusCode(), "Response", response)
	return response, nil
}
