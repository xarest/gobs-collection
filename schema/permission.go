package schema

import "github.com/uptrace/bun"

// CREATE TABLE IF NOT EXISTS "permission" (
// 	id  UUID DEFAULT uuid_generate_v4(),
// 	name VARCHAR(255) NOT NULL,
// 	description TEXT NOT NULL,
// 	updatedAt TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT NOW(),
// 	CONSTRAINT permission_pk PRIMARY KEY (id),
// 	CONSTRAINT permission_uk_name UNIQUE (name)
// );

type Permission struct {
	bun.BaseModel `bun:"table:role_permission"` // Specify the table name if needed

	ID          string  `bun:"id,pk,type:uuid,default:uuid_generate_v4()"`
	Name        string  `bun:"name,unique,notnull,type:varchar(255)"`
	Description *string `bun:"description,type:TEXT"`
}
