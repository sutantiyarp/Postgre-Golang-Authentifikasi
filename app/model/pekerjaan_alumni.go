package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type PekerjaanAlumni struct {
	ID                  bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID            bson.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string        `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string        `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string        `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string        `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           string        `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   time.Time     `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time     `bson:"tanggal_selesai_kerja" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string        `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string        `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time     `bson:"updated_at" json:"updated_at"`
}

type PekerjaanAlumniRequest struct {
	AlumniID            bson.ObjectID `json:"alumni_id"`
	NamaPerusahaan      string        `json:"nama_perusahaan"`
	PosisiJabatan       string        `json:"posisi_jabatan"`
	BidangIndustri      string        `json:"bidang_industri"`
	LokasiKerja         string        `json:"lokasi_kerja"`
	GajiRange           string        `json:"gaji_range"`
	TanggalMulaiKerja   time.Time     `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja time.Time     `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string        `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string        `json:"deskripsi_pekerjaan"`
}

type CreatePekerjaanAlumniRequest struct {
	AlumniID            bson.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string        `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string        `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string        `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string        `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           string        `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   CustomTime    `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja CustomTime    `bson:"tanggal_selesai_kerja" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string        `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string        `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanAlumniRequest struct {
	AlumniID            *bson.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string         `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string         `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string         `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string         `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           string         `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   *CustomTime    `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *CustomTime    `bson:"tanggal_selesai_kerja" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string         `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string         `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanAlumniSoftDelete struct {
	IsDelete string `bson:"is_delete" json:"is_delete"`
}

type Trash struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID  bson.ObjectID `bson:"alumni_id" json:"alumni_id"`
	IsDelete  string        `bson:"is_delete" json:"is_delete"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}