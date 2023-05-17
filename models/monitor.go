package models

import (
	"time"

	"github.com/uptrace/bun"
)

type (
	Monitor struct {
		bun.BaseModel `bun:"table:monitors,alias:m"`

		ID        int64     `bun:"id,pk,autoincrement" json:"id"`
		Name      string    `bun:",notnull" json:"name"`
		URL       string    `bun:",notnull" json:"url"`
		Protocol  string    `bun:",notnull" json:"protocol"`
		Interval  int64     `bun:",notnull" json:"interval"`
		UserID    int64     `bun:",notnull" json:"user_id"`
		User      *User     `bun:"rel:belongs-to,join:user_id=id"`
		CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
		UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	}

	MonitorData struct {
		bun.BaseModel `bun:"table:monitor_data,alias:md"`

		ID           int64     `bun:"id,pk,autoincrement" json:"id"`
		StatusCode   int       `bun:",notnull" json:"status_code"`
		ResponseTime float64   `bun:",notnull" json:"response_time"`
		MonitorID    int64     `bun:",notnull" json:"monitor_id"`
		Monitor      *Monitor  `bun:"rel:belongs-to,join:monitor_id=id"`
		CreatedAt    time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	}
)
