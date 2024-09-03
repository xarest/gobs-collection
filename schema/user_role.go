package schema

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// CREATE TABLE IF NOT EXISTS "user_role" (
// 	user_id UUID NOT NULL,
// 	role_id UUID NOT NULL,
// 	updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
// 	CONSTRAINT user_role_pk PRIMARY KEY (user_id, role_id),
// 	CONSTRAINT user_role_fk_user_id FOREIGN KEY (user_id) REFERENCES user(id),
// 	CONSTRAINT user_role_fk_role_id FOREIGN KEY (role_id) REFERENCES role(id)
// );

type UserRole struct {
	bun.BaseModel `bun:"table:user_role"`

	UserID    uuid.UUID `bun:"user_id,type:uuid,unique:user_role_uk_user_role"`
	RoleID    uuid.UUID `bun:"role_id,type:uuid,unique:user_role_uk_user_role"`
	UpdatedAt time.Time `bun:"updated_at,type:timestamp,nullzero,notnull,default:current_timestamp"`
	User      User      `bun:"rel:belongs-to,join:user_id=id"`
	Role      Role      `bun:"rel:belongs-to,join:role_id=id"`
}
