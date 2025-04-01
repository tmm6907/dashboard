package models

import "github.com/tmm6907/dashboard/utils"

type Collection struct {
	ID        uint            `db:"id" json:"id"`
	Name      string          `db:"name" json:"name"`
	CreatedAt utils.Timestamp `db:"created_at" json:"createdAt"`
}
