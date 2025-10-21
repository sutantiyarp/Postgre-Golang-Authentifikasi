package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Alumni struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NIM        string        `bson:"nim" json:"nim"`
	Nama       string        `bson:"nama" json:"nama"`
	Jurusan    string        `bson:"jurusan" json:"jurusan"`
	Angkatan   int           `bson:"angkatan" json:"angkatan"`
	TahunLulus int           `bson:"tahun_lulus" json:"tahun_lulus"`
	Email      string        `bson:"email" json:"email"`
	NoTelepon  string        `bson:"no_telepon" json:"no_telepon"`
	Alamat     string        `bson:"alamat" json:"alamat"`
	CreatedAt  time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time     `bson:"updated_at" json:"updated_at"`
}

type CreateAlumniRequest struct {
	NIM        string `json:"nim"`
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

type UpdateAlumniRequest struct {
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}
