package usecase

import (
	"context"
	"time"

	"github.com/fayzzzm/go-project/internal/domain"
)

type TeamServicer interface {
	Create(ctx context.Context, t *domain.Team) error
	GetByID(ctx context.Context, id string) (*domain.Team, error)
	List(ctx context.Context, limit, offset int, tenantID string) ([]domain.Team, error)
	Update(ctx context.Context, t *domain.Team) error
	Delete(ctx context.Context, id string) error
}

type CreateTeamInput struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	TenantID string `json:"tenant_id"`
}

type UpdateTeamInput struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type TeamOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	TenantID  string    `json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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
		TenantID: input.TenantID,
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

func (uc *TeamUseCase) List(ctx context.Context, limit, offset int, tenantID string) ([]TeamOutput, error) {
	teams, err := uc.svc.List(ctx, limit, offset, tenantID)
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
	team := &domain.Team{
		ID: id,
	}
	if input.Name != "" {
		team.Name = input.Name
	}
	if input.Status != "" {
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

func toTeamOutput(t *domain.Team) *TeamOutput {
	return &TeamOutput{
		ID:        t.ID,
		Name:      t.Name,
		Status:    t.Status,
		TenantID:  t.TenantID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
