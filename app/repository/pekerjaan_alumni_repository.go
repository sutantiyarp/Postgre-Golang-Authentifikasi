package repository

import (
	"hello-fiber/app/model"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type PekerjaanAlumniRepository interface {
	GetAll() ([]model.PekerjaanAlumni, error)
	GetByID(id int) (model.PekerjaanAlumni, error)
	GetTrashByID(id int) (model.Trash, error)
	GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error)
	GetPekerjaanAlumniWithPagination(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error)
	Create(pekerjaan model.CreatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error)
	Update(id int, pekerjaan model.UpdatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error)
	UpdateUser(id int, pekerjaan model.UpdatePekerjaanAlumniSoftDelete) (model.Trash, error)
	UpdateAdmin(id int, pekerjaan model.UpdatePekerjaanAlumniSoftDelete) (model.Trash, error)
	Delete(id int) error
	CountPekerjaanAlumni(search string) (int, error)
	GetTrashed() ([]model.Trash, error)
	HardDeleteIfTrashed(id int) error
	RestoreIfTrashed(id int) (model.Trash, error)
}

type pekerjaanAlumniRepository struct {
	db *sql.DB
}

func NewPekerjaanAlumniRepository(db *sql.DB) PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{db: db}
}

func (r *pekerjaanAlumniRepository) GetAll() ([]model.PekerjaanAlumni, error) {
	sqlStatement := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni 
		ORDER BY created_at DESC`
	
	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		log.Println("Error men-query pekerjaan alumni:", err)
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
			log.Println("Error men-scan pekerjaan alumni:", err)
			return nil, err
		}
		
		if tanggalSelesai.Valid {
			pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
		}
		
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}

func (r *pekerjaanAlumniRepository) GetByID(id int) (model.PekerjaanAlumni, error) {
	sqlStatement := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
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
		log.Println("Error menemukan pekerjaan alumni by ID:", err)
		return model.PekerjaanAlumni{}, err
	}
	
	if tanggalSelesai.Valid {
		pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
	}
	
	return pekerjaan, nil
}


func (r *pekerjaanAlumniRepository) GetTrashByID(id int) (model.Trash, error) {
	sqlStatement := `SELECT id, alumni_id, is_delete 
		FROM pekerjaan_alumni 
		WHERE id = $1 AND is_delete = 'hapus'`
	
	var trash model.Trash
	
	err := r.db.QueryRow(sqlStatement, id).Scan(
		&trash.ID, &trash.AlumniID, &trash.IsDelete,
	)
	if err != nil {
		log.Println("Error menemukan pekerjaan alumni by ID:", err)
		return model.Trash{}, err
	}
	
	return trash, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniID(alumniID int) ([]model.PekerjaanAlumni, error) {
	sqlStatement := `SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni 
		WHERE alumni_id = $1 
		ORDER BY tanggal_mulai_kerja DESC`
	
	rows, err := r.db.Query(sqlStatement, alumniID)
	if err != nil {
		log.Println("Error men-query pekerjaan alumni by alumni ID:", err)
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
			log.Println("Error men-scan pekerjaan alumni:", err)
			return nil, err
		}
		
		if tanggalSelesai.Valid {
			pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
		}
		
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}
	return pekerjaanList, nil
}


func (r *pekerjaanAlumniRepository) GetPekerjaanAlumniWithPagination(search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	validSortColumns := map[string]bool{
		"id": true, "nama_perusahaan": true, "posisi_jabatan": true, "bidang_industri": true,
		"lokasi_kerja": true, "tanggal_mulai_kerja": true, "status_pekerjaan": true, "created_at": true,
	}
	if !validSortColumns[sortBy] {
		sortBy = "id"
	}

	query := fmt.Sprintf(`SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		       deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := r.db.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
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
			return nil, err
		}
		
		if tanggalSelesai.Valid {
			pekerjaan.TanggalSelesaiKerja = tanggalSelesai.String
		}
		
		pekerjaanList = append(pekerjaanList, pekerjaan)
	}

	return pekerjaanList, nil
}

