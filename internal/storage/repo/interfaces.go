package repo

import (
	"context"
	"cook_book/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecipeRepositoryInterface interface {
	Create(ctx context.Context, recipe *domain.Recipe) error
	Update(
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
	Delete(ctx context.Context, id primitive.ObjectID, userId primitive.ObjectID) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipe, error)
	FindByUser(ctx context.Context, userId primitive.ObjectID) ([]*domain.Recipe, error)
	FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int64) ([]*domain.Recipe, int64, error)
	FindByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error)
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Update(ctx context.Context, id primitive.ObjectID, name *string, email *string, pass *string) error
	UpdateRole(ctx context.Context, newRole string, id primitive.ObjectID) error
	UpdateRefreshToken(ctx context.Context, id primitive.ObjectID, refreshToken string) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	DeleteAny(ctx context.Context, id primitive.ObjectID) error
	CleanExpiredRefreshTokens(ctx context.Context) error
}
