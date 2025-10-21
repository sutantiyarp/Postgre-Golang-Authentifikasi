package repository

import (
	"context"
	"errors"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type AlumniRepositoryMongo struct{}

func NewAlumniRepositoryMongo() *AlumniRepositoryMongo {
	return &AlumniRepositoryMongo{}
}

// CreateAlumni membuat alumni baru dan mengembalikan ID (hex string)
func (r *AlumniRepositoryMongo) CreateAlumni(alumni model.Alumni) (string, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := bson.NewObjectID()
	doc := bson.M{
		"_id":         id,
		"nim":         alumni.NIM,
		"nama":        alumni.Nama,
		"jurusan":     alumni.Jurusan,
		"angkatan":    alumni.Angkatan,
		"tahun_lulus": alumni.TahunLulus,
		"email":       alumni.Email,
		"no_telepon":  alumni.NoTelepon,
		"alamat":      alumni.Alamat,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

// GetAllAlumni mengambil semua alumni
func (r *AlumniRepositoryMongo) GetAllAlumni() ([]model.Alumni, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var alumni []model.Alumni
	if err = cursor.All(ctx, &alumni); err != nil {
		return nil, err
	}

	return alumni, nil
}

// GetAlumniByID mengambil alumni berdasarkan ID (hex string)
func (r *AlumniRepositoryMongo) GetAlumniByID(id string) (*model.Alumni, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var alumni model.Alumni
	if err := collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&alumni); err != nil {
		return nil, err
	}

	return &alumni, nil
}

// UpdateAlumni mengupdate data alumni berdasarkan ID (hex string)
func (r *AlumniRepositoryMongo) UpdateAlumni(id string, alumni model.Alumni) error {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	update := bson.M{
		"$set": bson.M{
			"nim":         alumni.NIM,
			"nama":        alumni.Nama,
			"jurusan":     alumni.Jurusan,
			"angkatan":    alumni.Angkatan,
			"tahun_lulus": alumni.TahunLulus,
			"email":       alumni.Email,
			"no_telepon":  alumni.NoTelepon,
			"alamat":      alumni.Alamat,
			"updated_at":  time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("alumni not found")
	}

	return nil
}

// DeleteAlumni menghapus alumni berdasarkan ID (hex string)
func (r *AlumniRepositoryMongo) DeleteAlumni(id string) error {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid ID format")
	}

	res, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("alumni not found")
	}

	return nil
}
