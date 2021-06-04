package fakes

import (
	"github.com/ZachiNachshon/anchor/config"
)

var FakeConfigLoader = func(content string) (*config.AnchorConfig, error) {
	if cfg, err := config.ViperConfigInMemoryLoader(content); err != nil {
		return nil, err
	} else {
		return cfg, nil
	}
}
