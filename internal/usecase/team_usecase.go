package usecase

import (
	"context"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
)

type TeamServicer interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID, userID string) ([]domain.Team, error)
	GetForUser(ctx context.Context, userID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
	AddMember(ctx context.Context, teamID, userID, role string) (*domain.Member, error)
	IsMember(ctx context.Context, teamID, userID string) (bool, error)
}

type CreateTeamInput struct {
	Name     string  `json:"name" binding:"required"`
	Status   *string `json:"status"`
	TenantID string  `json:"tenant_id"` // Injected by middleware
}

type UpdateTeamInput struct {
	Name   *string `json:"name"`
	Status *string `json:"status"`
}

type AddMemberInput struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Role   string `json:"role"`
}

type TeamOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	TenantID  string    `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MemberOutput struct {
	TeamID    string    `json:"team_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type TeamUseCase struct {
	svc TeamServicer
}

func NewTeamUseCase(svc TeamServicer) *TeamUseCase {
	return &TeamUseCase{svc: svc}
}

func (uc *TeamUseCase) Create(ctx context.Context, input CreateTeamInput) (*TeamOutput, error) {
	team := &domain.Team{
		Name:     input.Name,
		Status:   input.Status,
		TenantID: utils.StringPtrOrNil(input.TenantID),
	}

	if err := uc.svc.Create(ctx, team); err != nil {
		return nil, err
	}

	return toTeamOutput(team), nil
}

func (uc *TeamUseCase) GetByID(ctx context.Context, id string) (*TeamOutput, error) {
	team, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTeamOutput(team), nil
}

func (uc *TeamUseCase) List(ctx context.Context, p domain.Pagination, tenantID, userID string) ([]TeamOutput, error) {
	p.Normalize()
	teams, err := uc.svc.List(ctx, p.Limit, p.Offset, tenantID, userID)
	if err != nil {
		return nil, err
	}

	output := make([]TeamOutput, len(teams))
	for i, t := range teams {
		output[i] = *toTeamOutput(&t)
	}
	return output, nil
}

func (uc *TeamUseCase) ListForUser(ctx context.Context, userID string) ([]TeamOutput, error) {
	teams, err := uc.svc.GetForUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	output := make([]TeamOutput, len(teams))
	for i, t := range teams {
		output[i] = *toTeamOutput(&t)
	}
	return output, nil
}

func (uc *TeamUseCase) Update(ctx context.Context, id string, input UpdateTeamInput) (*TeamOutput, error) {
	team, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		team.Name = *input.Name
	}
	if input.Status != nil {
		team.Status = input.Status
	}

	if err := uc.svc.Update(ctx, team); err != nil {
		return nil, err
	}
	return toTeamOutput(team), nil
}

func (uc *TeamUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}

func (uc *TeamUseCase) AddMember(ctx context.Context, teamID string, input AddMemberInput) (*MemberOutput, error) {
	member, err := uc.svc.AddMember(ctx, teamID, input.UserID, input.Role)
	if err != nil {
		return nil, err
	}
	return &MemberOutput{
		TeamID:    member.TeamID,
		UserID:    member.UserID,
		Role:      member.Role,
		CreatedAt: member.CreatedAt,
	}, nil
}

func toTeamOutput(t *domain.Team) *TeamOutput {
	status := ""
	if t.Status != nil {
		status = *t.Status
	}
	tenantID := ""
	if t.TenantID != nil {
		tenantID = *t.TenantID
	}
	return &TeamOutput{
		ID:        t.ID,
		Name:      t.Name,
		Status:    status,
		TenantID:  tenantID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
