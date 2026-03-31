package recipes

import (
	"context"
	"cook_book/internal/domain"
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
		return err
	}

	recipe.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

// обновление рецепта
func (r *recipeRepo) Update(ctx context.Context, recipe *domain.Recipe) error {
	// устанавливаем временную метку
	recipe.UpdatedAt = time.Now()

	filter := bson.M{"_id": recipe.Id}
	update := bson.M{"$set": recipe}

	// обновляем запись
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// если ничего не обновили - ошибка
	if result.MatchedCount == 0 {
		return domain.ErrRecipeNotFound
	}

	return nil
}

// удаление записи из бд
func (r *recipeRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	// фильтруем по конкретному id
	filter := bson.M{"_id": id}

	// удаляем запись
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// если ничего не удалилось - ошибка
	if result.DeletedCount == 0 {
		return domain.ErrRecipeNotFound
	}

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
		return nil, err
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
		return nil, 0, err
	}

	// отображаем с пагинацией
	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

		// получаем указатель на найденные объекты
	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// из указателя выгружаем рецепты
	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, 0, err
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
		return nil, err
	}
	defer cursor.Close(ctx)

	// выгружаем рецепты из указателя
	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, err
	}

	return recipes, nil
}
