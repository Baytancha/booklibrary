package db

import (
	"os"
	"test/config"
	"testing"

	"go.uber.org/zap"
)

//"testing"

//"go.uber.org/zap"

func TestDB(t *testing.T) {
	cfg := config.NewConfig(config.WithPort(8080), config.WithDBname("test"), config.WithDSN(os.Getenv("DB_DSN2")))
	_, err := NewSqlDB(cfg, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}

}
