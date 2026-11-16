package storage

import (
	"context"
	"cook_book/internal/config"
	"cook_book/internal/domain"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func NewStorage(cfg *config.Config) (*DataBase, error) {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%d/?authSource=admin",
		cfg.Mongo.User,
		cfg.Mongo.Password,
		cfg.Mongo.Host,
		cfg.Mongo.Port,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, fmt.Errorf("NewStorage: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("NewStorage: %w", err)
	}

	return &DataBase{
		DB:     mongoClient.Database(cfg.Mongo.Name),
		Client: mongoClient,
	}, nil
}

// типа миграция для создания нулевого админа
func MigrateAdminUser(db *mongo.Database, cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := db.Collection("users")

	count, err := usersCollection.CountDocuments(ctx, bson.M{"role": "Admin"})
	if err != nil {
		return fmt.Errorf("MigrateAdminUser: %w", err)
	}

	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("MigrateAdminUser: %w", err)
		}

		admin := domain.User{
			Name:      cfg.Admin.Name,
			Email:     cfg.Admin.Email,
			HashPass:  string(hashedPassword),
			Role:      "Admin",
			CreatedAt: time.Now(),
		}

		if _, err = usersCollection.InsertOne(ctx, admin); err != nil {
			return fmt.Errorf("MigrateAdminUser: %w", err)
		}
	}

	return nil
}
