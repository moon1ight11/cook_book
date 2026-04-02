package recipes

import (
	"context"
	"cook_book/internal/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// создание записи рецепта
func (r *recipeRepo) Create(ctx context.Context, recipe *domain.Recipe) error {
	// устанавливаем метку создания
	recipe.CreatedAt = time.Now()

	// получаем результат из бд
	result, err := r.collection.InsertOne(ctx, recipe)
	if err != nil {
		return fmt.Errorf("error in recipeCreate: %w", err)
	}

	recipe.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

// поиск рецепта по id
func (r *recipeRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipe, error) {
	var recipe domain.Recipe

	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrRecipeNotFound
		}
		return nil, fmt.Errorf("error in recipeFindById: %w", err)
	}

	return &recipe, nil
}

// получение всех рецептов в сервисе с фильтром
func (r *recipeRepo) FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int64) ([]*domain.Recipe, int64, error) {
	// Подбраем фильтр
	mongoFilter := bson.M{}

	// фильтр по названию
	if title, ok := filter["title"]; ok {
		mongoFilter["title"] = bson.M{"$regex": title, "$options": "i"}
	}

	// фильтр по сложности
	if difficulty, ok := filter["difficulty"]; ok {
		mongoFilter["difficulty"] = difficulty
	}

	// фильтр по тегам
	if tags, ok := filter["tags"]; ok {
		mongoFilter["tags"] = bson.M{"$in": tags}
	}

	// получаем количество найденных рецептов
	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("error in recipeFindAll: %w", err)
	}

	// отображаем с пагинацией
	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	// получаем указатель на найденные объекты
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("error in recipeFindAll: %w", err)
	}
	defer cursor.Close(ctx)

	// из указателя выгружаем рецепты
	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, 0, fmt.Errorf("error in recipeFindAll: %w", err)
	}

	return recipes, total, nil
}

// поиск записей по тегам
func (r *recipeRepo) FindByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error) {
	// выбираем теги
	filter := bson.M{"tags": bson.M{"$in": tags}}

	// получаем указатель на рецепты
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error in recipeFindByTags: %w", err)
	}
	defer cursor.Close(ctx)

	// выгружаем рецепты из указателя
	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, fmt.Errorf("error in recipeFindByTags: %w", err)
	}

	return recipes, nil
}

// получение всех рецептов конкретного пользователя
func (r *recipeRepo) FindByUser(ctx context.Context, userId primitive.ObjectID) ([]*domain.Recipe, error) {
	// выбираем фильтр
	filter := bson.M{"owner_id": bson.M{"$in": userId}}

	// получаем указатель на рецепты
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error in recipeFindByUser: %w", err)
	}
	defer cursor.Close(ctx)

	// выгружаем рецепты из указателя
	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, fmt.Errorf("error in recipeFindByUser: %w", err)
	}

	return recipes, nil
}

// обновление записи
func (r *recipeRepo) Update(
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
	update := bson.M{}

	if title != nil {
		update["title"] = *title
	}
	if description != nil {
		update["description"] = *description
	}
	if ingredients != nil {
		update["ingredients"] = *ingredients
	}
	if steps != nil {
		update["steps"] = *steps
	}
	if prepTime != nil {
		update["preptime"] = *prepTime
	}
	if cooktime != nil {
		update["cooktime"] = *cooktime
	}
	if servings != nil {
		update["servings"] = *servings
	}
	if difficulty != nil {
		update["difficulty"] = *difficulty
	}
	if tags != nil {
		update["tags"] = *tags
	}

	if len(update) == 0 {
		return fmt.Errorf("error in recipeUpdate: %s", "nothing to update")
	}

	update["updatedat"] = time.Now()

	filter := bson.M{
		"_id":      id,
		"owner_id": userId,
	}

	result, err := r.collection.UpdateOne(ctx, filter, bson.M{
		"$set": update,
	})
	if err != nil {
		return fmt.Errorf("error in recipeUpdate: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("error in recipeUpdate: recipe with id=%s not found", id)
	}

	return nil
}

// удаление записи из бд
func (r *recipeRepo) Delete(ctx context.Context, id primitive.ObjectID, userId primitive.ObjectID) error {
	// фильтруем по конкретному id + id пользователя
	filter := bson.M{
		"_id":      id,
		"owner_id": userId,
	}

	// удаляем запись
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error in recipeDelete: %w", err)
	}

	// если ничего не удалилось - ошибка
	if result.DeletedCount == 0 {
		return domain.ErrRecipeNotFound
	}

	return nil
}
