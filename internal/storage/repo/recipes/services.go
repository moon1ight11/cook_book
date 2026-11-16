package recipes

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"cook_book/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *recipeRepo) Create(ctx context.Context, recipe *domain.Recipe) error {
	recipe.CreatedAt = time.Now()
	recipe.Id = primitive.NewObjectID()

	result, err := r.collection.InsertOne(ctx, recipe)
	if err != nil {
		return fmt.Errorf("recipeCreate: %w", err)
	}

	recipe.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *recipeRepo) FindByID(ctx context.Context, id primitive.ObjectID) (*domain.Recipe, error) {
	var recipe domain.Recipe

	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&recipe)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("recipeFindByID: recipe not found")
		}
		return nil, fmt.Errorf("recipeFindByID: %w", err)
	}

	return &recipe, nil
}

func (r *recipeRepo) FindAll(ctx context.Context, filter map[string]interface{}, limit, offset int64) ([]*domain.Recipe, int64, error) {
	mongoFilter := bson.M{}

	if title, ok := filter["title"]; ok {
		mongoFilter["title"] = bson.M{"$regex": regexp.QuoteMeta(title.(string)), "$options": "i"}
	}

	if difficulty, ok := filter["difficulty"]; ok {
		mongoFilter["difficulty"] = bson.M{"$regex": regexp.QuoteMeta(difficulty.(string)), "$options": "i"}
	}

	total, err := r.collection.CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("recipeFindAll: %w", err)
	}

	findOptions := options.Find().
		SetLimit(limit).
		SetSkip(offset).
		SetSort(bson.M{"created_at": -1})

	cursor, err := r.collection.Find(ctx, mongoFilter, findOptions)
	if err != nil {
		return nil, 0, fmt.Errorf("recipeFindAll: %w", err)
	}
	defer cursor.Close(ctx)

	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, 0, fmt.Errorf("recipeFindAll: %w", err)
	}

	return recipes, total, nil
}

func (r *recipeRepo) FindByTags(ctx context.Context, tags []string) ([]*domain.Recipe, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"tags": bson.M{"$in": tags}})
	if err != nil {
		return nil, fmt.Errorf("recipeFindByTags: %w", err)
	}
	defer cursor.Close(ctx)

	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, fmt.Errorf("recipeFindByTags: %w", err)
	}

	return recipes, nil
}

func (r *recipeRepo) FindByUser(ctx context.Context, userId primitive.ObjectID) ([]*domain.Recipe, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"owner_id": userId})
	if err != nil {
		return nil, fmt.Errorf("recipeFindByUser: %w", err)
	}
	defer cursor.Close(ctx)

	var recipes []*domain.Recipe
	if err = cursor.All(ctx, &recipes); err != nil {
		return nil, fmt.Errorf("recipeFindByUser: %w", err)
	}

	return recipes, nil
}

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
		update["prep_time"] = *prepTime
	}
	if cooktime != nil {
		update["cook_time"] = *cooktime
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
		return fmt.Errorf("recipeUpdate: nothing to update")
	}

	update["updated_at"] = time.Now()

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id, "owner_id": userId}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("recipeUpdate: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("recipeUpdate: recipe with id=%s not found", id.Hex())
	}

	return nil
}

func (r *recipeRepo) Delete(ctx context.Context, id primitive.ObjectID, userId primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id, "owner_id": userId})
	if err != nil {
		return fmt.Errorf("recipeDelete: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("recipeDelete: recipe not found")
	}

	return nil
}
