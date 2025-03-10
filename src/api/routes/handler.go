package routes

import (
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type Handler struct {
	db          *sqlx.DB
	oauthConfig *oauth2.Config
}

func NewHandler(db *sqlx.DB, oauthConfig *oauth2.Config) *Handler {
	return &Handler{
		db:          db,
		oauthConfig: oauthConfig,
	}
}

func (h *Handler) GetDB() *sqlx.DB {
	return h.db
}

func (h *Handler) GetOauthConfig() *oauth2.Config {
	return h.oauthConfig
}
