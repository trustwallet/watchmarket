package priority

import "github.com/trustwallet/blockatlas/pkg/errors"

type Controller struct {
	currentProvider uint
	providers       map[uint]string
}

func Init(providersList []string) (Controller, error) {
	if len(providersList) == 0 {
		return Controller{}, errors.E("empty providers list")
	}

	p := make(map[uint]string, len(providersList))
	for i, provider := range providersList {
		p[uint(i)] = provider
	}

	return Controller{currentProvider: 0, providers: p}, nil
}

func (c Controller) GetCurrentProvider() string {
	return c.providers[c.currentProvider]
}

func (c *Controller) GetNextProvider() (string, error) {
	p, ok := c.providers[c.currentProvider+1]
	if !ok {
		return "", errors.E("There is no next provider")
	}
	c.currentProvider = c.currentProvider + 1
	return p, nil
}

func (c Controller) GetAllProviders() []string {
	p := make([]string, len(c.providers))
	for i, c := range c.providers {
		p[i] = c
	}
	return p
}
