package models

import (
	"github.com/tmm6907/dashboard/utils"
)

type Feed struct {
	ID            uint            `db:"id" json:"id"`
	FeedID        utils.UUID      `db:"feed_id" json:"feedId"`
	Title         string          `db:"title" json:"title"`
	Link          string          `db:"link" json:"link"`
	Description   string          `db:"description" json:"description"`
	Language      string          `db:"language" json:"language"`
	LastBuildDate string          `db:"last_build_date" json:"lastBuildDate"`
	CreatedAt     utils.Timestamp `db:"created_at" json:"createdAt"`
}
