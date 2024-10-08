package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	openai "github.com/sashabaranov/go-openai"
)

var client *mongo.Client
var batchCollection *mongo.Collection

func InitMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://admin:password@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database("batchgpt")
	
	// Check if the collection exists
	collections, err := database.ListCollectionNames(ctx, bson.M{"name": "batch_logs"})
	if err != nil {
		log.Fatal(err)
	}

	if len(collections) == 0 {
		// Collection doesn't exist, so create it with the index
		err = database.CreateCollection(ctx, "batch_logs")
		if err != nil {
			log.Fatal(err)
		}

		batchCollection = database.Collection("batch_logs")

		_, err = batchCollection.Indexes().CreateOne(
			ctx,
			mongo.IndexModel{
				Keys:    bson.D{{Key: "batch.id", Value: 1}},
				Options: options.Index().SetUnique(false),
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Created 'batch_logs' collection with index on batch.id")
	} else {
		// Collection exists, assume index is present
		batchCollection = database.Collection("batch_logs")
		log.Println("'batch_logs' collection already exists, assuming index is present")
	}

	log.Println("Connected to MongoDB")
}

func LogBatchResponse(batchResponse openai.BatchResponse) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	document := bson.M{
		"batch":     batchResponse.Batch,
		"timestamp": time.Now(),
	}

	_, err := batchCollection.InsertOne(ctx, document)
	return err
}

func GetLatestBatchStatus(batchID string) (openai.Batch, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var result struct {
        Batch openai.Batch `bson:"batch"`
    }
    err := batchCollection.FindOne(
        ctx,
        bson.M{"batch.id": batchID},
        options.FindOne().SetSort(bson.D{{Key: "timestamp", Value: -1}}),
    ).Decode(&result)

    if err != nil {
        if err == mongo.ErrNoDocuments {
            return openai.Batch{}, mongo.ErrNoDocuments
        }
        return openai.Batch{}, err
    }

    return result.Batch, nil
}
