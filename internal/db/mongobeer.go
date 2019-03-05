package mongobeer

import (
	"context"

	"github.com/ctco-dev/go-api-template/internal/log"

	"github.com/ctco-dev/go-api-template/internal/beer"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type beerModel struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type repo struct {
	collection *mongo.Collection
}

// NewRepo returns a new mongodb beer repository
func NewRepo(ctx context.Context) beer.Repository {

	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")

	if err != nil {
		log.WithCtx(ctx).Panic(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.WithCtx(ctx).Panic(err)
	}

	log.WithCtx(ctx).Info("Connected to MongoDB!")
	return &repo{collection: client.Database("test").Collection("beer")}
}

func (r *repo) Read(ctx context.Context, id beer.ID) (*beer.Beer, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var res beerModel
	filter := bson.D{{"_id", oid}}
	err = r.collection.FindOne(ctx, filter).Decode(&res)
	if err != nil {
		return nil, err
	}

	return &beer.Beer{
		ID:   res.ID.Hex(),
		Name: res.Name,
	}, nil
}

func (r *repo) Write(ctx context.Context, beer beer.Beer) (beer.ID, error) {

	document := beerModel{Name: beer.Name}

	res, err := r.collection.InsertOne(ctx, beer)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *repo) Remove(ctx context.Context, id beer.ID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", oid}}
	_, err = r.collection.DeleteOne(ctx, filter)

	return err
}
