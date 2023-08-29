package storage

import (
	"context"
	"log"
	"goauth/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"errors"
	"time"
)

type UserStorage struct {
	db *mongo.Collection
}

func ConnectDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("authorizer")
}

func NewUserStorage(db *mongo.Database, collection string) *UserStorage {
	return &UserStorage{
		db: db.Collection(collection),
	}
}

func (r *UserStorage) CreateUser(ctx context.Context, user *entity.User) error {
	_, err := r.db.InsertOne(ctx, user)
	if err != nil {
		return errors.New("user already exists")
	}

	return nil
}

func (r *UserStorage) GetUser(ctx context.Context, username, password string) (*entity.User, error) {
	user := new(entity.User)

	if err := r.db.FindOne(ctx, bson.M{"_id": username, "password": password}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user not found")
		}

		return nil, err
	}

	return user, nil
}