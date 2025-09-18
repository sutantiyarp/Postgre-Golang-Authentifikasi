package repository

import (
	"hello-fiber/app/model"
	"database/sql"
	"log"
	"time"
)

// PekerjaanAlumniRepository interface untuk operasional pekerjaan alumni
type PekerjaanAlumniRepository interface {
	GetAll() ([]model.PekerjaanAlumni, error)
	GetByID(id int) (model.PekerjaanAlumni, error)
	GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error)
	Create(pekerjaan model.CreatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error)
	Update(id int, pekerjaan model.UpdatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error)
	Delete(id int) error
}

type pekerjaanAlumniRepository struct {
	db *sql.DB
}

// NewPekerjaanAlumniRepository membuat instance baru PekerjaanAlumniRepository
func NewPekerjaanAlumniRepository(db *sql.DB) PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{db: db}
}

// GetAll untuk mengambil semua data pekerjaan alumni
func (r *pekerjaanAlumniRepository) GetAll() ([]model.PekerjaanAlumni, error) {
	sqlStatement := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni 
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		log.Println("Error querying pekerjaan alumni:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		var tanggalSelesai sql.NullString
		
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan,
			&pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, 
			&pekerjaan.TanggalMulaiKerja, &tanggalSelesai, &pekerjaan.StatusPekerjaan,
			&pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning pekerjaan alumni:", err)
			return nil, err
		}
		
		if tanggalSelesai.Valid {
			pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
		}
		
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

// GetByID untuk mengambil data pekerjaan alumni berdasarkan ID
func (r *pekerjaanAlumniRepository) GetByID(id int) (model.PekerjaanAlumni, error) {
	sqlStatement := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni 
		WHERE id = $1`
	
	var pekerjaan model.PekerjaanAlumni
	var tanggalSelesai sql.NullString
	
	err := r.db.QueryRow(sqlStatement, id).Scan(
		&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan,
		&pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, 
		&pekerjaan.TanggalMulaiKerja, &tanggalSelesai, &pekerjaan.StatusPekerjaan,
		&pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		log.Println("Error finding pekerjaan alumni by ID:", err)
		return model.PekerjaanAlumni{}, err
	}
	
	if tanggalSelesai.Valid {
		pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
	}
	
	return pekerjaan, nil
}

// GetByAlumniID untuk mengambil semua pekerjaan berdasarkan alumni ID
func (r *pekerjaanAlumniRepository) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	sqlStatement := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at 
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 
		ORDER BY tanggal_mulai_kerja DESC`
	
	rows, err := r.db.Query(sqlStatement, alumniID)
	if err != nil {
		log.Println("Error querying pekerjaan alumni by alumni ID:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		var tanggalSelesai sql.NullString
		
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.AlumniID, &pekerjaan.NamaPerusahaan, &pekerjaan.PosisiJabatan,
			&pekerjaan.BidangIndustri, &pekerjaan.LokasiKerja, &pekerjaan.GajiRange, 
			&pekerjaan.TanggalMulaiKerja, &tanggalSelesai, &pekerjaan.StatusPekerjaan,
			&pekerjaan.DeskripsiPekerjaan, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning pekerjaan alumni:", err)
			return nil, err
		}
		
		if tanggalSelesai.Valid {
			pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
		}
		
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

// Create untuk menambah pekerjaan alumni baru
func (r *pekerjaanAlumniRepository) Create(req model.CreatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error) {
	sqlStatement := `
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		                             lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		                             status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) 
		RETURNING id, created_at, updated_at`
	
	var pekerjaan model.PekerjaanAlumni
	now := time.Now()
	
	var tanggalSelesai interface{}
	if req.TanggalSelesaiKerja != "" {
		tanggalSelesai = req.TanggalSelesaiKerja
	} else {
		tanggalSelesai = nil
	}
	
	err := r.db.QueryRow(
		sqlStatement, req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, tanggalSelesai,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, now, now,
	).Scan(&pekerjaan.ID, &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt)
	
	if err != nil {
		log.Println("Error inserting pekerjaan alumni:", err)
		return model.PekerjaanAlumni{}, err
	}
	
	// Set the other fields
	pekerjaan.AlumniID = req.AlumniID
	pekerjaan.NamaPerusahaan = req.NamaPerusahaan
	pekerjaan.PosisiJabatan = req.PosisiJabatan
	pekerjaan.BidangIndustri = req.BidangIndustri
	pekerjaan.LokasiKerja = req.LokasiKerja
	pekerjaan.GajiRange = req.GajiRange
	pekerjaan.TanggalMulaiKerja = req.TanggalMulaiKerja
	pekerjaan.TanggalSelesaiKerja = req.TanggalSelesaiKerja
	pekerjaan.StatusPekerjaan = req.StatusPekerjaan
	pekerjaan.DeskripsiPekerjaan = req.DeskripsiPekerjaan
	
	return pekerjaan, nil
}

// Update untuk mengupdate data pekerjaan alumni
func (r *pekerjaanAlumniRepository) Update(id int, req model.UpdatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error) {
	sqlStatement := `
		UPDATE pekerjaan_alumni 
		SET nama_perusahaan = $1, posisi_jabatan = $2, bidang_industri = $3, lokasi_kerja = $4, 
		    gaji_range = $5, tanggal_mulai_kerja = $6, tanggal_selesai_kerja = $7, 
		    status_pekerjaan = $8, deskripsi_pekerjaan = $9, updated_at = $10 
		WHERE id = $11`
	
	now := time.Now()
	
	var tanggalSelesai interface{}
	if req.TanggalSelesaiKerja != "" {
		tanggalSelesai = req.TanggalSelesaiKerja
	} else {
		tanggalSelesai = nil
	}
	
	result, err := r.db.Exec(
		sqlStatement, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, req.TanggalMulaiKerja, tanggalSelesai, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, now, id,
	)
	if err != nil {
		log.Println("Error updating pekerjaan alumni:", err)
		return model.PekerjaanAlumni{}, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.PekerjaanAlumni{}, sql.ErrNoRows
	}
	
	// Get updated pekerjaan alumni
	return r.GetByID(id)
}

// Delete untuk menghapus data pekerjaan alumni
func (r *pekerjaanAlumniRepository) Delete(id int) error {
	sqlStatement := `DELETE FROM pekerjaan_alumni WHERE id = $1`
	result, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		log.Println("Error deleting pekerjaan alumni:", err)
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}
