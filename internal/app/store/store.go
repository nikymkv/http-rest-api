package store

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Store ...
type Store struct {
	config                   *Config
	db                       *mongo.Database
	userRepository           *UserRepository
	client                   *mongo.Client
	refreshSessionRepository *RefreshSessionRepository
}

// New ...
func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

// Open ...
func (s *Store) Open() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(s.config.DatabaseURL))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	db := client.Database(s.config.DatabaseName)

	s.db = db
	s.client = client

	return nil
}

// User ...
func (s *Store) User() *UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store:           s,
		usersCollection: s.db.Collection("my_users"),
	}

	return s.userRepository
}

// RefreshSession ...
func (s *Store) RefreshSession() *RefreshSessionRepository {
	if s.refreshSessionRepository != nil {
		return s.refreshSessionRepository
	}

	s.refreshSessionRepository = &RefreshSessionRepository{
		store:                    s,
		refreshSessionCollection: s.db.Collection("refresh_sessions"),
	}

	return s.refreshSessionRepository
}
