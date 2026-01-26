package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

const (
	queryUserCreate     = "SELECT * FROM users.create($1::users.user_request)"
	queryUserGetByID    = "SELECT * FROM users.get_by_id($1::users.user_request)"
	queryUserList       = "SELECT * FROM users.list($1::users.user_request)"
	queryUserUpdate     = "SELECT * FROM users.update($1::users.user_request)"
	queryUserDelete     = "SELECT * FROM users.delete($1::users.user_request)"
	queryUserGetByEmail = "SELECT * FROM users.get_by_email($1::users.user_request)"
)

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	req := UserRequest{
		Email:        &u.Email,
		Name:         u.Name,
		Address:      u.Address,
		Phone:        u.Phone,
		Role:         &u.Role,
		AppMetadata:  u.AppMetadata,
		UserMetadata: u.UserMetadata,
	}

	val, err := ExecQueryOne[domain.User](ctx, r.pool, queryUserCreate, req)
	if err != nil {
		return err
	}
	*u = *val
	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	req := UserRequest{ID: &id}
	return ExecQueryOne[domain.User](ctx, r.pool, queryUserGetByID, req)
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	var role *string
	if u.Role != "" {
		val := u.Role
		role = &val
	}

	// AppMetadata and UserMetadata: if nil, likely means "don't update" or "empty"
	// But JSONB COALESCE default is only if NULL.
	// We pass them as is. domain.User maps are nil-able.

	var email *string
	if u.Email != "" {
		val := u.Email
		email = &val
	}

	req := UserRequest{
		ID:           &u.ID,
		Name:         u.Name,
		Address:      u.Address,
		Phone:        u.Phone,
		Role:         role,
		Email:        email,
		AppMetadata:  u.AppMetadata,
		UserMetadata: u.UserMetadata,
	}
	val, err := ExecQueryOne[domain.User](ctx, r.pool, queryUserUpdate, req)
	if err != nil {
		return err
	}
	*u = *val
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	req := UserRequest{ID: &id}
	_, err := r.pool.Exec(ctx, queryUserDelete, req)
	return err
}

func (r *UserRepo) List(ctx context.Context, limit, offset int) ([]domain.User, error) {
	req := UserRequest{LimitVal: &limit, OffsetVal: &offset}
	return ExecQueryList[domain.User](ctx, r.pool, queryUserList, req)
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	req := UserRequest{Email: &email}
	return ExecQueryOne[domain.User](ctx, r.pool, queryUserGetByEmail, req)
}
