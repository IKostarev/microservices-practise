package rest

import (
	rest "gateway/internal/api"
	"github.com/rs/zerolog"
)

type GatewayHandler struct {
	logger         *zerolog.Logger
	gatewayService rest.GatewayService
}

func NewGatewayHandler(
	logger *zerolog.Logger,
	gatewayService rest.GatewayService,
) *GatewayHandler {
	return &GatewayHandler{
		logger:         logger,
		gatewayService: gatewayService,
	}
}
