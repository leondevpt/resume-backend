//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/leondevpt/resume-backend/internal/app"
	"github.com/leondevpt/resume-backend/internal/biz"
	"github.com/leondevpt/resume-backend/internal/config"
	"github.com/leondevpt/resume-backend/internal/server"
	"github.com/leondevpt/resume-backend/internal/service"
)

func wireApp(c context.Context, cfg *config.Config) (*app.App, error) {
	panic(wire.Build(server.ProviderSet, biz.ProviderSet, service.ProviderSet, app.NewApp))
}
