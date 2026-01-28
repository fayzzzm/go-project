package usecase

import (
	"context"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/pkg/utils"
)

type CabinetServicer interface {
	Create(ctx context.Context, c *domain.Cabinet, userID string) error
	GetByID(ctx context.Context, id string) (*domain.Cabinet, error)
	List(ctx context.Context, limit, offset int) ([]domain.Cabinet, error)
	Update(ctx context.Context, c *domain.Cabinet) error
	Delete(ctx context.Context, id string) error
}

type CreateCabinetInput struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	TeamID   string `json:"team_id" binding:"required,uuid"`
	TenantID string `json:"tenant_id"`
}

type UpdateCabinetInput struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	TeamID   string `json:"team_id" binding:"omitempty,uuid"`
	TenantID string `json:"tenant_id"`
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
	return &CabinetUseCase{
		svc: svc,
	}
}

func (uc *CabinetUseCase) Create(ctx context.Context, input CreateCabinetInput, userID string) (*CabinetOutput, error) {
	cabinet := &domain.Cabinet{
		Name:     input.Name,
		Location: utils.StringPtrOrNil(input.Location),
		TeamID:   utils.StringPtrOrNil(input.TeamID),
		TenantID: utils.StringPtrOrNil(input.TenantID),
	}

	if err := uc.svc.Create(ctx, cabinet, userID); err != nil {
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

func (uc *CabinetUseCase) List(ctx context.Context, p domain.Pagination) ([]CabinetOutput, error) {
	p.Normalize()
	cabinets, err := uc.svc.List(ctx, p.Limit, p.Offset)
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
		ID:       id,
		Name:     input.Name,
		Location: utils.StringPtrOrNil(input.Location),
		TeamID:   utils.StringPtrOrNil(input.TeamID),
		TenantID: utils.StringPtrOrNil(input.TenantID),
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
	return &CabinetOutput{
		ID:       c.ID,
		Name:     c.Name,
		Location: utils.StringValue(c.Location),
		TeamID:   utils.StringValue(c.TeamID),
	}
}
