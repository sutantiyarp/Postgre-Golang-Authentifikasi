package repository

import (
	"hello-fiber/app/model"
	"database/sql"
	"log"
	"time"
	"fmt"
)

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

func NewAlumniRepository(db *sql.DB) AlumniRepository {
	return &alumniRepository{db: db}
}

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

	return r.GetByID(id)
}

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

func GetAlumniWithPagination(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	validSortColumns := map[string]bool{
		"id": true, "nim": true, "nama": true, "jurusan": true, 
		"angkatan": true, "tahun_lulus": true, "email": true, "created_at": true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "id"
	}

	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var alumni []model.Alumni
	for rows.Next() {
		var a model.Alumni
		err := rows.Scan(
			&a.ID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumni = append(alumni, a)
	}

	return alumni, nil
}

func CountAlumni(db *sql.DB, search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM alumni WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1 OR email ILIKE $1`
	err := db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}
