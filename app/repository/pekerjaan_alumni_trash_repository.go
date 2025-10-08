package repository

import (
	"database/sql"
	"hello-fiber/app/model"
	"log"
)

type PekerjaanAlumniTrashRepository interface {
	List() ([]model.PekerjaanAlumniTrashView, error)
	InsertIfNotExists(pekerjaanAlumniID int) error
	GetByTrashID(trashID int) (model.PekerjaanAlumniTrash, error)
	DeleteByTrashID(trashID int) error
	DeleteByPekerjaanID(pekerjaanAlumniID int) error
}

type pekerjaanAlumniTrashRepository struct {
	db *sql.DB
}

func NewPekerjaanAlumniTrashRepository(db *sql.DB) PekerjaanAlumniTrashRepository {
	return &pekerjaanAlumniTrashRepository{db: db}
}

func (r *pekerjaanAlumniTrashRepository) List() ([]model.PekerjaanAlumniTrashView, error) {
	query := `
		SELECT t.id AS trash_id, t.pekerjaan_alumni_id, t.status, t.created_at,
		       p.alumni_id, p.nama_perusahaan, p.posisi_jabatan, p.bidang_industri
		FROM pekerjaan_alumni_trash t
		JOIN pekerjaan_alumni p ON p.id = t.pekerjaan_alumni_id
		ORDER BY t.created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("Error list trash:", err)
		return nil, err
	}
	defer rows.Close()

	var out []model.PekerjaanAlumniTrashView
	for rows.Next() {
		var v model.PekerjaanAlumniTrashView
		if err := rows.Scan(
			&v.TrashID, &v.PekerjaanAlumniID, &v.Status, &v.CreatedAt,
			&v.AlumniID, &v.NamaPerusahaan, &v.PosisiJabatan, &v.BidangIndustri,
		); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func (r *pekerjaanAlumniTrashRepository) InsertIfNotExists(pekerjaanAlumniID int) error {
	// status selalu "hapus" untuk item masuk trash
	query := `
		INSERT INTO pekerjaan_alumni_trash (pekerjaan_alumni_id, status)
		VALUES ($1, 'hapus')
		ON CONFLICT (pekerjaan_alumni_id) DO NOTHING`
	_, err := r.db.Exec(query, pekerjaanAlumniID)
	if err != nil {
		log.Println("Error insert trash:", err)
	}
	return err
}

func (r *pekerjaanAlumniTrashRepository) GetByTrashID(trashID int) (model.PekerjaanAlumniTrash, error) {
	var t model.PekerjaanAlumniTrash
	query := `SELECT id, pekerjaan_alumni_id, status, created_at FROM pekerjaan_alumni_trash WHERE id = $1`
	err := r.db.QueryRow(query, trashID).Scan(&t.ID, &t.PekerjaanAlumniID, &t.Status, &t.CreatedAt)
	return t, err
}

func (r *pekerjaanAlumniTrashRepository) DeleteByTrashID(trashID int) error {
	_, err := r.db.Exec(`DELETE FROM pekerjaan_alumni_trash WHERE id = $1`, trashID)
	return err
}

func (r *pekerjaanAlumniTrashRepository) DeleteByPekerjaanID(pekerjaanAlumniID int) error {
	_, err := r.db.Exec(`DELETE FROM pekerjaan_alumni_trash WHERE pekerjaan_alumni_id = $1`, pekerjaanAlumniID)
	return err
}
