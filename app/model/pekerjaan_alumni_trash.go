package model

import "time"

const (
	TrashStatusHapus = "hapus"
)

type PekerjaanAlumniTrash struct {
	ID                  int       `json:"id"`
	PekerjaanAlumniID   int       `json:"pekerjaan_alumni_id"`
	Status              string    `json:"status"` // expected "hapus"
	CreatedAt           time.Time `json:"created_at"`
}

type PekerjaanAlumniTrashView struct {
    TrashID         int       `json:"trash_id"`
    PekerjaanAlumniID int     `json:"pekerjaan_alumni_id"`
    Status          string    `json:"status"`
    CreatedAt       time.Time `json:"created_at"`
    AlumniID        int       `json:"alumni_id"`
    NamaPerusahaan  string    `json:"nama_perusahaan"`
    PosisiJabatan   string    `json:"posisi_jabatan"`
    BidangIndustri  string    `json:"bidang_industri"`
}

type PekerjaanAlumniTrashListResponse struct {
	Data []PekerjaanAlumniTrashView `json:"data"`
}
