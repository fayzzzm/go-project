package http

import (
	"context"
	"net/http"

	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/fayzzzm/go-project/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	Create(ctx context.Context, input usecase.CreateUserInput) (*usecase.UserOutput, error)
	GetByID(ctx context.Context, id string) (*usecase.UserOutput, error)
	List(ctx context.Context, limit, offset int) ([]usecase.UserOutput, error)
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
		users.POST("", middleware.BindJSON[usecase.CreateUserInput](), utils.Handle(c.Create, http.StatusCreated))
		users.GET("", utils.Handle(c.List, http.StatusOK))
		users.GET("/:id", utils.Handle(c.GetByID, http.StatusOK))
		users.PUT("/:id", middleware.BindJSON[usecase.UpdateUserInput](), utils.Handle(c.Update, http.StatusOK))
		users.DELETE("/:id", utils.Handle(c.Delete, http.StatusNoContent))
	}
}

func (c *UserController) Create(ctx *gin.Context) (any, error) {
	input := middleware.GetBody[usecase.CreateUserInput](ctx)

	return c.uc.Create(ctx.Request.Context(), input)
}

func (c *UserController) GetByID(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return c.uc.GetByID(ctx.Request.Context(), id)
}

func (c *UserController) List(ctx *gin.Context) (any, error) {
	limit, offset := middleware.GetPagination(ctx)

	return c.uc.List(ctx.Request.Context(), limit, offset)
}

func (c *UserController) Update(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	input := middleware.GetBody[usecase.UpdateUserInput](ctx)

	return c.uc.Update(ctx.Request.Context(), id, input)
}

func (c *UserController) Delete(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	err := c.uc.Delete(ctx.Request.Context(), id)
	return nil, err
}
