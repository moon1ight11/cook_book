package users

import (
	"context"
	"fmt"
	"time"

	"cook_book/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (u *userRepo) Create(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.Id = primitive.NewObjectID()

	if _, err := u.collection.InsertOne(ctx, user); err != nil {
		return fmt.Errorf("userCreate: %w", err)
	}

	return nil
}

func (u *userRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User

	if err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, fmt.Errorf("userFindByID: %w", err)
	}

	return &user, nil
}

func (u *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	if err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, fmt.Errorf("userFindByEmail: %w", err)
	}

	return &user, nil
}

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
		return fmt.Errorf("userUpdate: nothing to update")
	}

	update["updated_at"] = time.Now()

	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("userUpdate: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("userUpdate: user with id=%s not found", id.Hex())
	}

	return nil
}

func (u *userRepo) UpdateRole(ctx context.Context, newRole string, id primitive.ObjectID) error {
	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": newRole}})
	if err != nil {
		return fmt.Errorf("userUpdateRole: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("userUpdateRole: user with id=%s not found", id.Hex())
	}

	return nil
}

func (u *userRepo) UpdateRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) error {
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

	result, err := u.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("userUpdateRefreshToken: %w", err)
	}

	if result.MatchedCount == 0 {
		u.logger.Warn("UpdateRefreshToken: user not found", "id", id.Hex())
	}

	return nil
}

func (u *userRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	if _, err := u.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return fmt.Errorf("userDelete: %w", err)
	}

	return nil
}

func (u *userRepo) DeleteAny(ctx context.Context, id primitive.ObjectID) error {
	if _, err := u.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return fmt.Errorf("userDeleteAny: %w", err)
	}

	return nil
}

func (u *userRepo) CleanExpiredRefreshTokens(ctx context.Context) error {
	filter := bson.M{
		"refresh_token_expires_at": bson.M{"$lt": time.Now()},
		"refresh_token":            bson.M{"$ne": ""},
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
		return fmt.Errorf("CleanExpiredRefreshTokens: %w", err)
	}

	if result.ModifiedCount > 0 {
		u.logger.Info("Cleaned expired refresh tokens", "count", result.ModifiedCount)
	}

	return nil
}
