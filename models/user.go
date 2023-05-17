package models

import (
	"time"

	"github.com/uptrace/bun"
)

type (
	User struct {
		bun.BaseModel `bun:"table:users,alias:u"`
		ID            int64 `bun:"id,pk,autoincrement"`
		Email         string
		Session       *Session  `bun:"rel:has-one"`
		CreatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"-"`
		UpdatedAt     time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"-"`
	}

	VertificationCode struct {
		bun.BaseModel `bun:"table:verification_codes,alias:codes"`
		Code          string
		UserID        int64
		CreatedAt     time.Time
	}

	Session struct {
		bun.BaseModel `bun:"table:sessions"`
		Token         string
		UserID        int64
		User          *User `bun:"rel:belongs-to,join:user_id=id"`
	}
)
