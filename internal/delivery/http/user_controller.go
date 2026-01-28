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

type UserUseCase interface {
	Create(ctx context.Context, input usecase.CreateUserInput) (*usecase.UserOutput, error)
	GetByID(ctx context.Context, input usecase.UserIDInput) (*usecase.UserOutput, error)
	List(ctx context.Context, input usecase.ListUserInput) ([]usecase.UserOutput, error)
	Update(ctx context.Context, input usecase.UpdateUserInput) (*usecase.UserOutput, error)
	Delete(ctx context.Context, input usecase.UserIDInput) error
}

type UserController struct {
	uc UserUseCase
}

func NewUserController(r gin.IRouter, uc UserUseCase) {
	c := &UserController{uc: uc}

	users := r.Group("/users")
	{
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(), middleware.RequireTenant())
		{
			protected.POST("", middleware.RequireRole(domain.RoleAdmin), middleware.BindJSON[usecase.CreateUserInput](), utils.Handle(c.Create, http.StatusCreated))
			protected.GET("", utils.Handle(c.List, http.StatusOK))
			protected.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
			protected.PUT("/:id", middleware.RequireRole(domain.RoleAdmin), middleware.BindJSON[usecase.UpdateUserInput](), utils.Handle(c.Update, http.StatusOK))
			protected.DELETE("/:id", middleware.RequireRole(domain.RoleAdmin), utils.Handle(c.Delete, http.StatusNoContent))
		}
	}
}

func (c *UserController) Create(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateUserInput](ctx))
}

func (c *UserController) List(ctx *gin.Context) ([]usecase.UserOutput, error) {
	return c.uc.List(ctx.Request.Context(), usecase.ListUserInput{
		Pagination: middleware.GetPagination(ctx),
		TeamID:     ctx.Query("team_id"),
	})
}

func (c *UserController) GetByID(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), usecase.UserIDInput{ID: ctx.Param("id")})
}

func (c *UserController) Update(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.Update(ctx.Request.Context(), middleware.GetBodyWithID[usecase.UpdateUserInput](ctx, "id"))
}

func (c *UserController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), usecase.UserIDInput{ID: ctx.Param("id")})
}
