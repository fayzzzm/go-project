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
	Name     string `json:"name" validate:"required"`
	Location string `json:"location"`
	TeamID   string `json:"team_id" validate:"required,uuid"`
}

type UpdateCabinetInput struct {
	ID       string  `json:"-"`
	Name     *string `json:"name"`
	Location *string `json:"location"`
	TeamID   *string `json:"team_id" validate:"omitempty,uuid"`
}

type CabinetIDInput struct {
	ID string `json:"-"`
}

type ListCabinetInput struct {
	Pagination domain.Pagination
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

func (uc *CabinetUseCase) Create(ctx context.Context, input CreateCabinetInput) (*CabinetOutput, error) {
	userID := domain.UserFromContext(ctx)

	cabinet := &domain.Cabinet{
		Name:      input.Name,
		Location:  utils.StringPtrOrNil(input.Location),
		TeamID:    utils.StringPtrOrNil(input.TeamID),
		CreatedBy: utils.StringPtrOrNil(userID),
		UpdatedBy: utils.StringPtrOrNil(userID),
	}

	if err := uc.svc.Create(ctx, cabinet, userID); err != nil {
		return nil, err
	}

	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) GetByID(ctx context.Context, input CabinetIDInput) (*CabinetOutput, error) {
	cabinet, err := uc.svc.GetByID(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) List(ctx context.Context, input ListCabinetInput) ([]CabinetOutput, error) {
	input.Pagination.Normalize()
	cabinets, err := uc.svc.List(ctx, input.Pagination.Limit, input.Pagination.Offset)
	if err != nil {
		return nil, err
	}

	output := make([]CabinetOutput, len(cabinets))
	for i, c := range cabinets {
		output[i] = *toCabinetOutput(&c)
	}
	return output, nil
}

func (uc *CabinetUseCase) Update(ctx context.Context, input UpdateCabinetInput) (*CabinetOutput, error) {
	cabinet, err := uc.svc.GetByID(ctx, input.ID)

	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Name != nil {
		cabinet.Name = *input.Name
	}
	if input.Location != nil {
		cabinet.Location = input.Location // Location is already *string in domain
	}
	if input.TeamID != nil {
		cabinet.TeamID = input.TeamID
	}

	if userID := domain.UserFromContext(ctx); userID != "" {
		cabinet.UpdatedBy = utils.StringPtrOrNil(userID)
	}

	if err := uc.svc.Update(ctx, cabinet); err != nil {
		return nil, err
	}
	return toCabinetOutput(cabinet), nil
}

func (uc *CabinetUseCase) Delete(ctx context.Context, input CabinetIDInput) error {
	return uc.svc.Delete(ctx, input.ID)
}

func toCabinetOutput(c *domain.Cabinet) *CabinetOutput {
	return &CabinetOutput{
		ID:       c.ID,
		Name:     c.Name,
		Location: utils.StringValue(c.Location),
		TeamID:   utils.StringValue(c.TeamID),
	}
}
