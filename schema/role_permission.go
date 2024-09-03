package schema

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// CREATE TABLE IF NOT EXISTS "role_permission" (
//
//	role_id UUID NOT NULL,
//	permission_id UUID NOT NULL,
//	updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
//	CONSTRAINT role_permission_pk PRIMARY KEY (role_id, permission_id),
//	CONSTRAINT role_permission_fk_role_id FOREIGN KEY (role_id) REFERENCES role(id),
//	CONSTRAINT role_permission_fk_permission_id FOREIGN KEY (permission_id) REFERENCES permission(id)
//
// );

type RolePermission struct {
	bun.BaseModel `bun:"table:role_permission"` // Specify the table name if needed

	RoleID       uuid.UUID   `bun:"role_id"`
	PermissionID uuid.UUID   `bun:"permission_id"`
	Role         *Role       `bun:"rel:belongs-to,join:role_id=id"`
	Permission   *Permission `bun:"rel:belongs-to,join:permission_id=id"`
}

const (
	RelRolePermissionHasRole       = "Role"
	RelRolePermissionHasPermission = "Permission"
)
