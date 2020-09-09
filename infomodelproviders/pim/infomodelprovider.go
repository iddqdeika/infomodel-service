package pim

import (
	"fmt"
	"github.com/iddqdeika/pim"
	"infomodel-service/definitions"
)

const (
	defaultTimeoutInSecs = 10

	infomodelStructureID = 9001
)

func NewInfomodelProvider(config definitions.Config) (definitions.InfomodelProvider, error) {

	if config == nil {
		return nil, fmt.Errorf("config must not be nil")
	}

	host, err := config.GetString("host")
	if err != nil {
		return nil, err
	}

	login, err := config.GetString("login")
	if err != nil {
		return nil, err
	}

	pass, err := config.GetString("password")
	if err != nil {
		return nil, err
	}

	time, err := config.GetInt("timeout_in_seconds")
	if err != nil {
		return nil, err
	}
	if time <= 0 {
		time = defaultTimeoutInSecs
	}

	c, err := pim.NewClient(pim.Config{
		Host:          host,
		Login:         login,
		Password:      pass,
		TimeoutInSecs: time,
	})
	return &pimRestInfomodelProvider{c: c}, nil
}

type pimRestInfomodelProvider struct {
	c *pim.Client
}

func (p *pimRestInfomodelProvider) GetByIdentifier(identifier string) (*pim.StructureGroup, error) {
	im, err := p.c.StructureGroupProvider().GetInfomodelByIdentifier(identifier, infomodelStructureID)
	if err != nil {
		return nil, err
	}
	return im, nil
}
