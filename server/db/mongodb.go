package db

import (
	"batch-gpt/server/logger"

	"context"
	"fmt"
	"log"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var batchCollection *mongo.Collection
var cachedResponsesCollection *mongo.Collection

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

	cachedResponsesCollection = database.Collection("cached_responses")

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

func GetDanglingBatches() ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // First, get all unique batch IDs
    pipeline := mongo.Pipeline{
        {{Key: "$group", Value: bson.D{{Key: "_id", Value: "$batch.id"}}}},
    }

    cursor, err := batchCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate unique batch IDs: %w", err)
    }
    defer cursor.Close(ctx)

    var results []struct {
        ID string `bson:"_id"`
    }
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode aggregate results: %w", err)
    }

    var danglingBatches []string

    for _, result := range results {
        batchID := result.ID
        latestStatus, err := GetLatestBatchStatus(batchID)
        if err != nil {
            logger.WarnLogger.Printf("Failed to get latest status for batch %s: %v", batchID, err)
            continue
        }

        if latestStatus.Status != "completed" && latestStatus.Status != "failed" {
            danglingBatches = append(danglingBatches, batchID)
        }
    }

    return danglingBatches, nil
}

// func CacheResponses(batchID string, responses []openai.ChatCompletionResponse) error {
//     ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//     defer cancel()

//     var documents []interface{}
//     for _, response := range responses {
//         document := bson.M{
//             "batch_id":  batchID,
//             "response":  response,
//             "timestamp": time.Now(),
//         }
//         documents = append(documents, document)
//     }

//     _, err := cachedResponsesCollection.InsertMany(ctx, documents)
//     if err != nil {
//         return fmt.Errorf("failed to insert cached responses: %w", err)
//     }

//     return nil
// }

func GetAllBatchStatuses() ([]openai.BatchResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    pipeline := mongo.Pipeline{
        {{Key: "$sort", Value: bson.D{{Key: "batch.created_at", Value: -1}}}},
        {{Key: "$group", Value: bson.D{
            {Key: "_id", Value: "$batch.id"},
            {Key: "batch", Value: bson.D{{Key: "$first", Value: "$batch"}}},
        }}},
    }

    cursor, err := batchCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, fmt.Errorf("failed to aggregate batch statuses: %w", err)
    }
    defer cursor.Close(ctx)

    var results []openai.BatchResponse
    if err = cursor.All(ctx, &results); err != nil {
        return nil, fmt.Errorf("failed to decode aggregate results: %w", err)
    }

    return results, nil
}

func GetCachedResponse(hash string) (openai.ChatCompletionResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var result struct {
        Response openai.ChatCompletionResponse `bson:"response"`
    }
    err := cachedResponsesCollection.FindOne(ctx, bson.M{"hash": hash}).Decode(&result)
    if err != nil {
        return openai.ChatCompletionResponse{}, err
    }
    return result.Response, nil
}

func CacheRequestResponse(hash string, request openai.ChatCompletionRequest, response openai.ChatCompletionResponse) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    document := bson.M{
        "hash":      hash,
        "request":   request,
        "response":  response,
        "timestamp": time.Now(),
    }

    _, err := cachedResponsesCollection.InsertOne(ctx, document)
    return err
}
