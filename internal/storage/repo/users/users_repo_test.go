package users

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"cook_book/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testRepo *userRepo
var testClient *mongo.Client

func TestMain(m *testing.M) {
	ctx := context.Background()

	mongoURI := "mongodb://admin:123@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	testClient = client
	testDB := client.Database("cook_book_test")
	testCollection := testDB.Collection("users_test")

	testRepo = &userRepo{
		collection: testCollection,
		logger:     slog.Default(),
	}

	code := m.Run()

	testCollection.Drop(ctx)
	client.Disconnect(ctx)

	os.Exit(code)
}

func TestUserRepo_Create(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:     "UserTest",
		Email:    "test1@example.com",
		HashPass: "password123",
	}

	err := testRepo.Create(ctx, user)

	assert.NoError(t, err)
	assert.NotEmpty(t, user.Id)
	assert.False(t, user.CreatedAt.IsZero())
}

func TestUserRepo_FindByID(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:     "UserTest",
		Email:    "test2@example.com",
		HashPass: "password123",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	found, err := testRepo.FindByID(ctx, user.Id)

	assert.NoError(t, err)
	assert.Equal(t, user.Name, found.Name)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepo_FindByEmail(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:     "UserTest",
		Email:    "test3@example.com",
		HashPass: "password123",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	found, err := testRepo.FindByEmail(ctx, user.Email)

	assert.NoError(t, err)
	assert.Equal(t, user.Id, found.Id)
	assert.Equal(t, user.Email, found.Email)
}

func TestUserRepo_Update(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:     "TestUser",
		Email:    "test4@example.com",
		HashPass: "pass123",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	newName := "UpdatedName"
	newEmail := "updatedemail@example.com"
	err = testRepo.Update(ctx, user.Id, &newName, &newEmail, nil)
	assert.NoError(t, err)

	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newEmail, updated.Email)
	assert.Equal(t, user.HashPass, updated.HashPass)
}

func TestUserRepo_UpdateRole(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:  "TestUser",
		Email: "test5@example.com",
		Role:  "user",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	err = testRepo.UpdateRole(ctx, "Admin", user.Id)
	assert.NoError(t, err)

	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	assert.Equal(t, domain.RoleAdmin, updated.Role)
}

func TestUserRepo_UpdateRefreshToken(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:  "UserTest",
		Email: "test6@example.com",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	token := "refresh_token_123"
	err = testRepo.UpdateRefreshToken(ctx, user.Id, token)
	assert.NoError(t, err)

	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	assert.Equal(t, token, updated.RefreshToken)
	assert.NotNil(t, updated.RefreshTokenExpiresAt)
}

func TestUserRepo_Delete(t *testing.T) {
	ctx := context.Background()

	user := &domain.User{
		Name:  "UserTest",
		Email: "test7@example.com",
	}

	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	err = testRepo.Delete(ctx, user.Id)
	assert.NoError(t, err)

	_, err = testRepo.FindByID(ctx, user.Id)
	assert.Error(t, err)
}

func TestUserRepo_Update_NotFound(t *testing.T) {
	ctx := context.Background()

	fakeID := primitive.NewObjectID()
	newName := "UserName"

	err := testRepo.Update(ctx, fakeID, &newName, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUserRepo_CleanExpiredRefreshTokens(t *testing.T) {
	ctx := context.Background()

	expiredTime := time.Now().Add(-24 * time.Hour)
	user := &domain.User{
		Name:                  "UserTest",
		Email:                 "test8@example.com",
		RefreshToken:          "expired_token",
		RefreshTokenExpiresAt: &expiredTime,
	}
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	err = testRepo.CleanExpiredRefreshTokens(ctx)
	assert.NoError(t, err)

	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	assert.Empty(t, updated.RefreshToken)
	assert.Nil(t, updated.RefreshTokenExpiresAt)
}
