package services

import (
	"context"
	"cook_book/internal/api/models"
	"cook_book/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeUseCaseInterface interface {
	CreateRecipe(ctx context.Context, recipe *domain.Recipe) error
	UpdateRecipe(
		ctx context.Context,
		id primitive.ObjectID,
		userId primitive.ObjectID,
		title *string,
		description *string,
		ingredients *[]domain.Ingredient,
		steps *[]string,
		prepTime *int,
		cooktime *int,
		servings *int,
		difficulty *string,
		tags *[]string,
	) error
	DeleteRecipe(ctx context.Context, id string, userId string) error
	GetRecipeByID(ctx context.Context, id string) (*domain.Recipe, error)
	GetRecipeByUser(ctx context.Context, userId string) ([]*domain.Recipe, error)
	ListRecipes(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]*domain.Recipe, int64, error)
	SearchByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error)
}

type UserUseCaseInterface interface {
	CreateUser(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*models.LoginResponse, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*domain.LoginResponse, error)
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	UpdateUser(ctx context.Context, id string, name *string, email *string, password *string) error
	UpdateUserRole(ctx context.Context, newRole domain.Role, id string) error
	DeleteUser(ctx context.Context, id string) error
	DeleteAnyUser(ctx context.Context, id string) error
	SignOut(ctx context.Context, userID string) error
	CleanExpiredTokens(ctx context.Context) error
}
