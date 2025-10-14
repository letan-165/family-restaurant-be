package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BaseRepository[T any] struct {
	Collection *mongo.Collection
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]T, error) {
	cursor, err := r.Collection.Find(ctx, filter, opts)
	fmt.Printf("%+v\n", opts)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *BaseRepository[T]) Insert(ctx context.Context, doc T) (*primitive.ObjectID, error) {
	res, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}

	return &oid, nil
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var result T
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *BaseRepository[T]) Update(ctx context.Context, id any, update bson.M) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id any) (*mongo.DeleteResult, error) {
	return r.Collection.DeleteOne(ctx, bson.M{"_id": id})
}

func (r *BaseRepository[T]) Count(ctx context.Context, filter bson.M) (int64, error) {
	return r.Collection.CountDocuments(ctx, filter)
}
