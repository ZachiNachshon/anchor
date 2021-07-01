package config

import (
	"bytes"
	"fmt"
	"github.com/ZachiNachshon/anchor/common"
	"github.com/ZachiNachshon/anchor/logger"
	"github.com/ZachiNachshon/anchor/pkg/utils/ioutils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	DefaultAuthor                 = "Zachi Nachshon <zachi.nachshon@gmail.com>"
	DefaultLicense                = "Apache"
	DefaultClonePath              = "${HOME}/.config/anchor/anchorfiles"
	DefaultRemoteBranch           = "master"
	defaultConfigFileName         = "anchorConfig"
	defaultConfigFileType         = "yaml"
	defaultConfigFolderPathFormat = "%s/.config/anchor"
)

func GetConfigFilePath() string {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		logger.Errorf("failed to resolve home folder. err: %s", err.Error())
		return ""
	} else {
		folderPath := fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder)
		return fmt.Sprintf("%s/%s.%s", folderPath, defaultConfigFileName, defaultConfigFileType)
	}
}

func FromContext(ctx common.Context) AnchorConfig {
	return ctx.Config().(AnchorConfig)
}

func SetInContext(ctx common.Context, config AnchorConfig) {
	ctx.(common.ConfigSetter).SetConfig(config)
}

type AnchorConfig struct {
	Config  *Config `yaml:"config"`
	Author  string  `yaml:"author"`
	License string  `yaml:"license"`
}

type Config struct {
	Repository *Repository `yaml:"repository"`
}

type Repository struct {
	Remote *Remote `yaml:"remote"`
	Local  *Local  `yaml:"local"`
}

type Remote struct {
	Url       string `yaml:"url"`
	Revision  string `yaml:"revision"`
	Branch    string `yaml:"branch"`
	ClonePath string `yaml:"clonePath"`
}

type Local struct {
	Path string `yaml:"path"`
}

func initConfigPath() error {
	if homeFolder, err := ioutils.GetUserHomeFolder(); err != nil {
		return err
	} else {
		viper.SetConfigName(defaultConfigFileName)
		viper.SetConfigType(defaultConfigFileType)
		viper.AddConfigPath(fmt.Sprintf(defaultConfigFolderPathFormat, homeFolder)) // path to look for the config file in
		//viper.AddConfigPath(".")                      		// optionally look for config in the working directory
		return nil
	}
}

func setDefaults() {
	viper.SetDefault("author", DefaultAuthor)
	viper.SetDefault("license", DefaultLicense)
}

func createConfigFileWithDefaults() {
	err := viper.SafeWriteConfig() // Write defaults
	if err != nil {
		logger.Errorf("Could not create config file with defaults: %s \n", err)
	}
}

func listenOnConfigFileChanges() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Suggest to git fetch the repository
		logger.Infof("Config file changed:", e.Name)
	})
}

func createConfigObject() *AnchorConfig {
	var config Config
	if err := viper.UnmarshalKey("config", &config); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to unmarshal configuration file. error: %s \n", err))
	}

	setDefaultsPostCreation(&config)

	return &AnchorConfig{
		Config:  &config,
		Author:  viper.GetString("author"),
		License: viper.GetString("license"),
	}
}

func setDefaultsPostCreation(config *Config) {
	if config.Repository != nil && config.Repository.Remote != nil {
		if config.Repository.Remote.ClonePath == "" {
			config.Repository.Remote.ClonePath = DefaultClonePath
		}

		if config.Repository.Remote.Branch == "" {
			config.Repository.Remote.Branch = DefaultRemoteBranch
		}
	}
}

var ViperConfigInMemoryLoader = func(yaml string) (*AnchorConfig, error) {
	viper.SetConfigType("yaml")
	setDefaults()

	if err := viper.ReadConfig(bytes.NewBuffer([]byte(yaml))); err != nil {
		logger.Errorf("Failed to read config from buffer. error: %s", err)
		return nil, err
	}

	return createConfigObject(), nil
}

var ViperConfigFileLoader = func() (*AnchorConfig, error) {
	err := initConfigPath()
	if err != nil {
		logger.Fatalf("Failed to initialize anchor configuration path. error: %s", err.Error())
		return nil, err
	}

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Handle errors reading the config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			setDefaults()
			createConfigFileWithDefaults()
			listenOnConfigFileChanges()
		} else {
			logger.Errorf("Config file was found but an error occurred. error: %s", err)
			return nil, err
		}
	}

	// Every viper.Get request auto checks for ANCHOR_<flag-name> before reading from config file
	viper.SetEnvPrefix("ANCHOR")
	viper.AutomaticEnv()

	return createConfigObject(), nil
}
