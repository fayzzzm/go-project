package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/domain"
	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type TeamUseCase interface {
	Create(ctx context.Context, input usecase.CreateTeamInput) (*usecase.TeamOutput, error)
	GetByID(ctx context.Context, input usecase.TeamIDInput) (*usecase.TeamOutput, error)
	List(ctx context.Context, p domain.Pagination) ([]usecase.TeamOutput, error)
	ListForUser(ctx context.Context, input usecase.ListForUserInput) ([]usecase.TeamOutput, error)
	Update(ctx context.Context, input usecase.UpdateTeamInput) (*usecase.TeamOutput, error)
	Delete(ctx context.Context, input usecase.TeamIDInput) error
	AddMember(ctx context.Context, input usecase.AddMemberInput) (*usecase.MemberOutput, error)
}

type TeamController struct {
	uc TeamUseCase
}

func NewTeamController(r gin.IRouter, uc TeamUseCase) {
	c := &TeamController{uc: uc}

	teams := r.Group("/teams")
	teams.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
	{
		teams.POST("", middleware.BindJSON[usecase.CreateTeamInput](), utils.Handle(c.Create, http.StatusCreated))
		teams.GET("", utils.Handle(c.List, http.StatusOK))
		teams.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		teams.PUT("/:id", middleware.BindJSON[usecase.UpdateTeamInput](), utils.Handle(c.Update, http.StatusOK))
		teams.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
		teams.POST("/:id/members", middleware.BindJSON[usecase.AddMemberInput](), utils.Handle(c.AddMember, http.StatusOK))
	}
}

func (c *TeamController) Create(ctx *gin.Context) (*usecase.TeamOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateTeamInput](ctx))
}

func (c *TeamController) GetByID(ctx *gin.Context) (*usecase.TeamOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), usecase.TeamIDInput{ID: ctx.Param("id")})
}

func (c *TeamController) List(ctx *gin.Context) ([]usecase.TeamOutput, error) {
	if filterUserID := ctx.Query("user_id"); filterUserID != "" {
		return c.uc.ListForUser(ctx.Request.Context(), usecase.ListForUserInput{UserID: filterUserID})
	}
	return c.uc.List(ctx.Request.Context(), middleware.GetPagination(ctx))
}

func (c *TeamController) Update(ctx *gin.Context) (*usecase.TeamOutput, error) {
	return c.uc.Update(ctx.Request.Context(), middleware.GetBodyWithID[usecase.UpdateTeamInput](ctx, "id"))
}

func (c *TeamController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), usecase.TeamIDInput{ID: ctx.Param("id")})
}

func (c *TeamController) AddMember(ctx *gin.Context) (*usecase.MemberOutput, error) {
	return c.uc.AddMember(ctx.Request.Context(), middleware.GetBodyWithID[usecase.AddMemberInput](ctx, "id"))
}
