package repository

import (
	"hello-fiber/app/model"
	"database/sql"
	"log"
	"time"
)

// AlumniRepository interface untuk operasional alumni
type AlumniRepository interface {
	GetAll() ([]model.Alumni, error)
	GetByID(id int) (model.Alumni, error)
	Create(alumni model.CreateAlumniRequest) (model.Alumni, error)
	Update(id int, alumni model.UpdateAlumniRequest) (model.Alumni, error)
	Delete(id int) error
}

type alumniRepository struct {
	db *sql.DB
}

// NewAlumniRepository membuat instance baru AlumniRepository
func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

// GetAll untuk mengambil semua data alumni
func (r *alumniRepository) GetAll() ([]model.Alumni, error) {
	sqlStatement := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
		FROM alumni 
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		log.Println("Error querying alumni:", err)
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni
	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
			&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon,
			&alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning alumni:", err)
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}
	return alumniList, nil
}

// GetByID untuk mengambil data alumni berdasarkan ID
func (r *alumniRepository) GetByID(id int) (model.Alumni, error) {
	sqlStatement := `
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at 
		FROM alumni 
		WHERE id = $1`
	
	var alumni model.Alumni
	err := r.db.QueryRow(sqlStatement, id).Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
		&alumni.Angkatan, &alumni.TahunLulus, &alumni.Email, &alumni.NoTelepon,
		&alumni.Alamat, &alumni.CreatedAt, &alumni.UpdatedAt,
	)
	if err != nil {
		log.Println("Error finding alumni by ID:", err)
		return model.Alumni{}, err
	}
	return alumni, nil
}

// Create untuk menambah alumni baru
func (r *alumniRepository) Create(req model.CreateAlumniRequest) (model.Alumni, error) {
	sqlStatement := `
		INSERT INTO alumni (nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at, updated_at`
	
	var alumni model.Alumni
	now := time.Now()
	err := r.db.QueryRow(
		sqlStatement, req.NIM, req.Nama, req.Jurusan, req.Angkatan, 
		req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, now, now,
	).Scan(&alumni.ID, &alumni.CreatedAt, &alumni.UpdatedAt)
	
	if err != nil {
		log.Println("Error inserting alumni:", err)
		return model.Alumni{}, err
	}
	
	// Set the other fields
	alumni.NIM = req.NIM
	alumni.Nama = req.Nama
	alumni.Jurusan = req.Jurusan
	alumni.Angkatan = req.Angkatan
	alumni.TahunLulus = req.TahunLulus
	alumni.Email = req.Email
	alumni.NoTelepon = req.NoTelepon
	alumni.Alamat = req.Alamat
	
	return alumni, nil
}

// Update untuk mengupdate data alumni
func (r *alumniRepository) Update(id int, req model.UpdateAlumniRequest) (model.Alumni, error) {
	sqlStatement := `
		UPDATE alumni 
		SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, email = $5, no_telepon = $6, alamat = $7, updated_at = $8 
		WHERE id = $9`
	
	now := time.Now()
	result, err := r.db.Exec(
		sqlStatement, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus,
		req.Email, req.NoTelepon, req.Alamat, now, id,
	)
	if err != nil {
		log.Println("Error updating alumni:", err)
		return model.Alumni{}, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.Alumni{}, sql.ErrNoRows
	}
	
	// Get updated alumni
	return r.GetByID(id)
}

// Delete untuk menghapus data alumni
func (r *alumniRepository) Delete(id int) error {
	sqlStatement := `DELETE FROM alumni WHERE id = $1`
	result, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		log.Println("Error deleting alumni:", err)
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}
