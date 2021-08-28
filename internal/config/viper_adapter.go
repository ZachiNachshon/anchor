package config

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/internal/logger"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ConfigViperAdapter interface {
	SetConfigPath(dirPath string) error

	LoadConfigFromFile() error
	LoadConfigFromText(yaml string) error
	RegisterConfigChangesListener(callback func(e fsnotify.Event))

	UpdateAll(cfgToUpdate *AnchorConfig) error
	UpdateEntry(entryName string, value interface{}) error
	GetConfigByKey(key string) string

	SetDefaults() error
	SetEnvVars() error
	AppendConfig(anchorConfig interface{}) error

	flushToNewConfigFile() error
	flush() error
}

type configViperAdapterImpl struct {
	ConfigViperAdapter
	viper *viper.Viper
}

func NewAdapter() ConfigViperAdapter {
	return &configViperAdapterImpl{
		viper: viper.New(),
	}
}

func (a *configViperAdapterImpl) SetConfigPath(dirPath string) error {
	a.viper.SetConfigName(defaultConfigFileName)
	a.viper.SetConfigType(defaultConfigFileType)
	a.viper.AddConfigPath(dirPath) // path to look for the config file in
	return nil
}

func (a *configViperAdapterImpl) LoadConfigFromFile() error {
	if err := a.viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			_ = a.SetDefaults()
			return a.flushToNewConfigFile()
		} else {
			errMsg := fmt.Errorf("could not read configuration from file. error: %s", err.Error())
			logger.Errorf(errMsg.Error())
			return errMsg
		}
	}
	return nil
}

func (a *configViperAdapterImpl) LoadConfigFromText(yaml string) error {
	a.viper.SetConfigType(defaultConfigFileType)
	_ = a.SetDefaults()
	if err := a.viper.ReadConfig(bytes.NewBuffer([]byte(yaml))); err != nil {
		errMsg := fmt.Errorf("failed to read config from buffer. error: %s", err.Error())
		logger.Errorf(errMsg.Error())
		return errMsg
	}
	return nil
}

func (a *configViperAdapterImpl) RegisterConfigChangesListener(callback func(e fsnotify.Event)) {
	a.viper.WatchConfig()
	a.viper.OnConfigChange(callback)
}

func (a *configViperAdapterImpl) flushToNewConfigFile() error {
	return a.viper.SafeWriteConfig()
}

func (a *configViperAdapterImpl) flush() error {
	return a.viper.WriteConfig()
}

func (a *configViperAdapterImpl) UpdateAll(cfgToUpdate *AnchorConfig) error {
	out, err := yaml.Marshal(cfgToUpdate)
	if err != nil {
		return err
	}
	err = a.viper.MergeConfig(bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	return a.flush()
}

func (a *configViperAdapterImpl) UpdateEntry(entryName string, value interface{}) error {
	a.viper.Set(entryName, value)
	if !a.viper.IsSet(entryName) {
		return fmt.Errorf("failed to set configuration entry. name: %s, value: %s", entryName, value)
	}
	return a.flush()
}

func (a *configViperAdapterImpl) GetConfigByKey(key string) string {
	return a.viper.GetString(key)
}

func (a *configViperAdapterImpl) SetDefaults() error {
	a.viper.SetDefault("author", DefaultAuthor)
	a.viper.SetDefault("license", DefaultLicense)
	return nil
}

func (a *configViperAdapterImpl) SetEnvVars() error {
	// Every viper.Get request auto checks for ANCHOR_<flag-name> before reading from config file
	a.viper.SetEnvPrefix("ANCHOR")
	a.viper.AutomaticEnv()
	return nil
}

func (a *configViperAdapterImpl) AppendConfig(anchorConfig interface{}) error {
	var cfg Config
	if err := a.viper.UnmarshalKey("config", &cfg); err != nil {
		return err
	} else {
		anchorConfig.(*AnchorConfig).Config = &cfg
		return nil
	}
}
