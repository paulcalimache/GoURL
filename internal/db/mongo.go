package db

import (
	"context"
	"errors"

	"github.com/paulcalimache/gourl/internal/model"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "urls"

type MongoDB struct {
	database *mongo.Database
}

func NewMongoDB() *MongoDB {
	credential := options.Credential{
		Username: "root",
		Password: "root",
	}
	mongoURI := "mongodb://localhost:27017"

	// Check the connection
	client, err := mongo.Connect(
		context.TODO(),
		options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		log.Fatal().Err(err).Send()
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	log.Info().Str("mongo_uri", mongoURI).Msg("Connected to Mongo database !")
	return &MongoDB{
		database: client.Database(databaseName),
	}
}

func (m *MongoDB) CreateShortURL(urlSchema model.UrlSchema) error {
	// Check the url doesn't exist
	filter := bson.D{{Key: "alias", Value: urlSchema.Alias}}
	err := m.database.Collection("url").FindOne(context.TODO(), filter).Decode(&urlSchema)
	if errors.Is(err, mongo.ErrNoDocuments) {
		_, err := m.database.Collection("url").InsertOne(context.TODO(), urlSchema)
		return err
	} else {
		return errors.New("Alias " + urlSchema.Alias + " already exist !")
	}
}

func (m *MongoDB) GetURL(alias string) (string, error) {
	filter := bson.D{{Key: "alias", Value: alias}}
	var urlSchema model.UrlSchema
	err := m.database.Collection("url").FindOne(context.TODO(), filter).Decode(&urlSchema)
	if err != nil {
		log.Error().Err(err)
	}
	return urlSchema.Url, err
}
