package app

import (
	"github.com/markgregr/bestHack_support_gRPC_server/internal/adapters/db/postgresql"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/adapters/db/redis"
	grpcapp "github.com/markgregr/bestHack_support_gRPC_server/internal/app/grpc"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/config"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/auth"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/user"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/workflow/cases"
	"github.com/markgregr/bestHack_support_gRPC_server/internal/services/workflow/tasks"
	"github.com/markgregr/bestHack_support_gRPC_server/pkg/gmiddleware"
	"github.com/sirupsen/logrus"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *logrus.Entry, cfg *config.Config) *App {
	postgre, err := postgresql.New(log.Logger, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	redis, err := redis.New(log.Logger, &cfg.Redis)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log.Logger, postgre, redis, postgre, postgre, cfg.JWT.TokenTTL)

	userService := user.New(log.Logger, postgre)

	taskService := tasks.New(log.Logger, postgre, postgre, postgre, postgre, *userService)

	caseService := cases.New(log.Logger, postgre, postgre, postgre, *userService)

	authMd := gmiddleware.NewAuthInterceptor(cfg.JWT.TokenKey, authService)

	grpcApp := grpcapp.New(log, authService, taskService, caseService, authMd, cfg.GRPC.Port, cfg.GRPC.Host)

	return &App{
		GRPCSrv: grpcApp,
	}

}
