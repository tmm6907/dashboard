package models

import (
	"fmt"

	"github.com/tmm6907/dashboard/utils"
)

type User struct {
	ID            utils.UUID      `db:"id" json:"id"`
	OAuthID       string          `db:"oauth_id" json:"oathID"`
	OAuthProvider string          `db:"oauth_provider" json:"oathProvider"`
	FirstName     string          `db:"first_name" json:"firstName"`
	LastName      string          `db:"last_name" json:"lastName"`
	Email         string          `db:"mashboard_email" json:"email"`
	EmailRSSLink  string          `db:"email_rss_link" json:"emailRSSLink"`
	CreatedAt     utils.Timestamp `db:"created_at" json:"createdAt"`
}

func (u User) FullName() string {
	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}
