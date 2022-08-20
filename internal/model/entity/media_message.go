package entity

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"crawl-worker/pkg/json"
)

type MediaType string

const (
	MediaTypePhoto     MediaType = "photo"
	MediaTypeVideo               = "video"
	MediaTypeAudio               = "audio"
	MediaTypeAnimation           = "animation"
)

type Message struct {
	Type   MediaType `json:"type"`
	FileID string    `json:"file_id"`
}

type Messages []Message

func (u *Messages) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		if len(v) > 0 {
			bytes = make([]byte, len(v))
			copy(bytes, v)
		}
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal json value:", value))
	}

	err := json.Unmarshal(bytes, u)
	return err
}

func (u *Messages) Value() (driver.Value, error) {
	bytes, err := json.Marshal(u)
	return string(bytes), err
}

// GormDataType gorm common data type
func (Messages) GormDataType() string {
	return "json"
}

func (Messages) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}

func (u Messages) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := json.Marshal(u)

	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}

	return gorm.Expr("?", string(data))
}

type MediaMessage struct {
	ID        int64 `gorm:"primaryKey"`
	Messages  Messages
	Type      int `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (MediaMessage) TableName() string {
	return "media_messages"
}

func (MediaMessage) IDField() string {
	return "id"
}
