package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
)

type CabinetServicer interface {
	Create(ctx context.Context, c *domain.Cabinet) error
	GetByID(ctx context.Context, id string) (*domain.Cabinet, error)
	List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error)
	Update(ctx context.Context, c *domain.Cabinet) error
	Delete(ctx context.Context, id string) error
}

type CreateCabinetInput struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	TeamID   string `json:"team_id" binding:"required,uuid"`
}

type UpdateCabinetInput struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	TeamID   string `json:"team_id" binding:"omitempty,uuid"`
}

type CabinetOutput struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	TeamID   string `json:"team_id"`
}

type CabinetUseCase struct {
	svc CabinetServicer
}

func NewCabinetUseCase(svc CabinetServicer) *CabinetUseCase {
	return &CabinetUseCase{svc: svc}
}

func (uc *CabinetUseCase) Create(ctx context.Context, input CreateCabinetInput) (*CabinetOutput, error) {
	var location, teamID *string
	if input.Location != "" {
		val := input.Location
		location = &val
	}
	if input.TeamID != "" {
		val := input.TeamID
		teamID = &val
	}

	cabinet := &domain.Cabinet{
		Name:     input.Name,
		Location: location,
		TeamID:   teamID,
	}

	if err := uc.svc.Create(ctx, cabinet); err != nil {
		return nil, err
	}

	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) GetByID(ctx context.Context, id string) (*CabinetOutput, error) {
	cabinet, err := uc.svc.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) List(ctx context.Context, limit, offset int) ([]CabinetOutput, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}
	cabinets, err := uc.svc.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	output := make([]CabinetOutput, len(cabinets))
	for i, c := range cabinets {
		output[i] = *toCabinetOutput(&c)
	}
	return output, nil
}

func (uc *CabinetUseCase) Update(ctx context.Context, id string, input UpdateCabinetInput) (*CabinetOutput, error) {
	cabinet := &domain.Cabinet{
		ID: id,
	}
	if input.Name != "" {
		cabinet.Name = input.Name
	}
	if input.Location != "" {
		val := input.Location
		cabinet.Location = &val
	}
	if input.TeamID != "" {
		val := input.TeamID
		cabinet.TeamID = &val
	}

	if err := uc.svc.Update(ctx, cabinet); err != nil {
		return nil, err
	}
	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) Delete(ctx context.Context, id string) error {
	return uc.svc.Delete(ctx, id)
}

func toCabinetOutput(c *domain.Cabinet) *CabinetOutput {
	location := ""
	if c.Location != nil {
		location = *c.Location
	}
	teamID := ""
	if c.TeamID != nil {
		teamID = *c.TeamID
	}
	return &CabinetOutput{
		ID:       c.ID,
		Name:     c.Name,
		Location: location,
		TeamID:   teamID,
	}
}
