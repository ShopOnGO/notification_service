package notifications

import (
	"context"
	"fmt"
	"time"

	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type NotificationRepository struct {
	collection *mongodriver.Collection
}

func NewNotificationRepository(db *mongodriver.Database) *NotificationRepository {
	return &NotificationRepository{
		collection: db.Collection("notifications"),
	}
}
func (repo *NotificationRepository) Add(n *Notification) (*Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n.ID = fmt.Sprintf("user_%d_%d", n.UserID, time.Now().UnixNano())

	fmt.Println("🧾 Генерируем ID:", n.ID)

	_, err := repo.collection.InsertOne(ctx, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}
