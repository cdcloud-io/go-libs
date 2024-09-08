package mongoclient

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client wraps the MongoDB client and provides additional functionality
// In a Hexagonal Architecture, this acts as the **Adapter** for MongoDB.
type Client struct {
	*mongo.Client
}

// ClientOptions represents options for creating a new Client
// These options abstract connection details that can be passed from outside
// the business logic, allowing for flexibility in different environments.
type ClientOptions struct {
	URI                    string
	ConnectTimeout         time.Duration
	ServerSelectionTimeout time.Duration
}

// QueryParams abstracts the MongoDB query parameters
// This struct can be used as a **Port**, as it defines how external systems
// can communicate with the MongoDB database without knowing its internal workings.
type QueryParams struct {
	Database   string
	Collection string
	Filter     bson.M
}

// NewClient creates and returns a new Client with the given options
// This function allows for external systems to create an instance of a MongoDB client.
// In a hexagonal architecture, this might be called from an Adapter that integrates with the infrastructure layer.
func NewClient(opts ClientOptions) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	clientOpts := options.Client().ApplyURI(opts.URI).
		SetServerSelectionTimeout(opts.ServerSelectionTimeout)

	// Connect to MongoDB using the specified options
	mongoClient, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to ensure the connection is successful
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		mongoClient.Disconnect(ctx)
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	// Return the wrapped MongoDB client
	return &Client{Client: mongoClient}, nil
}

// Close disconnects the client from MongoDB
// This is part of the infrastructure layer, closing the connection to the database.
func (c *Client) Close(ctx context.Context) error {
	return c.Disconnect(ctx)
}

// QueryOne executes a query to find a single document using QueryParams
// This abstracts the MongoDB-specific query logic, making it reusable by passing `QueryParams`.
// It acts as an **Adapter** method that can be called from the application core via Ports.
func (c *Client) QueryOne(ctx context.Context, params QueryParams, result interface{}) error {
	collection := c.Database(params.Database).Collection(params.Collection)

	// Execute the FindOne query based on the filter provided in QueryParams
	err := collection.FindOne(ctx, params.Filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		return nil // Return nil if no documents are found
	}
	if err != nil {
		return fmt.Errorf("failed to execute FindOne query: %w", err)
	}

	return nil
}

// QueryMany executes a query to find multiple documents using QueryParams
// This function can be used to find multiple documents and returns them as an array of interfaces.
// It's abstracted, so the core application does not need to handle MongoDB-specific logic.
func (c *Client) QueryMany(ctx context.Context, params QueryParams) ([]interface{}, error) {
	collection := c.Database(params.Database).Collection(params.Collection)

	// Execute the Find query and get a cursor to iterate over the results
	cursor, err := collection.Find(ctx, params.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Find query: %w", err)
	}
	defer cursor.Close(ctx)

	var results []interface{}

	// Decode all the documents returned by the query
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to decode query results: %w", err)
	}

	return results, nil
}

// InsertOne inserts a single document using QueryParams
// This function allows for inserting a document into MongoDB while abstracting the MongoDB-specific logic.
func (c *Client) InsertOne(ctx context.Context, params QueryParams, document interface{}) (*mongo.InsertOneResult, error) {
	// Insert the document into the specified collection
	result, err := c.Database(params.Database).Collection(params.Collection).InsertOne(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}
	return result, nil
}

// UpdateOne updates a single document using QueryParams
// This abstracts the update operation to ensure the core logic does not depend on MongoDB internals.
func (c *Client) UpdateOne(ctx context.Context, params QueryParams, update interface{}) (*mongo.UpdateResult, error) {
	// Update the document based on the filter provided in QueryParams
	result, err := c.Database(params.Database).Collection(params.Collection).UpdateOne(ctx, params.Filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to update document: %w", err)
	}
	return result, nil
}

// DeleteOne deletes a single document using QueryParams
// Abstracts the delete operation, keeping the core logic independent of the MongoDB implementation.
func (c *Client) DeleteOne(ctx context.Context, params QueryParams) (*mongo.DeleteResult, error) {
	// Delete the document based on the filter provided in QueryParams
	result, err := c.Database(params.Database).Collection(params.Collection).DeleteOne(ctx, params.Filter)
	if err != nil {
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}
	return result, nil
}

// QueryMongoDBStruct executes a MongoDB query with abstracted parameters
// and decodes the result directly into the provided struct.
func (c *Client) QueryMongoDBStruct(ctx context.Context, params QueryParams, result interface{}) error {
	collection := c.Database(params.Database).Collection(params.Collection)

	// Execute the query and decode the result into the provided struct
	err := collection.FindOne(ctx, params.Filter).Decode(result)
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("no documents found")
	}
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}

	return nil
}

/*

// QueryMongoDB executes a MongoDB query with abstracted parameters
// This method allows for more generic query operations, using the `map[string]interface{}` to handle unknown document structures.
func (c *Client) QueryMongoDB(ctx context.Context, params QueryParams) (*map[string]interface{}, error) {
	collection := c.Database(params.Database).Collection(params.Collection)

	// Store the result in a map[string]interface{} since the structure is unknown
	var result map[string]interface{}

	// Execute the query and decode the result into the generic map
	err := collection.FindOne(ctx, params.Filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no documents found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query MongoDB: %w", err)
	}

	return &result, nil
}
*/
