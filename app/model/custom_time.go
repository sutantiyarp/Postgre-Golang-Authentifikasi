package model

import (
	"time"
)

// CustomTime adalah custom type untuk parsing tanggal dengan format fleksibel
type CustomTime struct {
	time.Time
}

// UnmarshalJSON mengimplementasikan custom JSON unmarshaling untuk CustomTime
// Mendukung format: "2024-10-05" (date only) dan "2024-10-05T15:04:05Z07:00" (datetime)
func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	
	// Hapus quotes dari JSON string
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	// Coba parse format date saja (yyyy-mm-dd)
	t, err := time.Parse("2006-01-02", s)
	if err == nil {
		ct.Time = t
		return nil
	}

	// Coba parse format RFC3339 (datetime dengan timezone)
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		ct.Time = t
		return nil
	}

	// Jika kedua format gagal, return error
	return err
}
