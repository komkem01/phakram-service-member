package kafka

import (
	"phakram/internal/config"
)

type Module struct {
	Svc *Service
}

func New(conf *config.Config[Config]) *Module {
	return &Module{
		Svc: newService(conf),
	}
}
