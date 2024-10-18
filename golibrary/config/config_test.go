package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config == nil {
		t.Fatal("config is nil")
	}
}

func TestWithPort(t *testing.T) {
	port := 8080
	config := NewConfig(WithPort(port))

	if config.Port != port {
		t.Errorf("expected %d, got %d", port, config.Port)
	}
}

func TestWithEnv(t *testing.T) {
	env := "test"
	config := NewConfig(WithEnv(env))

	if config.Env != env {
		t.Errorf("expected %s, got %s", env, config.Env)
	}
}

func TestWithDSN(t *testing.T) {
	dsn := "test"
	config := NewConfig(WithDSN(dsn))

	if config.Db.Dsn != dsn {
		t.Errorf("expected %s, got %s", dsn, config.Db.Dsn)
	}
}

func TestWithMaxOpenConns(t *testing.T) {
	maxOpenConns := 10
	config := NewConfig(WithMaxOpenConns(maxOpenConns))

	if config.Db.MaxOpenConns != maxOpenConns {
		t.Errorf("expected %d, got %d", maxOpenConns, config.Db.MaxOpenConns)
	}
}

func TestWithMaxIdleConns(t *testing.T) {
	maxIdleConns := 10
	config := NewConfig(WithMaxIdleConns(maxIdleConns))

	if config.Db.MaxIdleConns != maxIdleConns {
		t.Errorf("expected %d, got %d", maxIdleConns, config.Db.MaxIdleConns)
	}
}

func TestWithMaxIdleTime(t *testing.T) {
	maxIdleTime := "10s"
	config := NewConfig(WithMaxIdleTime(maxIdleTime))

	if config.Db.MaxIdleTime != maxIdleTime {
		t.Errorf("expected %s, got %s", maxIdleTime, config.Db.MaxIdleTime)
	}
}
