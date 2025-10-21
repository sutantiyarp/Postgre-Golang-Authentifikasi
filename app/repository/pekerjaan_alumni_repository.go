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

type PekerjaanAlumniRepositoryMongo struct{}

func NewPekerjaanAlumniRepositoryMongo() *PekerjaanAlumniRepositoryMongo {
	return &PekerjaanAlumniRepositoryMongo{}
}

// CreatePekerjaanAlumni membuat pekerjaan alumni baru dan mengembalikan ID (bson.ObjectID)
func (r *PekerjaanAlumniRepositoryMongo) CreatePekerjaanAlumni(pekerjaan model.PekerjaanAlumni) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := bson.M{
		"alumni_id":             pekerjaan.AlumniID,
		"nama_perusahaan":       pekerjaan.NamaPerusahaan,
		"posisi_jabatan":        pekerjaan.PosisiJabatan,
		"bidang_industri":       pekerjaan.BidangIndustri,
		"lokasi_kerja":          pekerjaan.LokasiKerja,
		"gaji_range":            pekerjaan.GajiRange,
		"tanggal_mulai_kerja":   pekerjaan.TanggalMulaiKerja,
		"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
		"status_pekerjaan":      pekerjaan.StatusPekerjaan,
		"deskripsi_pekerjaan":   pekerjaan.DeskripsiPekerjaan,
		"is_delete":             "tidak",
		"created_at":            time.Now(),
		"updated_at":            time.Now(),
	}

	res, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return bson.ObjectID{}, err
	}

	insertedID, ok := res.InsertedID.(bson.ObjectID)
	if !ok {
		return bson.ObjectID{}, errors.New("failed to get inserted ID")
	}
	return insertedID, nil
}

// GetAllPekerjaanAlumni mengambil semua pekerjaan alumni yang belum dihapus
func (r *PekerjaanAlumniRepositoryMongo) GetAllPekerjaanAlumni() ([]model.PekerjaanAlumni, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_delete": "tidak"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

// GetPekerjaanAlumniByID mengambil pekerjaan alumni berdasarkan ID (bson.ObjectID)
func (r *PekerjaanAlumniRepositoryMongo) GetPekerjaanAlumniByID(id bson.ObjectID) (*model.PekerjaanAlumni, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pekerjaan model.PekerjaanAlumni
	if err := collection.FindOne(ctx, bson.M{"_id": id, "is_delete": "tidak"}).Decode(&pekerjaan); err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

// GetPekerjaanAlumniByAlumniID mengambil pekerjaan berdasarkan alumni_id (bson.ObjectID)
func (r *PekerjaanAlumniRepositoryMongo) GetPekerjaanAlumniByAlumniID(alumniID bson.ObjectID) ([]model.PekerjaanAlumni, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"alumni_id": alumniID, "is_delete": "tidak"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var pekerjaan []model.PekerjaanAlumni
	if err = cursor.All(ctx, &pekerjaan); err != nil {
		return nil, err
	}

	return pekerjaan, nil
}

// UpdatePekerjaanAlumni mengupdate pekerjaan alumni berdasarkan ID (bson.ObjectID)
func (r *PekerjaanAlumniRepositoryMongo) UpdatePekerjaanAlumni(id bson.ObjectID, pekerjaan model.PekerjaanAlumni) error {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":       pekerjaan.NamaPerusahaan,
			"posisi_jabatan":        pekerjaan.PosisiJabatan,
			"bidang_industri":       pekerjaan.BidangIndustri,
			"lokasi_kerja":          pekerjaan.LokasiKerja,
			"gaji_range":            pekerjaan.GajiRange,
			"tanggal_mulai_kerja":   pekerjaan.TanggalMulaiKerja,
			"tanggal_selesai_kerja": pekerjaan.TanggalSelesaiKerja,
			"status_pekerjaan":      pekerjaan.StatusPekerjaan,
			"deskripsi_pekerjaan":   pekerjaan.DeskripsiPekerjaan,
			"updated_at":            time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": id, "is_delete": "tidak"}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("pekerjaan alumni not found")
	}
	return nil
}

// SoftDeletePekerjaanAlumni soft delete pekerjaan alumni (set is_delete)
func (r *PekerjaanAlumniRepositoryMongo) SoftDeletePekerjaanAlumni(id bson.ObjectID, isDelete string) error {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_delete":  isDelete,
			"updated_at": time.Now(),
		},
	}

	res, err := collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("pekerjaan alumni not found")
	}
	return nil
}

// GetTrashedPekerjaanAlumni mengambil semua pekerjaan yang di-trash
func (r *PekerjaanAlumniRepositoryMongo) GetTrashedPekerjaanAlumni() ([]model.Trash, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"is_delete": "hapus"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var trash []model.Trash
	if err = cursor.All(ctx, &trash); err != nil {
		return nil, err
	}

	return trash, nil
}

// HardDeleteTrashedPekerjaanAlumni hard delete pekerjaan yang di-trash
func (r *PekerjaanAlumniRepositoryMongo) HardDeleteTrashedPekerjaanAlumni(id bson.ObjectID) (*model.PekerjaanAlumni, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pekerjaan model.PekerjaanAlumni
	if err := collection.FindOne(ctx, bson.M{"_id": id, "is_delete": "hapus"}).Decode(&pekerjaan); err != nil {
		return nil, err
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

// RestoreTrashedPekerjaanAlumni restore pekerjaan dari trash dan kembalikan dokumen terbaru
func (r *PekerjaanAlumniRepositoryMongo) RestoreTrashedPekerjaanAlumni(id bson.ObjectID) (*model.PekerjaanAlumni, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"is_delete":  "tidak",
			"updated_at": time.Now(),
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var pekerjaan model.PekerjaanAlumni
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": id, "is_delete": "hapus"}, update, opts).Decode(&pekerjaan); err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}
