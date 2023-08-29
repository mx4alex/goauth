package storage

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"goauth/internal/entity"
	"goauth/internal/config"
	"log"
	"time"
)

type UserStorage struct {
	db *mongo.Collection
}

func ConnectDB(cfg config.MongoConfig) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.Url))
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

	return client.Database(cfg.Name)
}

func NewUserStorage(db *mongo.Database, cfg config.MongoConfig) *UserStorage {
	return &UserStorage{
		db: db.Collection(cfg.Collection),
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