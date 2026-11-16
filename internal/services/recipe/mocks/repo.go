package mocks

import (
	"context"
	"cook_book/internal/domain"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockRecipeRepository struct {
	mock.Mock
}

func (m *MockRecipeRepository) Create(ctx context.Context, recipe *domain.Recipe) error {
	args := m.Called(ctx, recipe)
	return args.Error(0)
}

func (m *MockRecipeRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipe, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Recipe), args.Error(1)
}

func (m *MockRecipeRepository) FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int64) ([]*domain.Recipe, int64, error) {
	args := m.Called(ctx, filter, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Recipe), args.Get(1).(int64), args.Error(2)
}

func (m *MockRecipeRepository) FindByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error) {
	args := m.Called(ctx, tags)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Recipe), args.Error(1)
}

func (m *MockRecipeRepository) FindByUser(ctx context.Context, userId primitive.ObjectID) ([]*domain.Recipe, error) {
	args := m.Called(ctx, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Recipe), args.Error(1)
}

func (m *MockRecipeRepository) Update(
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
) error {
	args := m.Called(ctx, id, userId, title, description, ingredients, steps, prepTime, cooktime, servings, difficulty, tags)
	return args.Error(0)
}

func (m *MockRecipeRepository) Delete(ctx context.Context, id primitive.ObjectID, userId primitive.ObjectID) error {
	args := m.Called(ctx, id, userId)
	return args.Error(0)
}
