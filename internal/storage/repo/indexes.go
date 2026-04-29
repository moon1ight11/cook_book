package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// индексы для бд
type IndexManager struct {
	recipesCollection *mongo.Collection
	usersCollection   *mongo.Collection
}

// конструктор
func NewIndexManager(recipesCollection, usersCollection *mongo.Collection) *IndexManager {
	return &IndexManager{
		recipesCollection: recipesCollection,
		usersCollection:   usersCollection,
	}
}

// основная функция создания индексов
func (im *IndexManager) CreateAllIndexes(ctx context.Context) error {
	if err := im.createRecipeIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create recipe indexes: %w", err)
	}

	if err := im.createUserIndexes(ctx); err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	return nil
}

// индексы для рецептов
func (im *IndexManager) createRecipeIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "owner_id", Value: 1}},
			Options: options.Index().
				SetName("idx_owner_id"),
		},
		{
			Keys: bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().
				SetName("idx_tags"),
		},
		{
			Keys: bson.D{
				{Key: "created_at", Value: -1},
			},
			Options: options.Index().
				SetName("idx_created_at_desc"),
		},
	}

	_, err := im.recipesCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create recipe indexes: %w", err)
	}

	return nil
}

// индексы для юзеров
func (im *IndexManager) createUserIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().
				SetName("idx_email_unique").
				SetUnique(true),
		},
	}

	_, err := im.usersCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return fmt.Errorf("failed to create user indexes: %w", err)
	}

	return nil
}
