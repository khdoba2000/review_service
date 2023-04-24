package app

import (
	"context"
	"fmt"
	"monorepo/src/libs/log"
	"monorepo/src/review_service/configs"
	"monorepo/src/review_service/controller"
	"monorepo/src/review_service/handler/rpc"
	"monorepo/src/review_service/storage/repo"
	"monorepo/src/review_service/tracer"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loggerInit(config *configs.Configuration) log.Factory {
	fmt.Println("logInit")
	loggerForTracer, _ := zap.NewDevelopment(
		zap.AddStacktrace(zapcore.FatalLevel),
		zap.AddCallerSkip(1),
	)

	zapLogger := loggerForTracer.With(zap.String("service", config.ServiceName))
	logger := log.NewFactory(zapLogger)
	return logger
}

type App struct {
	engine *fx.App
}

// Start starts app with context spesified
func (a *App) Start(ctx context.Context) {
	a.engine.Start(ctx)
}

// Run starts the application, blocks on the signals channel, and then gracefully shuts the application down
func (a *App) Run() {
	a.engine.Run()
}

// New returns fx app
func New() App {

	engine := fx.New(
		fx.Provide(
			configs.Config,
			loggerInit,
			repo.New,
			controller.New,
			tracer.Load,
			rpc.New,
		),

		fx.Invoke(
			rpc.Start,
		),

		fx.WithLogger(
			func() fxevent.Logger {
				logger, _ := zap.NewDevelopment(
					zap.AddStacktrace(zapcore.FatalLevel),
					zap.AddCallerSkip(1),
				)

				return &fxevent.ZapLogger{
					Logger: logger,
				}
			},
		),
	)

	return App{engine: engine}
}
