package auth

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"github.com/xarest/gobs"
	"github.com/xarest/gobs-collection/lib/db"
	"github.com/xarest/gobs-collection/lib/logger"
	"github.com/xarest/gobs-collection/schema"
)

type UserRepository struct {
	log  logger.ILogger
	db   *db.DB
	user *schema.User
}

// Init implements gobs.IServiceInit.
func (u *UserRepository) Init(ctx context.Context) (*gobs.ServiceLifeCycle, error) {
	return &gobs.ServiceLifeCycle{
		Deps: gobs.Dependencies{logger.NewILogger(), &db.DB{}, &schema.User{}},
	}, nil
}

// Setup implements gobs.IServiceSetup.
func (u *UserRepository) Setup(ctx context.Context, deps ...gobs.IService) error {
	return gobs.Dependencies(deps).Assign(&u.log, &u.db, &u.user)
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
	if len(user.UserRoles) == 0 {
		return &user, nil
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
	res = &schema.User{
		Email:    user.Email,
		Password: user.Password,
	}
	return res, u.db.RunInTx(ctx,
		&sql.TxOptions{},
		func(ctx context.Context, tx bun.Tx) error {
			if _, err := tx.NewInsert().
				Model(res).
				Exec(ctx); err != nil {
				u.log.Error(err)
				return err
			}
			return nil
		},
	)
}

var _ gobs.IServiceInit = (*UserRepository)(nil)
var _ gobs.IServiceSetup = (*UserRepository)(nil)
