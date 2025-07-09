package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func InitMongoDB(uri string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ MongoDB connection error:", err)
	}

	// Optional: check connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	MongoClient = client
	log.Println("✅ Connected to MongoDB")

	//return client.Database("Notification-system-v1") // or make dbName a parameter if needed

	db := client.Database("Notification-system-v1")

    // 🔐 Create unique index on notification_id
    ensureUniqueIndex(db, "Notifications", "notification_id")

    return db
}

func ensureUniqueIndex(db *mongo.Database, collectionName, field string) {
    coll := db.Collection(collectionName)
    indexModel := mongo.IndexModel{
        Keys: bson.D{{Key: field, Value: 1}}, // 1 = ascending index
        Options: options.Index().SetUnique(true),
    }

    _, err := coll.Indexes().CreateOne(context.TODO(), indexModel)
    if err != nil {
        log.Fatalf("❌ Failed to create unique index on %s.%s: %v", collectionName, field, err)
    }
    log.Printf("✅ Unique index ensured on %s.%s\n", collectionName, field)
}

