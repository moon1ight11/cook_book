package users

import (
	"context"
	"cook_book/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

// содание пользователя
func (u *userRepo) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()

	_, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error in userCreate: %w", err)
	}

	return nil
}

// поиск пользователя по id
func (u *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User

	filter := bson.M{"_id": id}

	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error in userFindById: %w", err)
	}

	return &user, err
}

// поиск пользоателя по емайл
func (u *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	filter := bson.M{"email": email}

	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("error in userFindByEmail: %w", err)
	}

	return &user, err
}

// обновление пользователя
func (u *userRepo) Update(
	ctx context.Context,
	id primitive.ObjectID,
	name *string,
	email *string,
	pass *string,
) error {
	update := bson.M{}

	if name != nil {
		update["name"] = *name
	}
	if email != nil {
		update["email"] = *email
	}
	if pass != nil {
		update["hashpass"] = *pass
	}

	if len(update) == 0 {
		return fmt.Errorf("error in userUpdate: %s", "nothing to update")
	}

	update["updatedat"] = time.Now()

	filter := bson.M{"_id": id}

	result, err := u.collection.UpdateOne(ctx, filter, bson.M{
		"$set": update,
	})
	if err != nil {
		return fmt.Errorf("error in userUpdate: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("error in userUpdate: user with id=%s not found", id)
	}

	return nil
}

// обновление рефрештокена
func (u *userRepo) UpdateRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) error {
	_, err := u.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"refresh_token": refreshToken,
				"updated_at":    time.Now(),
			},
		},
	)
	return err
}

// удаление пользователя
func (u *userRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := u.collection.DeleteOne(ctx, filter)

	return err
}
