package mongobeer

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/ctco-dev/go-api-template/internal/beer"
	"github.com/ctco-dev/go-api-template/internal/log"
)

type beerModel struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}

type repo struct {
	collection *mongo.Collection
}

// NewRepo returns a new mongodb beer repository
func NewRepo(ctx context.Context, host string, db string, collection string) beer.Repository {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.WithCtx(ctx).Panic(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.WithCtx(ctx).Panic(err)
	}

	err = client.Ping(ctx, nil)

	if err != nil {
		log.WithCtx(ctx).Panic(err)
	}

	log.WithCtx(ctx).Info("Connected to MongoDB!")
	return &repo{collection: client.Database(db).Collection(collection)}
}

func (r *repo) Read(ctx context.Context, id beer.ID) (*beer.Beer, error) {
	log.WithCtx(ctx).Info("Reading a beer")
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

func (r *repo) ReadAll(ctx context.Context) ([]*beer.Beer, error) {
	log.WithCtx(ctx).Info("Reading all beers")
	cur, err := r.collection.Find(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	var res []*beer.Beer

	for cur.Next(ctx) {
		var elem beerModel
		err = cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		res = append(res, &beer.Beer{
			ID:   elem.ID.Hex(),
			Name: elem.Name,
		})
	}

	return res, nil
}

func (r *repo) Write(ctx context.Context, beer beer.Beer) (beer.ID, error) {
	log.WithCtx(ctx).Info("Writing a beer")
	document := bson.D{
		{"name", beer.Name},
	}

	res, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *repo) Remove(ctx context.Context, id beer.ID) error {
	log.WithCtx(ctx).Info("Removing a beer")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", oid}}
	_, err = r.collection.DeleteOne(ctx, filter)

	return err
}
