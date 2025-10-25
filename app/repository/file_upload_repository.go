package repository

import (
	"context"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type FileUploadRepository interface {
	Create(file *model.FileUpload) error
	FindAll() ([]model.FileUpload, error)
	FindByID(id bson.ObjectID) (*model.FileUpload, error)
	Delete(id bson.ObjectID) error
}

type fileUploadRepository struct {
	collection *mongo.Collection
}

func NewFileUploadRepository() FileUploadRepository {
	return &fileUploadRepository{
		collection: database.MongoDB.Collection("files"),
	}
}

func (r *fileUploadRepository) Create(file *model.FileUpload) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file.UploadedAt = time.Now()
	if file.ID.IsZero() {
		file.ID = bson.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, file)
	return err
}

func (r *fileUploadRepository) FindAll() ([]model.FileUpload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var files []model.FileUpload
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}

	return files, nil
}

func (r *fileUploadRepository) FindByID(id bson.ObjectID) (*model.FileUpload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var file model.FileUpload
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (r *fileUploadRepository) Delete(id bson.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
