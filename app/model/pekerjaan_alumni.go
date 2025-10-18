package model

import (
	"time"
)

type PekerjaanAlumni struct {
    ID                 int       `json:"id"`
    AlumniID           int       `json:"alumni_id"`
    NamaPerusahaan    string    `json:"nama_perusahaan"`
    PosisiJabatan     string    `json:"posisi_jabatan"`
    BidangIndustri    string    `json:"bidang_industri"`
    LokasiKerja       string    `json:"lokasi_kerja"`
    GajiRange         string    `json:"gaji_range"`
    TanggalMulaiKerja string    `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja string    `json:"tanggal_selesai_kerja"`
    StatusPekerjaan   string    `json:"status_pekerjaan"`
    DeskripsiPekerjaan string    `json:"deskripsi_pekerjaan"`
    CreatedAt         time.Time `json:"created_at"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type CreatePekerjaanAlumniRequest struct {
    AlumniID          int    `json:"alumni_id"`
    NamaPerusahaan   string `json:"nama_perusahaan"`
    PosisiJabatan    string `json:"posisi_jabatan"`
    BidangIndustri   string `json:"bidang_industri"`
    LokasiKerja      string `json:"lokasi_kerja"`
    GajiRange        string `json:"gaji_range"`
    TanggalMulaiKerja string `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
    StatusPekerjaan  string `json:"status_pekerjaan"`
    DeskripsiPekerjaan string `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanAlumniRequest struct {
    NamaPerusahaan    string `json:"nama_perusahaan"`
    PosisiJabatan     string `json:"posisi_jabatan"`
    BidangIndustri    string `json:"bidang_industri"`
    LokasiKerja       string `json:"lokasi_kerja"`
    GajiRange         string `json:"gaji_range"`
    TanggalMulaiKerja string `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
    StatusPekerjaan   string `json:"status_pekerjaan"`
    DeskripsiPekerjaan string `json:"deskripsi_pekerjaan"`
    UpdatedAt         time.Time `json:"updated_at"`
}

type UpdatePekerjaanAlumniSoftDelete struct {
    IsDelete    string `json:"is_delete"`
}

type Trash struct {
    ID          int    `json:"id"`
    AlumniID    int    `json:"alumni_id"`
    IsDelete    string `json:"is_delete"`
    UpdatedAt         time.Time `json:"updated_at"`
}