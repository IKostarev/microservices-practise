package handlers

import "github.com/rs/zerolog"

type GatewayHandler struct {
	logger         *zerolog.Logger
	gatewayService GatewayService
}

func NewGatewayHandler(
	logger *zerolog.Logger,
	gatewayService GatewayService,
) *GatewayHandler {
	return &GatewayHandler{
		logger:         logger,
		gatewayService: gatewayService,
	}
}
