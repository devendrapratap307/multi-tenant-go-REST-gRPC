//go:build wireinject
// +build wireinject

package app

import (
	"go-multitenant/internal/db"
	"go-multitenant/internal/handler"
	"go-multitenant/internal/repo"
	"go-multitenant/internal/service"

	"github.com/google/wire"
)

type App struct {
	RestHandler *handler.RestHandler
	GRPCServer  *handler.GRPCServer
}

func NewApp(rest *handler.RestHandler, grpc *handler.GRPCServer) *App {
	return &App{RestHandler: rest, GRPCServer: grpc}
}

func InitializeApp() (*App, error) {
	wire.Build(
		db.NewMasterDB,
		db.NewDBProvider,
		repo.NewUserRepo,
		service.NewUserService,
		handler.NewRestHandler,
		handler.NewGRPCServer,
		NewApp,
	)
	return nil, nil
}
