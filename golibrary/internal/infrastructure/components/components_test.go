package components

import (
	"os"
	"testing"

	"test/config"
	"test/internal/db"
	"test/internal/infrastructure/responder"

	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

func TestNewComponents(t *testing.T) {
	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	responseManager := responder.NewResponder(decoder, zap.NewNop())
	cfg := config.NewConfig(config.WithPort(8080), config.WithDBname("test"), config.WithDSN(os.Getenv("DB_DSN2")))
	dbx, err := db.NewSqlDB(cfg, zap.NewNop())
	if err != nil {
		t.Fatal(err)
	}
	components := NewComponents(responseManager, decoder, zap.NewNop(), dbx)

	if components == nil {
		t.Fatal("components is nil")
	}
}
