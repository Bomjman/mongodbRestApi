package helper

import (
	"context"
	"log"
	"time"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Configuration struct {
	Username         string
	Password		 string
	Host			 string
	Database 		 string
}

func getEnvVars() Configuration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configuration := Configuration{
		os.Getenv("USERNAME"),
		os.Getenv("PASSWORD"),
		os.Getenv("HOST"),
		os.Getenv("DATABASE"),

	}

	return configuration
}

var client *mongo.Client

func Connect() *mongo.Collection {
	// export credentials
	getEnvVars()

	config := getEnvVars()
	uri := "mongodb+srv://" + config.Username + ":" + config.Password + "@" + config.Host

	clientOptions := options.Client().ApplyURI(uri)
	client, _ = mongo.NewClient(clientOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Could't connect to the database", err)
	} else {
		log.Println("Connected to MongoDB")
	}

	collection := client.Database(config.Database).Collection("profiles")

	return collection
}
