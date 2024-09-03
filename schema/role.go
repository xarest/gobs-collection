package schema

import (
	"time"

	"github.com/uptrace/bun"
)

// CREATE TABLE IF NOT EXISTS "role" (
// 	id  UUID DEFAULT uuid_generate_v4(),
// 	name VARCHAR(255) NOT NULL,
// 	description TEXT NULL,
// 	updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
// 	CONSTRAINT role_pk PRIMARY KEY (id),
// 	CONSTRAINT role_uk_name UNIQUE (name)
// );

type Role struct {
	bun.BaseModel `bun:"table:role"`

	ID          string        `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Name        string        `bun:"name,unique,notnull,type:varchar(255)"`
	Description *string       `bun:"description,type:TEXT"`
	UpdatedAt   time.Time     `bun:"updated_at,type:timestamp,notnull,default:now()"`
	Permissions []*Permission `bun:"-"`
}
