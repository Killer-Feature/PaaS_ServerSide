package main

import (
	handler "KillerFeature/ServerSide/internal/deploy_app/delivery"
	"KillerFeature/ServerSide/internal/deploy_app/usecase"
	"KillerFeature/ServerSide/pkg/client_conn/ssh"
	"context"
	"log"
	"net/http"

	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"KillerFeature/ServerSide/internal/handlers"
	"KillerFeature/ServerSide/pkg/taskmanager"
)

func main() {

	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, err := config.Build()
	if err != nil {
		log.Fatal("zap logger build error")
	}
	logger := prLogger
	defer func(prLogger *zap.Logger) {
		err = prLogger.Sync()
		if err != nil {
			log.Fatal(err)
		}
	}(prLogger)

	server := echo.New()

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	tm := taskmanager.NewTaskManager(ctx)

	u := usecase.NewDeployAppUsecase(ssh.NewSSHBuilder(), tm)
	h := handler.NewDeployAppHandler(logger, u)
	if err := handlers.Register(server, *h); err != nil {
		log.Fatal(err)
	}
	// metrics := monitoring.RegisterMonitoring(server)

	//m := middleware.NewMiddleware(p.Logger, metrics)
	//m.Register(server)

	g.Go(func() error {
		return server.Start(":8090")
	})

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("error shutdown with error", zap.Error(err))
			ctx := context.Background()
			//u.CloseService()
			_ = server.Shutdown(ctx)
		}
	}

}
