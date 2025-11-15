package model

import (
	"messaging-system/app/dto"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID               primitive.ObjectID `json:"id" bson:"_id"`
	WebhookMessageID string             `json:"webhookMessageId" bson:"webhookMessageId,omitempty"`
	PhoneNumber      string             `json:"phoneNumber" bson:"phoneNumber"`
	Content          string             `json:"content" bson:"content"`
	Status           string             `json:"status" bson:"status"`
	SentAt           time.Time          `bson:"sentAt" json:"sentAt"`
}

type CacheMessage struct {
	MessageID string `json:"messageId"`
	SentAt    string `json:"sentAt"`
}

func (m *Message) ConvertToRequest() *dto.MessageRequest {
	return &dto.MessageRequest{
		To:      m.PhoneNumber,
		Content: m.Content,
	}
}
