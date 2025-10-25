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
func (r *PekerjaanAlumniRepositoryMongo) CreatePekerjaanAlumni(req model.CreatePekerjaanAlumniRequest) (bson.ObjectID, error) {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if req.AlumniID.IsZero() {
		return bson.ObjectID{}, errors.New("alumni ID kosong/tidak valid")
	}

	doc := bson.M{
		"alumni_id":             req.AlumniID,
		"nama_perusahaan":       req.NamaPerusahaan,
		"posisi_jabatan":        req.PosisiJabatan,
		"bidang_industri":       req.BidangIndustri,
		"lokasi_kerja":          req.LokasiKerja,
		"gaji_range":            req.GajiRange,
		"tanggal_mulai_kerja":   req.TanggalMulaiKerja.Time,
		"tanggal_selesai_kerja": req.TanggalSelesaiKerja.Time,
		"status_pekerjaan":      req.StatusPekerjaan,
		"deskripsi_pekerjaan":   req.DeskripsiPekerjaan,
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

// repository/pekerjaan_alumni_repository.go

func (r *PekerjaanAlumniRepositoryMongo) GetAllPekerjaanAlumni() ([]model.PekerjaanAlumni, error) {
    collection := database.MongoDB.Collection("pekerjaan_alumni")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Ambil hanya dokumen yang belum dihapus DAN alumni_id benar-benar ObjectId
    filter := bson.M{
        "is_delete": "tidak",
        "alumni_id": bson.M{"$type": "objectId"},
    }

    cursor, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var out []model.PekerjaanAlumni
    for cursor.Next(ctx) {
        var row model.PekerjaanAlumni
        if err := cursor.Decode(&row); err != nil {
            // Optional: log error lalu lanjut, mirip pola di GetAllUsers
            // fmt.Printf("[WARNING] decode pekerjaan_alumni gagal: %v\n", err)
            continue
        }
        out = append(out, row)
    }
    if err := cursor.Err(); err != nil {
        return nil, err
    }
    return out, nil
}


// GetPekerjaanAlumniByID: pastikan dokumen punya alumni_id bertipe ObjectId juga.
func (r *PekerjaanAlumniRepositoryMongo) GetPekerjaanAlumniByID(id bson.ObjectID) (*model.PekerjaanAlumni, error) {
    collection := database.MongoDB.Collection("pekerjaan_alumni")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var pekerjaan model.PekerjaanAlumni
    filter := bson.M{
        "_id":       id,
        "is_delete": "tidak",
        // cegah decode error bila field salah tipe
        "alumni_id": bson.M{"$type": "objectId"},
    }
    if err := collection.FindOne(ctx, filter).Decode(&pekerjaan); err != nil {
        return nil, err
    }
    return &pekerjaan, nil
}

// GetPekerjaanAlumniByAlumniID: equality ke ObjectID sudah aman, tetap iterasi manual biar konsisten.
func (r *PekerjaanAlumniRepositoryMongo) GetPekerjaanAlumniByAlumniID(alumniID bson.ObjectID) ([]model.PekerjaanAlumni, error) {
    collection := database.MongoDB.Collection("pekerjaan_alumni")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{
        "alumni_id": alumniID,    // match ObjectID -> aman
        "is_delete": "tidak",
    }

    cur, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)

    var out []model.PekerjaanAlumni
    for cur.Next(ctx) {
        var row model.PekerjaanAlumni
        if err := cur.Decode(&row); err != nil {
            // TODO: log lalu lanjut
            continue
        }
        out = append(out, row)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return out, nil
}

// repository/pekerjaan_alumni_repository.go

func (r *PekerjaanAlumniRepositoryMongo) UpdatePekerjaanAlumni(id bson.ObjectID, req model.UpdatePekerjaanAlumniRequest) error {
    collection := database.MongoDB.Collection("pekerjaan_alumni")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    set := bson.M{}

    if req.AlumniID != nil {
        set["alumni_id"] = *req.AlumniID
    }
    if req.NamaPerusahaan != "" {
        set["nama_perusahaan"] = req.NamaPerusahaan
    }
    if req.PosisiJabatan != "" {
        set["posisi_jabatan"] = req.PosisiJabatan
    }
    if req.BidangIndustri != "" {
        set["bidang_industri"] = req.BidangIndustri
    }
    if req.LokasiKerja != "" {
        set["lokasi_kerja"] = req.LokasiKerja
    }
    if req.GajiRange != "" {
        set["gaji_range"] = req.GajiRange
    }
    if req.TanggalMulaiKerja != nil {
        set["tanggal_mulai_kerja"] = req.TanggalMulaiKerja.Time
    }
    if req.TanggalSelesaiKerja != nil {
        set["tanggal_selesai_kerja"] = req.TanggalSelesaiKerja.Time
    }
    if req.StatusPekerjaan != "" {
        set["status_pekerjaan"] = req.StatusPekerjaan
    }
    if req.DeskripsiPekerjaan != "" {
        set["deskripsi_pekerjaan"] = req.DeskripsiPekerjaan
    }

    if len(set) == 0 {
        return errors.New("tidak ada field yang diupdate")
    }
    set["updated_at"] = time.Now()

    update := bson.M{"$set": set}
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
func (r *PekerjaanAlumniRepositoryMongo) SoftDeletePekerjaanAlumni(id bson.ObjectID, isDelete bool) error {
	collection := database.MongoDB.Collection("pekerjaan_alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	deleteStatus := "tidak"
	if isDelete {
		deleteStatus = "hapus"
	}

	update := bson.M{
		"$set": bson.M{
			"is_delete":  deleteStatus,
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

// GetTrashedPekerjaanAlumni: versi "trash" juga difilter tipe-nya karena struct Trash punya AlumniID: ObjectID.
func (r *PekerjaanAlumniRepositoryMongo) GetTrashedPekerjaanAlumni() ([]model.Trash, error) {
    collection := database.MongoDB.Collection("pekerjaan_alumni")
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    filter := bson.M{
        "is_delete": "hapus",
        "alumni_id": bson.M{"$type": "objectId"},
    }

    cur, err := collection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    defer cur.Close(ctx)

    var out []model.Trash
    for cur.Next(ctx) {
        var row model.Trash
        if err := cur.Decode(&row); err != nil {
            // TODO: log lalu lanjut
            continue
        }
        out = append(out, row)
    }
    if err := cur.Err(); err != nil {
        return nil, err
    }
    return out, nil
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
