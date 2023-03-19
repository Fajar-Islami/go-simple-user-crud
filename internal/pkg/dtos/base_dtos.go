package dtos

import "time"

type (
	DtosModel struct {
		ID        int64     `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	Pagination struct {
		Page       int `json:"page"`
		Rows       int `json:"rows"`
		TotalRows  int `json:"total_rows"`
		TotalPages int `json:"total_pages"`
	}
)
