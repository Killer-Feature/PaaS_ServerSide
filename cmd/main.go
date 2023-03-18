package main

import (
	"context"
	"errors"
	handler "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/delivery"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/usecase"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/logger/zaplogger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
)

func main() {
	//config := zap.NewDevelopmentConfig()
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      true,
		Encoding:         "json",
		OutputPaths:      []string{"log"},
		ErrorOutputPaths: []string{"stderr"},

		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
	logger, err := zaplogger.NewZapLogger(&config)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal("Error occurred in logger sync")
		}
	}()

	servLogger := servlog.NewServLogger(logger)

	server := echo.New()

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	tm := taskmanager.NewTaskManager(ctx, servLogger)

	u := usecase.NewDeployAppUsecase(ssh.NewSSHBuilder(), tm)
	h := handler.NewDeployAppHandler(servLogger, u)

	middlewares := middleware.NewCommonMiddleware(servLogger)
	if err := handlers.Register(server, *h, middlewares); err != nil {
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
