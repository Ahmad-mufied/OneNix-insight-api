package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google-custom-search/model"
	"log"
)

type MongoRepository struct {
	newsCollection *mongo.Collection
}

func NewMongoRepository(collection *mongo.Collection) *MongoRepository {
	return &MongoRepository{newsCollection: collection}
}

func (r *MongoRepository) SaveNews(query string, country, degree, major string, results []*model.News) error {

	ctx := context.TODO()

	// Create a compound unique index
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "country", Value: 1}, // 1 for ascending order
			{Key: "degree", Value: 1},
			{Key: "major", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	// Create the index
	_, err := r.newsCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("Error creating index: %v", err)
	}

	document := &model.SearchResult{
		Country: country,
		Degree:  degree,
		Major:   major,
		Results: results,
	}

	_, err = r.newsCollection.InsertOne(ctx, document)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Printf("Results for query %s already exist", query)
			return nil
		} else {
			return fmt.Errorf("error saving results for query %s: %v", query, err)
		}
	}

	log.Printf("Saved results for query %s", query)

	return nil
}

func (r *MongoRepository) DropCollection() error {
	ctx := context.TODO()
	err := r.newsCollection.Drop(ctx)
	if err != nil {
		return fmt.Errorf("error dropping collection: %v", err)
	}
	return nil
}

// List retrieves all news records matching the provided filters
func (r *MongoRepository) List(ctx context.Context, filters map[string]string) ([]model.News, error) {
	filter := bson.M{}
	for key, value := range filters {
		if value != "" {
			filter[key] = value
		}
	}

	cursor, err := r.newsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var newsList []model.News
	var results []*model.SearchResult
	if err := cursor.All(ctx, &results); err != nil {

		return nil, err
	}

	for _, result := range results {
		for _, news := range result.Results {
			newsList = append(newsList, *news)
		}
	}

	return newsList, nil
}

// GetByID fetches a news record by its ID
func (r *MongoRepository) GetByID(ctx context.Context, id string) (*model.News, error) {
	var news model.News
	err := r.newsCollection.FindOne(ctx, bson.M{"id": id}).Decode(&news)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &news, err
}

// Update modifies an existing news record by its ID
func (r *MongoRepository) Update(ctx context.Context, id string, updatedData model.News) error {
	filter := bson.M{"id": id}
	update := bson.M{"$set": updatedData}
	result, err := r.newsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("news not found")
	}
	return nil
}

// Delete removes a single news record by its ID
func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	result, err := r.newsCollection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("news not found")
	}
	return nil
}

// DeleteAll removes all news records from the collection
func (r *MongoRepository) DeleteAll(ctx context.Context) error {
	result, err := r.newsCollection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no news records found to delete")
	}
	return nil
}
