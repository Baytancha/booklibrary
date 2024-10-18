package components

import (
	"test/internal/infrastructure/responder"

	"github.com/jmoiron/sqlx"
	"github.com/ptflp/godecoder"
	"go.uber.org/zap"
)

type Components struct {
	Responder responder.Responder
	Decoder   godecoder.Decoder
	Logger    *zap.Logger
	DB        *sqlx.DB
}

func NewComponents(responder responder.Responder, decoder godecoder.Decoder, logger *zap.Logger, db *sqlx.DB) *Components {
	return &Components{
		Responder: responder,
		Decoder:   decoder,
		Logger:    logger,
		DB:        db,
	}
}
