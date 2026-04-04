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
	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error in NewStorage: %w", err)
	}

	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("error in NewStorage: %w", err)
	}

	db := mongoClient.Database(cfg.Mongo.Name)

	Database := DataBase{
		DB:     db,
		Client: mongoClient,
	}

	return &Database, nil
}

// типа миграция для создания нулевого админа
func MigrateAdminUser(db *mongo.Database, cfg *config.Config) error {
	// создаем базовый контекст и выбираем коллекцию
	ctx := context.Background()
	usersCollection := db.Collection("users")

	// выбираем записи с ролью админ
	count, err := usersCollection.CountDocuments(ctx, bson.M{"role": "admin"})
	if err != nil {
		return fmt.Errorf("error in MigrateAdmin: %w", err)
	}

	// если админа нет — создаём
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), bcrypt.DefaultCost)

		admin := domain.User{
			Name:      cfg.Admin.Name,
			Email:     cfg.Admin.Email,
			HashPass:  string(hashedPassword),
			Role:      "Admin",
			CreatedAt: time.Now(),
		}

		_, err := usersCollection.InsertOne(ctx, admin)
		if err != nil {
			return fmt.Errorf("error in MigrateAdmin: %w", err)
		}
	}

	return nil
}
