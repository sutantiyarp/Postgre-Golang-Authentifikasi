package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hello-fiber/app/model"
	"hello-fiber/database"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type UserRepositoryMongo struct{}

func NewUserRepositoryMongo() *UserRepositoryMongo {
	return &UserRepositoryMongo{}
}

func (r *UserRepositoryMongo) Register(req model.RegisterRequest) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password sebelum disimpan
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return bson.ObjectID{}, err
	}

	user := model.User{
		ID:        bson.NewObjectID(),
		Username:  strings.TrimSpace(req.Username),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  string(hashed),
		RoleID:    req.RoleID,
		CreatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.ObjectID{}, errors.New("email atau username sudah terdaftar")
		}
		return bson.ObjectID{}, err
	}

	return result.InsertedID.(bson.ObjectID), nil
}

func (r *UserRepositoryMongo) CreateUser(req model.CreateUserRequest) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hash password sebelum disimpan
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return bson.ObjectID{}, err
	}

	user := model.User{
		ID:        bson.NewObjectID(),
		Username:  strings.TrimSpace(req.Username),
		Email:     strings.ToLower(strings.TrimSpace(req.Email)),
		Password:  string(hashed),
		RoleID:    req.RoleID,
		AlumniID:  req.AlumniID,
		CreatedAt: time.Now(),
	}

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return bson.ObjectID{}, errors.New("email atau username sudah terdaftar")
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

func (r *UserRepositoryMongo) GetUserByUsername(username string) (*model.User, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user model.User
	if err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Username belum ada, return nil (tidak error)
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryMongo) GetAllUsers(page, limit int64) ([]model.User, int64, error) {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Hitung total users
	total, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	// Pagination
	opts := options.Find().SetSkip((page - 1) * limit).SetLimit(limit)
	cursor, err := collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			// Log error but continue processing other documents
			fmt.Printf("[WARNING] Failed to decode user document: %v\n", err)
			continue
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepositoryMongo) UpdateUser(id bson.ObjectID, req model.UpdateUserRequest) error {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{}
	if req.Username != "" {
		update["username"] = strings.TrimSpace(req.Username)
	}
	if req.Email != "" {
		update["email"] = strings.ToLower(strings.TrimSpace(req.Email))
	}
	if req.Password != "" {
		// Hash password sebelum disimpan
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		update["password"] = string(hashed)
	}
	if req.RoleID != bson.NilObjectID {
		update["role_id"] = req.RoleID
	}
	if req.AlumniID != nil && !req.AlumniID.IsZero() {
		update["alumni_id"] = req.AlumniID
	}

	if len(update) == 0 {
		return errors.New("tidak ada data yang diupdate")
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("email atau username sudah terdaftar")
		}
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

func (r *UserRepositoryMongo) DeleteUser(id bson.ObjectID) error {
	collection := database.MongoDB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("user tidak ditemukan")
	}

	return nil
}

func (r *UserRepositoryMongo) GetRoleByID(id bson.ObjectID) (*model.Role, error) {
	collection := database.MongoDB.Collection("roles")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var role model.Role
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&role); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role tidak ditemukan")
		}
		return nil, err
	}
	return &role, nil
}