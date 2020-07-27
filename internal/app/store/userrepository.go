package store

import (
	"context"

	"github.com/http-rest-api/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository ...
type UserRepository struct {
	store           *Store
	usersCollection *mongo.Collection
}

// Create ...
func (r *UserRepository) Create(ctx context.Context, u *model.User) (string, error) {
	err := u.BeforeCreate()
	if err != nil {
		return "", nil
	}

	result, err := r.usersCollection.InsertOne(ctx, bson.D{
		{Key: "email", Value: u.Email},
		{Key: "password", Value: u.Password},
	})
	if err != nil {
		return "", err
	}

	objectID := result.InsertedID.(primitive.ObjectID)

	return objectID.Hex(), nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	filter := bson.M{"email": email}
	err := r.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByID ...
func (r *UserRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
	user := &model.User{}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": userObjectID}
	err = r.usersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
