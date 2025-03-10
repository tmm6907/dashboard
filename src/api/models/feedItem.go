package models

import (
	"github.com/tmm6907/dashboard/utils"
)

type FeedItem struct {
	ID          uint            `db:"id" json:"id"`
	FeedID      utils.UUID      `db:"feed_id" json:"feedId"`
	Title       string          `db:"title" json:"title"`
	Link        string          `db:"link" json:"link"`
	Description string          `db:"description" json:"description"`
	Image       string          `db:"image" json:"image"`
	PubDate     string          `db:"pub_date" json:"pubDate"`
	GUID        string          `db:"guid" json:"guid"`
	CreatedAt   utils.Timestamp `db:"created_at" json:"createdAt"`
}
