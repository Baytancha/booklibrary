package run

import (
	"test/config"
	"testing"

	"go.uber.org/zap"
)

func TestRoutes(t *testing.T) {
	cfg := config.NewConfig(config.WithPort(8080), config.WithDSN("sqlite3"))

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	app := NewApp(cfg, logger)

	app.Run()

}
