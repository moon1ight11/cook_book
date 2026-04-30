package user

import (
	"context"
	"cook_book/internal/api/jwt"
	"cook_book/internal/domain"
	"cook_book/internal/services/user/mocks"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// создание тестового юзкейса
func setupTest() (*userUseCase, *mocks.MockUserRepository, *mocks.MockJWTService) {
	mockRepo := new(mocks.MockUserRepository)
	mockJWT := new(mocks.MockJWTService)

	usecase := &userUseCase{
		repo:       mockRepo,
		jwtService: mockJWT,
	}

	return usecase, mockRepo, mockJWT
}

// создание тестового пользователя
func createTestUser() *domain.User {
	return &domain.User{
		Id:    primitive.NewObjectID(),
		Name:  "TestUser",
		Email: "test@example.com",
		Role:  domain.RoleBasic,
	}
}

// успешное создание пользователя
func TestCreateUser_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	req := &domain.CreateUserRequest{
		Name:     "TestUser1",
		Email:    "test@example.com",
		Password: "password123",
		Role:     domain.RoleBasic,
	}

	// настраиваем моки
	mockRepo.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.User")).Return(nil)

	// выполняем тест
	user, err := usecase.CreateUser(ctx, req)

	// проверяем
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, req.Name, user.Name)
	assert.Equal(t, req.Email, user.Email)
	assert.NotEmpty(t, user.HashPass)
	assert.Equal(t, req.Role, user.Role)

	// хэширован ли пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(req.Password))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// создание пользователя с существующим емайлом
func TestCreateUser_DuplicateEmail(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	existingUser := createTestUser()

	req := &domain.CreateUserRequest{
		Name:     "TestUser2",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", ctx, req.Email).Return(existingUser, nil)

	user, err := usecase.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "already exist")
	mockRepo.AssertExpectations(t)
}

// успеный логин
func TestLogin_Success(t *testing.T) {
	usecase, mockRepo, mockJWT := setupTest()
	ctx := context.Background()

	// приводим пароль в правильный вид
	password := "password123"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// настраиваем пользователя
	user := &domain.User{
		Id:       primitive.NewObjectID(),
		Name:     "TestUser3",
		Email:    "test3@example.com",
		HashPass: string(hashedPass),
		Role:     domain.RoleBasic,
	}

	// запрос на логин
	req := &domain.LoginRequest{
		Email:    "test3@example.com",
		Password: password,
	}

	// токены
	tokens := &jwt.TokenDetails{
		AccessToken:  "access_token_123",
		RefreshToken: "refresh_token_123",
	}

	// настраиваем моки
	mockRepo.On("FindByEmail", ctx, req.Email).Return(user, nil)
	mockRepo.On("UpdateRefreshToken", ctx, user.Id, "").Return(nil)
	mockJWT.On("GenerateTokens", user.Id.Hex(), user.Email, string(user.Role)).Return(tokens, nil)
	mockRepo.On("UpdateRefreshToken", ctx, user.Id, mock.AnythingOfType("string")).Return(nil)

	// тестируем
	response, err := usecase.Login(ctx, req)

	// проверяем
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, tokens.AccessToken, response.AccessToken)
	assert.Equal(t, tokens.RefreshToken, response.RefreshToken)
	assert.Equal(t, user.Id, response.User.Id)

	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

// логин с неправильным паролем
func TestLogin_InvalidPassword(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	// правильный пароль
	password := "password1"
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// делаем пользователя
	user := &domain.User{
		Id:       primitive.NewObjectID(),
		Email:    "test5@example.com",
		HashPass: string(hashedPass),
	}

	// запрос на логин с не тем паролем
	req := &domain.LoginRequest{
		Email:    "test5@example.com",
		Password: "password2",
	}

	// настраиваем моки
	mockRepo.On("FindByEmail", ctx, req.Email).Return(user, nil)

	// делаем запрос
	response, err := usecase.Login(ctx, req)

	// проверяем
	assert.Error(t, err)
	assert.Nil(t, response)
	mockRepo.AssertExpectations(t)
}

