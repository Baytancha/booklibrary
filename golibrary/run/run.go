package run

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test/config"
	"test/internal/db"
	"test/internal/infrastructure/responder"
	"test/internal/modules"
	"test/internal/router"
	"time"

	"test/internal/infrastructure/components"

	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

// type Config struct {
// 	port int
// 	env  string
// 	db   struct {
// 		dsn          string
// 		maxOpenConns int
// 		maxIdleConns int
// 		maxIdleTime  string
// 	}
// }

// App - структура приложения

type App struct {
	cfg      *config.Config
	logger   *zap.Logger
	server   *http.Server
	services *modules.Services
}

func NewApp(conf *config.Config, logger *zap.Logger) *App {
	return &App{
		cfg:    conf,
		logger: logger}
}

func (app *App) Serve() error {

	shutdownErr := make(chan error)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigChan

		app.logger.Info("shutting down", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr <- app.server.Shutdown(ctx)
	}()

	app.logger.Info("starting server", zap.String("addr", app.server.Addr))

	err := app.server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownErr
	if err != nil {
		return err
	}

	app.logger.Info("graceful exit", zap.String("addr", app.server.Addr))
	return nil
}

func (a *App) Run() {
	decoder := godecoder.NewDecoder(jsoniter.Config{})
	responseManager := responder.NewResponder(decoder, a.logger)

	dbx, err := db.NewSqlDB(a.cfg, a.logger)
	if err != nil {
		a.logger.Fatal("error init db", zap.Error(err))
	}
	components := components.NewComponents(responseManager, decoder, a.logger, dbx)
	storages := modules.NewStorages(dbx, a.logger)
	services := modules.NewServices(components, storages)
	a.services = services
	controllers := modules.NewControllers(services, components)

	r := router.Routes(controllers, components)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%v", a.cfg.Port),
		Handler: r,
	}

	a.server = server

}
