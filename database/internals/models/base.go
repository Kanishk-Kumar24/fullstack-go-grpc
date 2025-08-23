package models

import (
	"time"

	"github.com/google/uuid"
	// "github.com/uptrace/bun"
)

// BaseModel defines common fields for all models.
type BaseModel struct {
	ID        int64     `bun:"id,pk,autoincrement"`
	UniqueID  uuid.UUID `bun:"unique_id,type:uuid,default:gen_random_uuid()"`
	CreatedAt time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,notnull,default:current_timestamp"`
	DeletedAt time.Time `bun:"deleted_at,soft_delete,nullzero"`
}