package main

import (
	"context"
	"errors"
	handler "github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/delivery"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/deploy_app/usecase"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/handlers/middleware"
	"github.com/Killer-Feature/PaaS_ServerSide/internal/models"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/client_conn/ssh"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/key_value_storage/map_storage"
	servlog "github.com/Killer-Feature/PaaS_ServerSide/pkg/logger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/logger/zaplogger"
	"github.com/Killer-Feature/PaaS_ServerSide/pkg/taskmanager"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"net/netip"
)

func main() {
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

	tm := taskmanager.NewTaskManager[netip.AddrPort](ctx, servLogger)
	processStorage := map_storage.NewMapStorage[netip.Addr, models.TaskProgressMsg]()

	u := usecase.NewDeployAppUsecase(ssh.NewSSHBuilder(), tm, processStorage)
	h := handler.NewDeployAppHandler(servLogger, u)

	middlewares := middleware.NewCommonMiddleware(servLogger)
	mwChain := []echo.MiddlewareFunc{
		middlewares.PanicMiddleware,
		middlewares.RequestIdMiddleware,
		middlewares.AccessLogMiddleware,
		middlewares.PanicMiddleware,
		echomiddleware.CORSWithConfig(middleware.GetCorsConfig([]string{"", "http://localhost:8080"}, 86400)),
	}
	server.Use(mwChain...)

	if err := handlers.Register(server, *h); err != nil {
		log.Fatal(err)
	}
	// metrics := monitoring.RegisterMonitoring(server)

	//m := middleware.NewMiddleware(p.Logger, metrics)
	//m.Register(server)

	g.Go(func() error {
		return server.Start(":8091")
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
