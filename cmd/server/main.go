package main

import (
	"context"
	"github.com/leondevpt/resume-backend/internal/app"
	"github.com/leondevpt/resume-backend/internal/config"
	"github.com/leondevpt/resume-backend/pkg/logging"
	"github.com/leondevpt/resume-backend/version"
	"os/signal"
	"syscall"
)

func main() {
	ctx, done := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	/*
		logger := logging.NewLoggerFromEnv().
			With("build_id", version.GitCommit).
			With("build_tag", version.GitTag)

	*/

	err := config.Init()
	if err != nil {
		panic(err)
	}
	logger := logging.NewLogger("info", true).
		With("build_id", version.GitCommit).
		With("build_tag", version.GitTag)

	defer logger.Sync()

	ctx = logging.WithLogger(ctx, logger)

	defer func() {
		done()
		if r := recover(); r != nil {
			logger.Fatalw("application panic", "panic", r)
		}
	}()

	a, err := wireApp(ctx, config.Conf)
	if err != nil {
		logger.Fatal(err)
	}

	err = app.Run(ctx, a)
	done()

	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("successful shutdown")

}
