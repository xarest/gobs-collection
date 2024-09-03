package schema

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/xarest/gobs-template/utils"
)

const (
	RelUserHasUserRoles string = "UserRoles"
)

type User struct {
	bun.BaseModel `bun:"table:user"`

	ID        uuid.UUID   `bun:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Email     string      `bun:"email,notnull"`
	Password  string      `bun:"password,notnull"`
	UpdatedAt time.Time   `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	UserRoles []*UserRole `bun:"rel:has-many,join:id=user_id"`
	Roles     []*Role     `bun:"-"`
}

type UserFilterParams struct {
	ID    *string `bun:"id" mapstructure:"id,omitempty"`
	Email *string `bun:"email" mapstructure:"email,omitempty"`
}

func (u *UserFilterParams) ToQueryParams() (query string, args []interface{}) {
	if u.ID != nil {
		query = query + "id = ?"
		args = append(args, *u.ID)
	}
	if u.Email != nil {
		query = query + " AND email = ?"
		args = append(args, *u.Email)
	}
	return query, args
}

type UserCreateParams struct {
	Email    string `bun:"email" mapstructure:"email"`
	Password string `bun:"password" mapstructure:"password"`
}

type UserUpdateParams struct{}

func (u *User) GetMany(ctx context.Context,
	tx bun.IDB,
	conds UserFilterParams,
	page Page,
	selectFields []string,
	relations []string,
) ([]User, error) {

	query := tx.NewSelect().Table("user").Model(conds)

	if len(selectFields) > 0 {
		query = query.Column(selectFields...)
	}

	cols, err := utils.GetListFields(conds)
	if err != nil {
		return nil, err
	}
	if len(cols) > 0 {
		query = query.WherePK(cols...)
	}
	for _, relation := range relations {
		query = query.Relation(string(relation))
	}

	page.LoadDefault()
	query = query.Offset(page.Offset).Limit(page.Limit).OrderExpr("? ?", bun.Ident(page.OrderBy), page.SortBy)

	var users []User
	return users, query.Scan(ctx, &users)
}

func (u *User) GetOne(ctx context.Context, tx bun.IDB, params UserFilterParams, relations []string) (*User, error) {
	users, err := u.GetMany(ctx, tx, params, Page{Limit: 1}, nil, relations)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}
	return &users[0], nil
}
