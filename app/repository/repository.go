package repository

import (
	"context"
	"messaging-system/app/model"
	"messaging-system/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	client            *mongo.Client
	database          *mongo.Database
	messageCollection *mongo.Collection
}

func New(ctx context.Context, conf *config.Mongo) (*Repository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.URI))
	if err != nil {
		return nil, err
	}

	return &Repository{
		client:            client,
		database:          client.Database(conf.Database),
		messageCollection: client.Database(conf.Database).Collection(conf.MessageCollection),
	}, nil
}

func (r *Repository) GetMessages(ctx context.Context, status string, limit int) ([]model.Message, error) {
	var messages []model.Message

	filter := bson.M{
		"status": status,
	}

	opts := options.Find()
	opts.SetLimit(int64(limit))

	result, err := r.messageCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := result.All(ctx, &messages); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *Repository) MarkMessageAsSent(ctx context.Context, messageID primitive.ObjectID, webhookMessageID string) error {
	filter := bson.M{"_id": messageID}
	update := bson.M{
		"$set": bson.M{
			"status":           "sent",
			"webhookMessageId": webhookMessageID,
		},
	}
	_, err := r.messageCollection.UpdateOne(ctx, filter, update)
	return err
}
