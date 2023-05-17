package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
	"uptime/models"

	"github.com/uptrace/bun"
)

const EmailCodeLifeSpan = 20 * time.Minute

func SendVerificationEmail(ctx context.Context, db *bun.DB, email string) error {
	user := models.User{}
	const upsert_user_query = `
		insert into users (email) values (?) on conflict (email) do update
		set updated_at = current_timestamp returning *
	`
	if err := db.NewRaw(upsert_user_query, email, email).
		Scan(ctx, &user); err != nil {
		return err
	}

	var code string
	const upsert_code_query = `
		insert into verification_codes (user_id, code, created_at)
		values (?, lower(hex(randomblob(16))), current_timestamp)
		on conflict (user_id) do update
		set code = lower(hex(randomblob(16))), created_at = current_timestamp
		returning code
	`
	if err := db.NewRaw(upsert_code_query, user.ID, user.ID).
		Scan(ctx, &code); err != nil && err != sql.ErrNoRows {
		return err
	}

	fmt.Println("Magic link: ", magicLink(code))
	return nil
}

func ValidateEmailCode(ctx context.Context, db *bun.DB, code string) (models.User, error) {
	user := models.User{}
	const query = `
		select users.*
		from users
			inner join verification_codes vc on vc.user_id = users.id
		where vc.code = ? and vc.created_at > ?
	`
	elapstedTime := time.Now().Add(-EmailCodeLifeSpan)
	if err := db.NewRaw(query, code, elapstedTime).
		Scan(ctx, &user); err != nil && err != sql.ErrNoRows {
		return models.User{}, fmt.Errorf("Invalid code")
	}

	if err := expireEmailCode(ctx, db, code); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func expireEmailCode(ctx context.Context, db *bun.DB, code string) error {
	_, err := db.Exec("delete from verification_codes where code = ?", code)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func RandomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func magicLink(code string) string {
	return "http://localhost:3000/login/verify?code=" + code
}
