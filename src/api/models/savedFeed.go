package models

import "github.com/tmm6907/dashboard/utils"

type SavedFeed struct {
	UserID     utils.UUID      `db:"user_id" json:"userID"`
	FeedItemID uint            `db:"feed_item_id" json:"feedItemID"`
	CreatedAt  utils.Timestamp `db:"created_at" json:"createdAt"`
}
