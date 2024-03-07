package kb

import (
	"context"
	"errors"
	"fmt"
	"ruby/pkg/user"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type kb struct {
	collection *mongo.Collection
}

func (d *kb) Create(ctx context.Context, arg user.CreateSessionParams) (string, error) {
  logrus.Debug("create session")
	result, err := d.collection.InsertOne(ctx, arg)
	if err != nil {
		return "", fmt.Errorf("failed to create session due error: %v", err)
	}
	logrus.Debug("convert InsertedID to ObjectID")
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}
	logrus.Trace(arg)
	return "", fmt.Errorf("failed to objectid to hex, probably oid: %s", oid)
}

 func (d *kb) Delete(ctx context.Context, id uuid.UUID) error {
 	filter := bson.M{"sessionid": id}
 	result, err := d.collection.DeleteOne(ctx, filter)
 	if err != nil {
 		return fmt.Errorf("failed execute query. error: %v", err)
 	}
 	if result.DeletedCount == 0 {
 		return nil
 	}
	logrus.Tracef("Deleted %d documents", result.DeletedCount)
 	return nil
 }

func (d *kb) FindOne(ctx context.Context, id uuid.UUID) (s user.CreateSessionParams, err error) {
 	filter := bson.M{"sessionid": id}

 	result := d.collection.FindOne(ctx, filter)
 	if result.Err() != nil {
 		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
 			return s, nil
 		}
 		return s, fmt.Errorf("failed to find one user by id: %s due to error: %v", id, err)
 	}
 	if err = result.Decode(&s); err != nil {
 		return s, fmt.Errorf("failed to decode user(id:%s) from DB due to error: %v", id, err)
 	}
 	return s, nil
 }

func NewStorage(database *mongo.Database, collection string) user.Storage {
	return &kb{
		collection: database.Collection(collection),
	}
}
