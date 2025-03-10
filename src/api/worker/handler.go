package worker

import "github.com/jmoiron/sqlx"

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetDB() *sqlx.DB {
	return h.db
}