func (r *pekerjaanAlumniRepository) GetTrashed() ([]model.Trash, error) {
	sqlStatement := `SELECT id, alumni_id, is_delete, updated_at 
		FROM pekerjaan_alumni 
		WHERE is_delete = 'hapus'
		ORDER BY updated_at DESC`

	rows, err := r.db.Query(sqlStatement)
	if err != nil {
		log.Println("Error memanggil trash pekerjaan_alumni:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Trash
	for rows.Next() {
		var trash model.Trash
		if err := rows.Scan(
			&trash.ID, &trash.AlumniID, &trash.IsDelete, &trash.UpdatedAt,
		); err != nil {
			log.Println("Error men-scan trash pekerjaan_alumni:", err)
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, trash)
	}
	return pekerjaanList, nil
}

func (r *pekerjaanAlumniRepository) Create(req model.CreatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error) {
	sqlStatement := `INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
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
	// pekerjaan.IsDelete = req.IsDelete
	
	return pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) Update(id int, req model.UpdatePekerjaanAlumniRequest) (model.PekerjaanAlumni, error) {
	sqlStatement := `UPDATE pekerjaan_alumni 
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
	
	return r.GetByID(id)
}

func (r *pekerjaanAlumniRepository) UpdateUser(id int, req model.UpdatePekerjaanAlumniSoftDelete) (model.Trash, error) {
	sqlStatement := `UPDATE pekerjaan_alumni 
		SET is_delete = $1, updated_at = $2
		WHERE id = $3`
	
	now := time.Now()
	
	result, err := r.db.Exec(
		sqlStatement, req.IsDelete, now, id,
	)
	if err != nil {
		log.Println("Error memperbarui pekerjaan alumni user to user:", err)
		return model.Trash{}, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.Trash{}, sql.ErrNoRows
	}
	
	return r.GetTrashByID(id)
}

func (r *pekerjaanAlumniRepository) UpdateAdmin(id int, req model.UpdatePekerjaanAlumniSoftDelete) (model.Trash, error) {
	sqlStatement := `UPDATE pekerjaan_alumni 
		SET is_delete = $1, updated_at = $2
		WHERE id = $3`
	
	now := time.Now()
	
	result, err := r.db.Exec(
		sqlStatement, req.IsDelete, now, id,
	)
	if err != nil {
		log.Println("Error memperbarui pekerjaan alumni admin to user:", err)
		return model.Trash{}, err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.Trash{}, sql.ErrNoRows
	}
	
	return r.GetTrashByID(id)
}

func (r *pekerjaanAlumniRepository) Delete(id int) error {
	sqlStatement := `UPDATE pekerjaan_alumni SET is_delete = 'hapus' WHERE id = $1`
	result, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		log.Println("Error menghapus pekerjaan alumni:", err)
		return err
	}
	
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	
	return nil
}

func (r *pekerjaanAlumniRepository) CountPekerjaanAlumni(search string) (int, error) {
	var total int
	countQuery := `SELECT COUNT(*) FROM pekerjaan_alumni WHERE nama_perusahaan ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1 OR lokasi_kerja ILIKE $1`
	err := r.db.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func (r *pekerjaanAlumniRepository) HardDeleteIfTrashed(id int) error {
	sqlStatement := `DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_delete = 'hapus'`
	result, err := r.db.Exec(sqlStatement, id)
	if err != nil {
		log.Println("Error hard deleting pekerjaan_alumni:", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *pekerjaanAlumniRepository) RestoreIfTrashed(id int) (model.Trash, error) {
	sqlStatement := `UPDATE pekerjaan_alumni
		SET is_delete = 'tidak', updated_at = $2
		WHERE id = $1 AND is_delete = 'hapus'`
	now := time.Now()

	result, err := r.db.Exec(sqlStatement, id, now)
	if err != nil {
		log.Println("Error untuk mengembalikan pekerjaan_alumni:", err)
		return model.Trash{}, err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return model.Trash{}, sql.ErrNoRows
	}
	
	sqlSelect := `SELECT id, alumni_id, is_delete, updated_at 
		FROM pekerjaan_alumni 
		WHERE id = $1`
	
	var trash model.Trash
	err = r.db.QueryRow(sqlSelect, id).Scan(
		&trash.ID, &trash.AlumniID, &trash.IsDelete, &trash.UpdatedAt,
	)
	if err != nil {
		log.Println("Error mengambil restored pekerjaan_alumni:", err)
		return model.Trash{}, err
	}
	
	return trash, nil
}
