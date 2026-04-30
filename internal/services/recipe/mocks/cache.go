package mocks

import (
	"cook_book/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockRecipeCache struct {
	mock.Mock
}

func (m *MockRecipeCache) Get(id string) (*domain.Recipe, bool) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*domain.Recipe), args.Bool(1)
}

func (m *MockRecipeCache) Set(id string, recipe *domain.Recipe) {
	m.Called(id, recipe)
}

func (m *MockRecipeCache) Delete(id string) {
	m.Called(id)
}

func (m *MockRecipeCache) GetByTags(tagsKey string) ([]*domain.Recipe, bool) {
	args := m.Called(tagsKey)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).([]*domain.Recipe), args.Bool(1)
}

func (m *MockRecipeCache) SetByTags(tagsKey string, recipes []*domain.Recipe) {
	m.Called(tagsKey, recipes)
}

func (m *MockRecipeCache) Clear() {}

func (m *MockRecipeCache) Stop() {}
