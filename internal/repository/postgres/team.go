package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

const (
	queryTeamCreate     = "SELECT * FROM teams.create($1::teams.team_request)"
	queryTeamGetByID    = "SELECT * FROM teams.get_by_id($1::teams.team_request)"
	queryTeamList       = "SELECT * FROM teams.list($1::teams.team_request)"
	queryTeamUpdate     = "SELECT * FROM teams.update($1::teams.team_request)"
	queryTeamDelete     = "SELECT * FROM teams.delete($1::teams.team_request)"
	queryTeamAddMember  = "SELECT * FROM teams.add_member($1::teams.member_request)"
	queryTeamIsMember   = "SELECT EXISTS(SELECT 1 FROM teams.members WHERE team_id=$1::uuid AND user_id=$2::uuid)"
	queryTeamGetForUser = "SELECT * FROM teams.get_for_user($1::uuid)"
)

func (r *TeamRepo) IsMember(ctx context.Context, teamID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, queryTeamIsMember, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *TeamRepo) AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error) {
	return ExecQueryOne[domain.Member](ctx, r.pool, queryTeamAddMember, MemberRequest{
		TeamID: &teamID, UserID: &userID, Role: &role,
	})
}

func (r *TeamRepo) Create(ctx context.Context, t *domain.Team) error {
	return ExecQueryUpdate(ctx, r.pool, t, queryTeamCreate, NewTeamRequest(t))
}

func (r *TeamRepo) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	return ExecQueryOne[domain.Team](ctx, r.pool, queryTeamGetByID, TeamRequest{ID: &id})
}

func (r *TeamRepo) Update(ctx context.Context, t *domain.Team) error {
	return ExecQueryUpdate(ctx, r.pool, t, queryTeamUpdate, NewTeamRequest(t))
}

func (r *TeamRepo) Delete(ctx context.Context, id string) error {
	return Exec(ctx, r.pool, queryTeamDelete, TeamRequest{ID: &id})
}

func (r *TeamRepo) List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error) {
	return ExecQueryList[domain.Team](ctx, r.pool, queryTeamList, TeamRequest{
		LimitVal: &limit, OffsetVal: &offset, TenantID: utils.StringPtrOrNil(tenantID), UserID: utils.StringPtrOrNil(userID),
	})
}

func (r *TeamRepo) GetForUser(ctx context.Context, userID string) ([]domain.Team, error) {
	return ExecQueryList[domain.Team](ctx, r.pool, queryTeamGetForUser, userID)
}
