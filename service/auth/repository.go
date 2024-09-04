package auth

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/db"
	"github.com/xarest/gobs-collection/schema"
)

type UserRepository struct {
	db   *db.DB
	user *schema.User
}

// Init implements gobs.IServiceInit.
func (u *UserRepository) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{&db.DB{}, &schema.User{}},
	}, nil
}

// Setup implements gobs.IServiceSetup.
func (u *UserRepository) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&u.db, &u.user)
}

func (u *UserRepository) GetUserForSignIn(ctx context.Context, email string) (*schema.User, error) {
	var user schema.User
	if err := u.db.NewSelect().
		Model(&user).
		Column("id", "email", "password").
		Where("email = ?", email).
		Limit(1).
		Relation(schema.RelUserHasUserRoles).
		Scan(ctx, &user); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	roleIDs := make([]uuid.UUID, len(user.UserRoles))
	for i, role := range user.UserRoles {
		roleIDs[i] = role.RoleID
	}

	var rolePermissions []schema.RolePermission
	if err := u.db.NewSelect().
		Model(&rolePermissions).
		Relation(schema.RelRolePermissionHasPermission).
		Relation(schema.RelRolePermissionHasRole).
		Where("role_id IN ?", bun.In(roleIDs)).
		Scan(ctx); err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	mRoles := map[uuid.UUID]*schema.Role{}
	for _, rp := range rolePermissions {
		role, ok := mRoles[rp.RoleID]
		if !ok {
			mRoles[rp.RoleID] = rp.Role
			user.Roles = append(user.Roles, rp.Role)
			role = rp.Role
		}
		role.Permissions = append(role.Permissions, rp.Permission)
	}
	return &user, nil
}

func (u *UserRepository) AddUser(ctx context.Context, user schema.UserCreateParams) (res *schema.User, err error) {
	res = new(schema.User)
	return res, u.db.RunInTx(ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if _, err := tx.NewInsert().
				Table("user").
				Model(&user).
				Exec(ctx, res); err != nil {
				return err
			}
			return nil
		},
	)
}

var _ gobs.IServiceInit = (*UserRepository)(nil)
var _ gobs.IServiceSetup = (*UserRepository)(nil)
