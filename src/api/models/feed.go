package models

import (
	"github.com/tmm6907/dashboard/utils"
)

type Feed struct {
	FeedID        utils.UUID      `db:"feed_id" json:"feedId"`
	Title         string          `db:"title" json:"title"`
	Link          string          `db:"link" json:"link"`
	Image         string          `db:"image" json:"image"`
	AltText       string          `db:"alt_text" json:"altText"`
	MediaType     string          `db:"media_type" json:"mediaType"`
	Categories    string          `db:"categories" json:"categories"`
	Description   string          `db:"description" json:"description"`
	Language      string          `db:"language" json:"language"`
	LastBuildDate utils.Timestamp `db:"last_build_date" json:"lastBuildDate"`
	CreatedAt     utils.Timestamp `db:"created_at" json:"createdAt"`
}
