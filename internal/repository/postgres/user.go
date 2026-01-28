package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

const (
	queryUserCreate  = "SELECT * FROM users.create($1::users.user_request)"
	queryUserGetByID = "SELECT * FROM users.get_by_id($1::users.user_request)"
	queryUserList    = "SELECT * FROM users.list($1::users.user_request)"
	queryUserUpdate  = "SELECT * FROM users.update($1::users.user_request)"
	queryUserDelete  = "SELECT * FROM users.delete($1::users.user_request)"

	queryUserGetForLogin = "SELECT * FROM users.get_for_login($1::text)"
)

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	return ExecQueryUpdate(ctx, r.pool, u, queryUserCreate, NewUserRequest(u))
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return ExecQueryOne[domain.User](ctx, r.pool, queryUserGetByID, UserRequest{ID: &id})
}

func (r *UserRepo) Update(ctx context.Context, u *domain.User) error {
	return ExecQueryUpdate(ctx, r.pool, u, queryUserUpdate, NewUserRequest(u))
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	return Exec(ctx, r.pool, queryUserDelete, UserRequest{ID: &id})
}

func (r *UserRepo) List(ctx context.Context, limit, offset int, teamID, tenantID string) ([]domain.User, error) {
	return ExecQueryList[domain.User](ctx, r.pool, queryUserList, UserRequest{
		LimitVal: &limit, OffsetVal: &offset, TeamID: utils.StringPtrOrNil(teamID), TenantID: utils.StringPtrOrNil(tenantID),
	})
}

func (r *UserRepo) GetForLogin(ctx context.Context, email string) (*domain.User, error) {
	return ExecQueryOne[domain.User](ctx, r.pool, queryUserGetForLogin, email)
}
