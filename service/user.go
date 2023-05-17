package service

import (
	"context"
	"database/sql"
	"uptime/models"

	"github.com/uptrace/bun"
)

func FindUserByToken(ctx context.Context, db *bun.DB, token string, user *models.User) error {
	if err := db.NewSelect().
		Model(user).
		Join("join sessions on sessions.user_id = u.id").
		Where("sessions.token = ?", token).
		Scan(ctx); err != nil || err != sql.ErrNoRows {
		return err
	}
	return nil
}
