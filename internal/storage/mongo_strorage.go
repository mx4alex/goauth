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

func (r *UserStorage) CreateUser(ctx context.Context, user *entity.UserInput, refreshToken *entity.RefreshToken) error {
	newUser := bson.M{
		"username":      user.Username,
        "password":      user.Password,
        "refresh_token": refreshToken.Token,
		"expires_at":    refreshToken.ExpiresAt,
	}

	_, err := r.db.InsertOne(ctx, newUser)
	if err != nil {
		return errors.New("user already exists")
	}

	return nil
}

func (r *UserStorage) GetUser(ctx context.Context, username, password string) (string, error) {
	user := new(entity.UserDB)

	err := r.db.FindOne(ctx, bson.M{"username": username, "password": password}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", errors.New("user not found")
		}

		return "", err
	}

	return user.Username, nil
}

func (r *UserStorage) UpdateUser(ctx context.Context, username string, newToken *entity.RefreshToken) error {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"refresh_token": newToken.Token, "expires_at": newToken.ExpiresAt}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserStorage) Refresh(ctx context.Context, oldToken string, newToken *entity.RefreshToken) error {
	filter := bson.M{"refresh_token": oldToken}
	update := bson.M{"$set": bson.M{
		"refresh_token": newToken.Token,
		"expires_at": newToken.ExpiresAt,
	}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserStorage) GetUsername(ctx context.Context, refreshToken string) (string, time.Time, error) {
	user := new(entity.UserDB)

	err := r.db.FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", time.Time{}, errors.New("user not found")
		}

		return "", time.Time{}, err
	}

	return user.Username, user.ExpiresAt, nil
}