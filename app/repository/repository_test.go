package repository

import (
	"context"
	"messaging-system/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	mongoImage     = "mongo:7.0.4"
	mockDB         = "message"
	mockCollection = "messages"
)

func createTestContainer(ctx context.Context) (repo *Repository, clean func()) {
	mongodbContainer, err := mongodb.Run(ctx, mongoImage)
	if err != nil {
		panic(err)
	}

	cleanFunc := func() {
		if cleanErr := mongodbContainer.Terminate(ctx); cleanErr != nil {
			panic(err)
		}
	}

	uri, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	repo, err = New(ctx, &config.Mongo{
		URI:               uri,
		Database:          mockDB,
		MessageCollection: mockCollection,
	})
	if err != nil {
		panic(err)
	}

	return repo, cleanFunc
}

func TestRepository_GetMessages(t *testing.T) {
	ctx := context.Background()
	repo, clean := createTestContainer(ctx)
	defer clean()

	t.Run("", func(t *testing.T) {
		messages, err := repo.GetMessages(ctx, "unsent", 10)
		assert.NoError(t, err)
		assert.Len(t, messages, 0)
	})

	t.Run("returns only unsent messages with limit", func(t *testing.T) {
		testData := []interface{}{
			bson.M{"_id": primitive.NewObjectID(), "status": "unsent", "to": "+9053", "content": "A"},
			bson.M{"_id": primitive.NewObjectID(), "status": "unsent", "to": "+9053", "content": "B"},
			bson.M{"_id": primitive.NewObjectID(), "status": "sent", "to": "+9053", "content": "C"},
		}

		_, err := repo.messageCollection.InsertMany(ctx, testData)
		assert.NoError(t, err)

		messages, err := repo.GetMessages(ctx, "unsent", 1)
		assert.NoError(t, err)

		assert.Len(t, messages, 1)
		assert.Equal(t, "unsent", messages[0].Status)
	})
}

func TestRepository_MarkMessageAsSent(t *testing.T) {
	ctx := context.Background()
	repo, clean := createTestContainer(ctx)
	defer clean()

	t.Run("successfully updates message status and webhookMessageId", func(t *testing.T) {
		id := primitive.NewObjectID()
		_, err := repo.messageCollection.InsertOne(ctx, bson.M{
			"_id":     id,
			"status":  "unsent",
			"to":      "+9053",
			"content": "Hello",
		})
		assert.NoError(t, err)

		webhookID := "webhook-123"
		err = repo.MarkMessageAsSent(ctx, id, webhookID)
		assert.NoError(t, err)

		var result bson.M
		err = repo.messageCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
		assert.NoError(t, err)

		assert.Equal(t, "sent", result["status"])
		assert.Equal(t, webhookID, result["webhookMessageId"])
	})
}
