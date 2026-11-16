package user

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"cook_book/internal/api/models"
	"cook_book/internal/domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func (u *userUseCase) CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error) {
	existingUser, _ := u.repo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("CreateUser: user with email %s already exists", req.Email)
	}

	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("CreateUser: invalid email format")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		HashPass: string(hashedPass),
		Role:     req.Role,
	}

	if err = u.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}

	return user, nil
}

func (u *userUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*models.LoginResponse, error) {
	user, err := u.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.HashPass), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	if err = u.repo.UpdateRefreshToken(ctx, user.Id, ""); err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	tokens, err := u.jwtService.GenerateTokens(user.Id.Hex(), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	hash := sha256.Sum256([]byte(tokens.RefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	if err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken); err != nil {
		return nil, fmt.Errorf("Login: %w", err)
	}

	return &models.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (u *userUseCase) RefreshTokens(ctx context.Context, refreshToken string) (*domain.LoginResponse, error) {
	claims, err := u.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	objID, err := primitive.ObjectIDFromHex(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	hash := sha256.Sum256([]byte(refreshToken))
	requestHash := hex.EncodeToString(hash[:])

	if subtle.ConstantTimeCompare([]byte(user.RefreshToken), []byte(requestHash)) != 1 {
		return nil, fmt.Errorf("RefreshTokens: invalid token")
	}

	if err = u.repo.UpdateRefreshToken(ctx, user.Id, ""); err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	tokens, err := u.jwtService.GenerateTokens(user.Id.Hex(), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	hash = sha256.Sum256([]byte(tokens.RefreshToken))
	hashedToken := hex.EncodeToString(hash[:])

	if err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken); err != nil {
		return nil, fmt.Errorf("RefreshTokens: %w", err)
	}

	return &domain.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		User:         user,
	}, nil
}

func (u *userUseCase) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: %w", err)
	}

	return user, nil
}

func (u *userUseCase) UpdateUser(
	ctx context.Context,
	id string,
	name *string,
	email *string,
	password *string,
) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	user, err := u.repo.FindByID(ctx, objID)
	if err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	needsNewTokens := false
	var hashedPass string

	if name != nil && strings.TrimSpace(*name) == "" {
		return fmt.Errorf("UpdateUser: new name cannot be empty")
	}

	finalEmail := user.Email
	if email != nil && *email != "" && *email != user.Email {
		existingUser, _ := u.repo.FindByEmail(ctx, *email)
		if existingUser != nil && existingUser.Id != user.Id {
			return fmt.Errorf("UpdateUser: email %s already taken", *email)
		}

		if !isValidEmail(*email) {
			return fmt.Errorf("UpdateUser: invalid email format")
		}

		needsNewTokens = true
		finalEmail = *email
	}

	var hashedPassPtr *string
	if password != nil {
		if strings.TrimSpace(*password) == "" {
			return fmt.Errorf("UpdateUser: new password cannot be empty")
		}

		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("UpdateUser: %w", err)
		}
		hashedPass = string(hashedBytes)
		hashedPassPtr = &hashedPass
		needsNewTokens = true
	}

	if err = u.repo.Update(ctx, user.Id, name, email, hashedPassPtr); err != nil {
		return fmt.Errorf("UpdateUser: %w", err)
	}

	if needsNewTokens {
		tokens, err := u.jwtService.GenerateTokens(user.Id.Hex(), finalEmail, string(user.Role))
		if err != nil {
			return fmt.Errorf("UpdateUser: %w", err)
		}

		hash := sha256.Sum256([]byte(tokens.RefreshToken))
		hashedToken := hex.EncodeToString(hash[:])

		if err = u.repo.UpdateRefreshToken(ctx, user.Id, hashedToken); err != nil {
			return fmt.Errorf("UpdateUser: %w", err)
		}
	}

	return nil
}

func (u *userUseCase) UpdateUserRole(ctx context.Context, newRole domain.Role, id string) error {
	if newRole != domain.RoleAdmin && newRole != domain.RoleBasic {
		return fmt.Errorf("UpdateUserRole: invalid role: %s", string(newRole))
	}

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("UpdateUserRole: %w", err)
	}

	if err = u.repo.UpdateRole(ctx, string(newRole), objID); err != nil {
		return fmt.Errorf("UpdateUserRole: %w", err)
	}

	return nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	if err = u.repo.Delete(ctx, objID); err != nil {
		return fmt.Errorf("DeleteUser: %w", err)
	}

	return nil
}

func (u *userUseCase) DeleteAnyUser(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("DeleteAnyUser: %w", err)
	}

	if err = u.repo.DeleteAny(ctx, objID); err != nil {
		return fmt.Errorf("DeleteAnyUser: %w", err)
	}

	return nil
}

func (u *userUseCase) SignOut(ctx context.Context, userID string) error {
	objID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("SignOut: %w", err)
	}

	if err = u.repo.UpdateRefreshToken(ctx, objID, ""); err != nil {
		return fmt.Errorf("SignOut: %w", err)
	}

	return nil
}

func (u *userUseCase) CleanExpiredTokens(ctx context.Context) error {
	if err := u.repo.CleanExpiredRefreshTokens(ctx); err != nil {
		return fmt.Errorf("CleanExpiredTokens: %w", err)
	}
	return nil
}