// успешное получение пользователя по айди
func TestGetUserByID_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	user := createTestUser()

	mockRepo.On("FindByID", ctx, user.Id).Return(user, nil)

	found, err := usecase.GetUserByID(ctx, user.Id.Hex())

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, user.Id, found.Id)
	assert.Equal(t, user.Email, found.Email)
	mockRepo.AssertExpectations(t)
}

// получение пользователя по не тому айди
func TestGetUserByID_InvalidID(t *testing.T) {
	usecase, _, _ := setupTest()
	ctx := context.Background()

	user, err := usecase.GetUserByID(ctx, "USERIDNOTFOUND")

	assert.Error(t, err)
	assert.Nil(t, user)
}

// успешное обновление пользователя
func TestUpdateUser_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	user := createTestUser()
	newName := "NewName"
	var nilEmail *string = nil
	var nilPass *string = nil

	mockRepo.On("FindByID", ctx, user.Id).Return(user, nil)
	mockRepo.On("Update", ctx, user.Id, &newName, nilEmail, nilPass).Return(nil)

	err := usecase.UpdateUser(ctx, user.Id.Hex(), &newName, nilEmail, nilPass)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// успешное обновление роли пользователя
func TestUpdateUserRole_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	user := createTestUser()
	newRole := domain.RoleAdmin

	mockRepo.On("UpdateRole", ctx, string(newRole), user.Id).Return(nil)

	err := usecase.UpdateUserRole(ctx, newRole, user.Id.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// обновление роли с неправильной входной ролью
func TestUpdateUserRole_InvalidRole(t *testing.T) {
	usecase, _, _ := setupTest()
	ctx := context.Background()

	invalidRole := domain.Role("ROLENOTFOUND")

	err := usecase.UpdateUserRole(ctx, invalidRole, "some_id")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not valid")
}

// успешное удаление пользователя
func TestDeleteUser_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	user := createTestUser()

	mockRepo.On("Delete", ctx, user.Id).Return(nil)

	err := usecase.DeleteUser(ctx, user.Id.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// успешный разлогин пользователя
func TestSignOut_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	user := createTestUser()

	mockRepo.On("UpdateRefreshToken", ctx, user.Id, "").Return(nil)

	err := usecase.SignOut(ctx, user.Id.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// успешнпая очистка протухших токенов
func TestCleanExpiredTokens_Success(t *testing.T) {
	usecase, mockRepo, _ := setupTest()
	ctx := context.Background()

	mockRepo.On("CleanExpiredRefreshTokens", ctx).Return(nil)

	err := usecase.CleanExpiredTokens(ctx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// успешное обновление токенов
func TestRefreshTokens_Success(t *testing.T) {
	usecase, mockRepo, mockJWT := setupTest()
	ctx := context.Background()

	user := createTestUser()
	refreshToken := "valid_refresh_token"

	hash := sha256.Sum256([]byte(refreshToken))
	user.RefreshToken = hex.EncodeToString(hash[:])

	claims := &jwt.RefreshTokenClaims{
		UserID: user.Id.Hex(),
	}

	newTokens := &jwt.TokenDetails{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	mockJWT.On("ValidateRefreshToken", refreshToken).Return(claims, nil)
	mockRepo.On("FindByID", ctx, user.Id).Return(user, nil)
	mockRepo.On("UpdateRefreshToken", ctx, user.Id, "").Return(nil)
	mockJWT.On("GenerateTokens", user.Id.Hex(), user.Email, string(user.Role)).Return(newTokens, nil)
	mockRepo.On("UpdateRefreshToken", ctx, user.Id, mock.AnythingOfType("string")).Return(nil)

	response, err := usecase.RefreshTokens(ctx, refreshToken)

	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, newTokens.AccessToken, response.AccessToken)
	assert.Equal(t, newTokens.RefreshToken, response.RefreshToken)
}

// неудачное обновление токенов изза неправильного рефреша
func TestRefreshTokens_InvalidToken(t *testing.T) {
	usecase, _, mockJWT := setupTest()
	ctx := context.Background()

	invalidToken := "invalid_token"

	mockJWT.On("ValidateRefreshToken", invalidToken).Return(nil, errors.New("invalid token"))

	response, err := usecase.RefreshTokens(ctx, invalidToken)

	assert.Error(t, err)
	assert.Nil(t, response)
	mockJWT.AssertExpectations(t)
}
