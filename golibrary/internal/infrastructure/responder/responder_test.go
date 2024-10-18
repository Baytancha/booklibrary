package responder

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

func TestResponder(t *testing.T) {

	decoder := godecoder.NewDecoder(jsoniter.Config{
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		DisallowUnknownFields:  true,
	})
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	r := NewResponder(decoder, logger)
	if r == nil {
		t.Fatal("r is nil")
	}
}
