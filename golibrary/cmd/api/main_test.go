package main

import (
	"fmt"
	"os"
	"test/config"
	"test/run"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestMainfunc(t *testing.T) {
	go func() {
		cfg := config.NewConfig(config.WithPort(8080), config.WithDBname("postgres"), config.WithDSN(os.Getenv("DB_DSN")))

		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}

		app := run.NewApp(cfg, logger)
		fmt.Println(app)
	}()
	time.Sleep(2 * time.Second)
	t.Log("main finished")
}
