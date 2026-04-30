package users

import (
	"context"
	"cook_book/internal/domain"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var testRepo *userRepo
var testClient *mongo.Client

// точка входа для тестов
func TestMain(m *testing.M) {
	ctx := context.Background()

	// соединяемся с тестовой бд
	mongoURI := "mongodb://admin:123@localhost:27017"
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}

	// создаем тестовые бд и коллекцию пользователей
	testClient = client
	testDB := client.Database("cook_book_test")
	testCollection := testDB.Collection("users_test")

	// создаем репозиторий
	testRepo = &userRepo{
		collection: testCollection,
		logger:     slog.Default(),
	}

	// запускаем тесты
	code := m.Run()

	// после тестов все удаляем
	testCollection.Drop(ctx)
	client.Disconnect(ctx)

	// выходим из тестов
	os.Exit(code)
}

// тест на создание пользователя
func TestUserRepo_Create(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:     "UserTest",
		Email:    "test1@example.com",
		HashPass: "password123",
	}

	// добавляем в репо
	err := testRepo.Create(ctx, user)

	// проверка выходных данных
	assert.NoError(t, err)
	assert.NotEmpty(t, user.Id)
	assert.False(t, user.CreatedAt.IsZero())
}

// тест на поиск пользователя по айди
func TestUserRepo_FindByID(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:     "UserTest",
		Email:    "test2@example.com",
		HashPass: "password123",
	}

	// добавляем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// потом ищем
	found, err := testRepo.FindByID(ctx, user.Id)

	// проверка выходных данных
	assert.NoError(t, err)
	assert.Equal(t, user.Name, found.Name)
	assert.Equal(t, user.Email, found.Email)
}

// тест на поиск пользователя по имейл
func TestUserRepo_FindByEmail(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:     "UserTest",
		Email:    "test3@example.com",
		HashPass: "password123",
	}

	// создаем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// ищем по имейлу
	found, err := testRepo.FindByEmail(ctx, user.Email)

	// проверка выходных данных
	assert.NoError(t, err)
	assert.Equal(t, user.Id, found.Id)
	assert.Equal(t, user.Email, found.Email)
}

// тест на обновление полей пользователя
func TestUserRepo_Update(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:     "TestUser",
		Email:    "test4@example.com",
		HashPass: "pass123",
	}

	// создаем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// обновляем имя и емайл
	newName := "UpdatedName"
	newEmail := "updatedemail@example.com"
	err = testRepo.Update(ctx, user.Id, &newName, &newEmail, nil)
	assert.NoError(t, err)

	// получаем обновленного пользователя
	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	// проверка выходных данных
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, newEmail, updated.Email)
	assert.Equal(t, user.HashPass, updated.HashPass)
}

// тест на обновление роли пользователя
func TestUserRepo_UpdateRole(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:  "TestUser",
		Email: "test5@example.com",
		Role:  "user",
	}

	// создаем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// обновляем роль
	err = testRepo.UpdateRole(ctx, "Admin", user.Id)
	assert.NoError(t, err)

	// получаем обновленного пользователя
	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	// проверка выходных данных
	assert.Equal(t, domain.RoleAdmin, updated.Role)
}

// тест на обновление рефреш токена
func TestUserRepo_UpdateRefreshToken(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:  "UserTest",
		Email: "test6@example.com",
	}

	// создаем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// обновляем токен
	token := "refresh_token_123"
	err = testRepo.UpdateRefreshToken(ctx, user.Id, token)
	assert.NoError(t, err)

	// получаем пользователя
	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	// проверка выходных данных
	assert.Equal(t, token, updated.RefreshToken)
	assert.NotNil(t, updated.RefreshTokenExpiresAt)
}

// тест на удаление пользователя
func TestUserRepo_Delete(t *testing.T) {
	ctx := context.Background()

	// задаем пользователя
	user := &domain.User{
		Name:  "UserTest",
		Email: "test7@example.com",
	}

	// создаем его
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// удаляем
	err = testRepo.Delete(ctx, user.Id)
	assert.NoError(t, err)

	// проверяем, что пользователь удален
	_, err = testRepo.FindByID(ctx, user.Id)
	assert.Error(t, err)
}

// тест на обновление несуществующего пользователя
func TestUserRepo_Update_NotFound(t *testing.T) {
	ctx := context.Background()

	//  задаем параметры
	fakeID := primitive.NewObjectID()
	newName := "UserName"

	// обновляем с неправильным айди
	err := testRepo.Update(ctx, fakeID, &newName, nil, nil)

	// ожидаем ошибку типа нот фаунд
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

// тест на очистку просроченного токена
func TestUserRepo_CleanExpiredRefreshTokens(t *testing.T) {
	ctx := context.Background()

	// создаем пользователя с истекшим токеном
	expiredTime := time.Now().Add(-24 * time.Hour)
	user := &domain.User{
		Name:                  "UserTest",
		Email:                 "test8@example.com",
		RefreshToken:          "expired_token",
		RefreshTokenExpiresAt: &expiredTime,
	}
	err := testRepo.Create(ctx, user)
	require.NoError(t, err)

	// очищаем токен
	err = testRepo.CleanExpiredRefreshTokens(ctx)
	assert.NoError(t, err)

	// находим пользователя
	updated, err := testRepo.FindByID(ctx, user.Id)
	assert.NoError(t, err)

	// проверяем, что токен очищен
	assert.Empty(t, updated.RefreshToken)
	assert.Nil(t, updated.RefreshTokenExpiresAt)
}
