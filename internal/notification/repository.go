package notifications

import (
	"context"
	"time"

	"notification/internal/model"
	"notification/pkg/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepository struct {
	collection *mongodriver.Collection
}

func NewNotificationRepository(db *mongo.Mongo) *NotificationRepository {
	return &NotificationRepository{
		collection: db.Database.Collection("notifications"),
	}
}

func (repo *NotificationRepository) Add(n *model.Notification) (*model.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := repo.collection.InsertOne(ctx, n, options.InsertOne())
	if err != nil {
		return nil, err
	}

	// Присваиваем сгенерированный ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		n.ID = oid
	}

	return n, nil
}
