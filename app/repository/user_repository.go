package repository

import (
	"context"
	"errors"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryMongo struct{}

func NewUserRepositoryMongo() *UserRepositoryMongo {
	return &UserRepositoryMongo{}
}

func (r *UserRepositoryMongo) CreateUser(user model.User) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password sebelum disimpan
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return bson.ObjectID{}, err
	}

	user.ID = bson.NewObjectID()
	user.Password = string(hashed)
	user.CreatedAt = time.Now()

	// Cek jika AlumniID kosong, set ke nil
	if user.AlumniID == nil {
		// AlumniID di-set menjadi nil jika kosong
		user.AlumniID = nil
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.ObjectID{}, errors.New("email sudah terdaftar")
		}
		return bson.ObjectID{}, err
	}

	return result.InsertedID.(bson.ObjectID), nil
}


// GetUserByEmail mengambil user berdasarkan email (untuk login)
func (r *UserRepositoryMongo) GetUserByEmail(email string) (*model.User, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	if err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByID mengambil user berdasarkan ID
func (r *UserRepositoryMongo) GetUserByID(id bson.ObjectID) (*model.User, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("user tidak ditemukan")
		}
		return nil, err
	}
	return &user, nil
}
