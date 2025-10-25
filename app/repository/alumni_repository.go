package repository

import (
	"context"
	"errors"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type AlumniRepositoryMongo struct{}

func NewAlumniRepositoryMongo() *AlumniRepositoryMongo {
	return &AlumniRepositoryMongo{}
}

// CreateAlumni membuat alumni baru dan mengembalikan ID (bson.ObjectID)
func (r *AlumniRepositoryMongo) CreateAlumni(req model.CreateAlumniRequest) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := bson.NewObjectID()
	doc := bson.M{
		"_id":         id,
		"nim":         req.NIM,
		"nama":        req.Nama,
		"jurusan":     req.Jurusan,
		"angkatan":    req.Angkatan,
		"tahun_lulus": req.TahunLulus,
		"email":       req.Email,
		"no_telepon":  req.NoTelepon,
		"alamat":      req.Alamat,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}

	_, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return bson.ObjectID{}, err
	}

	return id, nil
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

// GetAllAlumniWithPagination untuk mendukung pagination
func (r *AlumniRepositoryMongo) GetAllAlumniWithPagination(page, limit int64) ([]model.Alumni, int64, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hitung total dokumen
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Hitung skip
	skip := (page - 1) * limit

	opts := options.Find().SetSkip(skip).SetLimit(limit)

	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var alumni []model.Alumni
	if err = cursor.All(ctx, &alumni); err != nil {
		return nil, 0, err
	}

	return alumni, total, nil
}

// GetAlumniByID mengambil alumni berdasarkan ID (bson.ObjectID)
func (r *AlumniRepositoryMongo) GetAlumniByID(id bson.ObjectID) (*model.Alumni, error) {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var alumni model.Alumni
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&alumni); err != nil {
		return nil, err
	}

	return &alumni, nil
}

// UpdateAlumni mengupdate data alumni berdasarkan ID (bson.ObjectID)
func (r *AlumniRepositoryMongo) UpdateAlumni(id bson.ObjectID, req model.UpdateAlumniRequest) error {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateFields := bson.M{}
	
	if req.NIM != nil {
		updateFields["nim"] = *req.NIM
	}
	if req.Nama != nil {
		updateFields["nama"] = *req.Nama
	}
	if req.Jurusan != nil {
		updateFields["jurusan"] = *req.Jurusan
	}
	if req.Angkatan != nil {
		updateFields["angkatan"] = *req.Angkatan
	}
	if req.TahunLulus != nil {
		updateFields["tahun_lulus"] = *req.TahunLulus
	}
	if req.Email != nil {
		updateFields["email"] = *req.Email
	}
	if req.NoTelepon != nil {
		updateFields["no_telepon"] = *req.NoTelepon
	}
	if req.Alamat != nil {
		updateFields["alamat"] = *req.Alamat
	}
	
	updateFields["updated_at"] = time.Now()

	update := bson.M{
		"$set": updateFields,
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("alumni not found")
	}

	return nil
}

// DeleteAlumni menghapus alumni berdasarkan ID (bson.ObjectID)
func (r *AlumniRepositoryMongo) DeleteAlumni(id bson.ObjectID) error {
	collection := database.MongoDB.Collection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("alumni not found")
	}

	return nil
}