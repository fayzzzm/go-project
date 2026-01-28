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
	GetByID(ctx context.Context, id string) (*usecase.UserOutput, error)
	List(ctx context.Context, p domain.Pagination, teamID, tenantID string) ([]usecase.UserOutput, error)
	Update(ctx context.Context, id string, input usecase.UpdateUserInput) (*usecase.UserOutput, error)
	Delete(ctx context.Context, id string) error
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
			protected.POST("", middleware.RequireRole("admin"), middleware.BindJSON[usecase.CreateUserInput](), utils.Handle(c.Create, http.StatusCreated))
			protected.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
			protected.PUT("/:id", middleware.RequireRole("admin"), middleware.BindJSON[usecase.UpdateUserInput](), utils.Handle(c.Update, http.StatusOK))
			protected.DELETE("/:id", middleware.RequireRole("admin"), utils.Handle(c.Delete, http.StatusNoContent))
		}
	}
}

func (c *UserController) Create(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.Create(ctx.Request.Context(), middleware.GetBody[usecase.CreateUserInput](ctx))
}

func (c *UserController) GetByID(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.GetByID(ctx.Request.Context(), ctx.Param("id"))
}

func (c *UserController) Update(ctx *gin.Context) (*usecase.UserOutput, error) {
	return c.uc.Update(ctx.Request.Context(), ctx.Param("id"), middleware.GetBody[usecase.UpdateUserInput](ctx))
}

func (c *UserController) Delete(ctx *gin.Context) (struct{}, error) {
	return struct{}{}, c.uc.Delete(ctx.Request.Context(), ctx.Param("id"))
}
