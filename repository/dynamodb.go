package repository

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"google-custom-search/model"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DynamoDBClient represents a wrapper for DynamoDB service
type DynamoDBClient struct {
	Client *dynamodb.Client
	Table  string
}

// NewDynamoDBClient initializes and returns a DynamoDBClient
func NewDynamoDBClient(region string) (*DynamoDBClient, error) {
	// Load AWS configuration
	awsConfig, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Printf("Unable to load AWS config: %v", err)
		return nil, err
	}

	// Initialize DynamoDB client
	client := dynamodb.NewFromConfig(awsConfig)

	return &DynamoDBClient{
		Client: client,
		Table:  "News", // Replace with your DynamoDB table name
	}, nil
}

// Check if the news already exists in the database
func (db *DynamoDBClient) CheckNewsExists(id string) (bool, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(db.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	result, err := db.Client.GetItem(context.TODO(), input)
	if err != nil {
		return false, err
	}

	// If no item is returned, it doesn't exist
	return result.Item != nil, nil
}

// Save news to DynamoDB
func (db *DynamoDBClient) SaveNews(news model.News) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String(db.Table),
		Item: map[string]types.AttributeValue{
			"id":        &types.AttributeValueMemberS{Value: news.ID},
			"title":     &types.AttributeValueMemberS{Value: news.Title},
			"date":      &types.AttributeValueMemberS{Value: news.Date},
			"thumbnail": &types.AttributeValueMemberS{Value: news.Thumbnail},
			"snippet":   &types.AttributeValueMemberS{Value: news.Snippet},
			"link":      &types.AttributeValueMemberS{Value: news.Link},
		},
	}

	_, err := db.Client.PutItem(context.TODO(), input)
	if err != nil {
		log.Printf("Error saving news: %v\n", err)
		return err
	}
	return nil
}

func (db *DynamoDBClient) GetAllNews() ([]model.News, error) {
	var newsList []model.News

	// Scan the DynamoDB table
	input := &dynamodb.ScanInput{
		TableName: aws.String(db.Table),
	}

	result, err := db.Client.Scan(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to scan table: %w", err)
	}

	// Process each item in the result
	for _, item := range result.Items {
		var news model.News

		// Safely access and type assert each attribute
		if titleAttr, ok := item["title"].(*types.AttributeValueMemberS); ok {
			news.Title = titleAttr.Value
		} else {
			news.Title = "Unknown Title" // Default or fallback
		}

		if linkAttr, ok := item["link"].(*types.AttributeValueMemberS); ok {
			news.Link = linkAttr.Value
		} else {
			news.Link = "Unknown Link"
		}

		if snippetAttr, ok := item["snippet"].(*types.AttributeValueMemberS); ok {
			news.Snippet = snippetAttr.Value
		} else {
			news.Snippet = "No snippet available"
		}

		if thumbnailAttr, ok := item["thumbnail"].(*types.AttributeValueMemberS); ok {
			news.Thumbnail = thumbnailAttr.Value
		} else {
			news.Thumbnail = ""
		}

		if dateAttr, ok := item["date"].(*types.AttributeValueMemberS); ok {
			news.Date = dateAttr.Value
		} else {
			news.Date = "Unknown Date"
		}

		if idAttr, ok := item["newsID"].(*types.AttributeValueMemberS); ok {
			news.ID = idAttr.Value
		} else {
			news.ID = "Unknown ID"
		}

		// Append to the list
		newsList = append(newsList, news)
	}

	return newsList, nil
}
