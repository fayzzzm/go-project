package main

import (
	"context"
	"log"
	"net/http"
	"os"

	delivery "github.com/fayzzzm/go-project/internal/delivery/http"
	"github.com/fayzzzm/go-project/internal/middleware"
	"github.com/fayzzzm/go-project/internal/repository/postgres"
	"github.com/fayzzzm/go-project/internal/service"
	"github.com/fayzzzm/go-project/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	_ = godotenv.Load()

	fx.New(
		fx.Provide(
			NewDatabasePool,
			NewGinEngine,

			// 1. Repositories -> Service Interfaces (Consumer defined)
			fx.Annotate(
				postgres.NewDeviceRepo,
				fx.As(new(service.DeviceRepository)),
			),
			fx.Annotate(
				postgres.NewUserRepo,
				fx.As(new(service.UserRepository)),
			),
			fx.Annotate(
				postgres.NewCabinetRepo,
				fx.As(new(service.CabinetRepository)),
			),
			fx.Annotate(
				postgres.NewTeamRepo,
				fx.As(new(service.TeamRepository)),
			),
			fx.Annotate(
				postgres.NewDeviceProfileRepo,
				fx.As(new(service.DeviceProfileRepository)),
			),

			// 2. Services -> UseCase Interfaces (Consumer defined)
			fx.Annotate(
				service.NewDeviceService,
				fx.As(new(usecase.DeviceServicer)),
			),
			fx.Annotate(
				service.NewUserService,
				fx.As(new(usecase.UserServicer)),
			),
			fx.Annotate(
				service.NewCabinetService,
				fx.As(new(usecase.CabinetServicer)),
			),
			fx.Annotate(
				service.NewTeamService,
				fx.As(new(usecase.TeamServicer)),
			),
			fx.Annotate(
				service.NewDeviceProfileService,
				fx.As(new(usecase.DeviceProfileServicer)),
			),

			// 3. UseCases -> Delivery Interfaces (Consumer defined)
			fx.Annotate(
				usecase.NewDeviceUseCase,
				fx.As(new(delivery.DeviceUseCase)),
			),
			fx.Annotate(
				usecase.NewUserUseCase,
				fx.As(new(delivery.UserUseCase)),
			),
			fx.Annotate(
				usecase.NewCabinetUseCase,
				fx.As(new(delivery.CabinetUseCase)),
			),
			fx.Annotate(
				usecase.NewTeamUseCase,
				fx.As(new(delivery.TeamUseCase)),
			),
			fx.Annotate(
				usecase.NewDeviceProfileUseCase,
				fx.As(new(delivery.DeviceProfileUseCase)),
			),

			// 4. Controllers
		),
		fx.Invoke(
			RegisterRoutes,
		),
	).Run()
}

func NewDatabasePool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/itsware?sslmode=disable"
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		types := []string{
			"devices.device_request",
			"users.user_request",
			"cabinets.cabinet_request",
			"teams.team_request",
			"device_profiles.device_profile_request",
		}
		for _, t := range types {
			dt, err := conn.LoadType(ctx, t)
			if err != nil {
				log.Printf("Warning: Failed to load type %s: %v", t, err)
				continue
			}
			conn.TypeMap().RegisterType(dt)
		}

		conn.TypeMap().RegisterDefaultPgType(postgres.DeviceRequest{}, "devices.device_request")
		conn.TypeMap().RegisterDefaultPgType(postgres.UserRequest{}, "users.user_request")
		conn.TypeMap().RegisterDefaultPgType(postgres.CabinetRequest{}, "cabinets.cabinet_request")
		conn.TypeMap().RegisterDefaultPgType(postgres.TeamRequest{}, "teams.team_request")
		conn.TypeMap().RegisterDefaultPgType(postgres.DeviceProfileRequest{}, "device_profiles.device_profile_request")

		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func NewGinEngine() *gin.Engine {
	r := gin.Default()
	// Global security middleware
	r.Use(middleware.AuditMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	return r
}

func RegisterRoutes(
	lc fx.Lifecycle,
	r *gin.Engine,
	deviceUC delivery.DeviceUseCase,
	userUC delivery.UserUseCase,
	cabinetUC delivery.CabinetUseCase,
	teamUC delivery.TeamUseCase,
	deviceProfileUC delivery.DeviceProfileUseCase,
) {
	api := r.Group("/api/v1")

	delivery.NewDeviceController(api, deviceUC)
	delivery.NewUserController(api, userUC)
	delivery.NewCabinetController(api, cabinetUC)
	delivery.NewTeamController(api, teamUC)
	delivery.NewDeviceProfileController(api, deviceProfileUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server starting on %s", port)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
