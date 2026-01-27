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
	Name     string `json:"name" binding:"required"`
	Status   string `json:"status"`
	TenantID string `json:"tenant_id" binding:"required,uuid"`
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
	var status *string
	if input.Status != "" {
		val := input.Status
		status = &val
	}
	var tenantID *string
	if input.TenantID != "" {
		val := input.TenantID
		tenantID = &val
	}
	team := &domain.Team{
		Name:     input.Name,
		Status:   status,
		TenantID: tenantID,
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
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
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
		val := input.Status
		team.Status = &val
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
