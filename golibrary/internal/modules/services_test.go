package modules

import (
	"test/internal/infrastructure/components"
	"test/internal/infrastructure/responder"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

func TestNewServices(t *testing.T) {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	decoder := godecoder.NewDecoder(jsoniter.Config{})
	responseManager := responder.NewResponder(decoder, logger)

	components := components.NewComponents(responseManager, decoder, logger, nil)
	storages := NewStorages(nil, nil)
	services := NewServices(components, storages)
	if services == nil {
		t.Fatal("services is nil")
	}
}
