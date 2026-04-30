package users

import (
	"context"
	"cook_book/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// содание пользователя
func (u *userRepo) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	id := primitive.NewObjectID()
	user.Id = id

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
		return nil, fmt.Errorf("error in userFindById: %w", err)
	}

	return &user, nil
}

// поиск пользоателя по емайл
func (u *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	filter := bson.M{"email": email}

	err := u.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("error in userFindByEmail: %w", err)
	}

	return &user, nil
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
		update["password"] = *pass
	}

	if len(update) == 0 {
		return fmt.Errorf("error in userUpdate: %s", "nothing to update")
	}

	update["updated_at"] = time.Now()

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

// обновление роли пользователя
func (u *userRepo) UpdateRole(ctx context.Context, newRole string, id primitive.ObjectID) error {
	// ставим фильтр
	filter := bson.M{"_id": id}

	// выбираем что обновлять
	update := bson.M{"role": newRole}

	// обновляем
	result, err := u.collection.UpdateOne(ctx, filter, bson.M{
		"$set": update,
	})
	if err != nil {
		return fmt.Errorf("error in userUpdateRole: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("error in userUpdateRole: user with id=%s not found", id)
	}

	return nil
}

// обновление рефрештокена
func (u *userRepo) UpdateRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) error {
	filter := bson.M{"_id": id}

	var expiresAt *time.Time
	if refreshToken != "" {
		expiry := time.Now().Add(7 * 24 * time.Hour)
		expiresAt = &expiry
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_token":            refreshToken,
			"refresh_token_expires_at": expiresAt,
			"updated_at":               time.Now(),
		},
	}

	result, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error in userUpdateRefreshToken: %w", err)
	}

	if result.MatchedCount == 0 {
		u.logger.Warn("UpdateRefreshToken: user not found", "id", id.Hex())
	}

	return nil
}

// удаление пользователя
func (u *userRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error in userDelete: %w", err)
	}

	return nil
}

// удаление пользователя для админа
func (u *userRepo) DeleteAny(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}

	_, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error in userDeleteAny: %w", err)
	}

	return nil
}

// очистка просроченных токенов
func (u *userRepo) CleanExpiredRefreshTokens(ctx context.Context) error {
	// Находим всех пользователей с истекшими токенами
	filter := bson.M{
		"refresh_token_expires_at": bson.M{
			"$lt": time.Now(),
		},
		"refresh_token": bson.M{
			"$ne": "",
		},
	}

	update := bson.M{
		"$set": bson.M{
			"refresh_token":            "",
			"refresh_token_expires_at": nil,
			"updated_at":               time.Now(),
		},
	}

	result, err := u.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error in CleanExpiredRefreshTokens: %w", err)
	}

	if result.ModifiedCount > 0 {
		u.logger.Info("Cleaned expired refresh tokens", "count", result.ModifiedCount)
	}

	return nil
}
