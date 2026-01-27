package postgres

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepo struct {
	pool *pgxpool.Pool
}

func NewTeamRepo(pool *pgxpool.Pool) *TeamRepo {
	return &TeamRepo{pool: pool}
}

const (
	queryTeamCreate    = "SELECT * FROM teams.create($1::teams.team_request)"
	queryTeamGetByID   = "SELECT * FROM teams.get_by_id($1::teams.team_request)"
	queryTeamList      = "SELECT * FROM teams.list($1::teams.team_request)"
	queryTeamUpdate    = "SELECT * FROM teams.update($1::teams.team_request)"
	queryTeamDelete    = "SELECT * FROM teams.delete($1::teams.team_request)"
	queryTeamAddMember = "SELECT * FROM teams.add_member($1::teams.member_request)"
	queryTeamIsMember  = "SELECT EXISTS(SELECT 1 FROM teams.members WHERE team_id=$1::uuid AND user_id=$2::uuid)"
)

func (r *TeamRepo) IsMember(ctx context.Context, teamID, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, queryTeamIsMember, teamID, userID).Scan(&exists)
	return exists, err
}

func (r *TeamRepo) AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error) {
	req := MemberRequest{
		TeamID: &teamID,
		UserID: &userID,
		Role:   &role,
	}
	return ExecQueryOne[domain.Member](ctx, r.pool, queryTeamAddMember, req)
}

func (r *TeamRepo) Create(ctx context.Context, t *domain.Team) error {
	req := TeamRequest{
		Name:     &t.Name,
		Status:   t.Status,
		TenantID: t.TenantID,
	}

	val, err := ExecQueryOne[domain.Team](ctx, r.pool, queryTeamCreate, req)
	if err != nil {
		return err
	}
	*t = *val
	return nil
}

func (r *TeamRepo) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	req := TeamRequest{ID: &id}
	return ExecQueryOne[domain.Team](ctx, r.pool, queryTeamGetByID, req)
}

func (r *TeamRepo) Update(ctx context.Context, t *domain.Team) error {
	req := TeamRequest{
		ID:     &t.ID,
		Name:   &t.Name,
		Status: t.Status,
		// TenantID not updatable usually or passed if needed
	}
	val, err := ExecQueryOne[domain.Team](ctx, r.pool, queryTeamUpdate, req)
	if err != nil {
		return err
	}
	*t = *val
	return nil
}

func (r *TeamRepo) Delete(ctx context.Context, id string) error {
	req := TeamRequest{ID: &id}
	_, err := r.pool.Exec(ctx, queryTeamDelete, req)
	return err
}

func (r *TeamRepo) List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error) {
	var tID, uID *string
	if tenantID != "" {
		tID = &tenantID
	}
	if userID != "" {
		uID = &userID
	}
	req := TeamRequest{LimitVal: &limit, OffsetVal: &offset, TenantID: tID, UserID: uID}
	return ExecQueryList[domain.Team](ctx, r.pool, queryTeamList, req)
}
